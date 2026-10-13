package tuyile

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// AG-UI 事件类型（https://docs.ag-ui.com/concepts/events，1.0 Reasoning 事件取代旧 THINKING_*）
const (
	EventRunStarted              = "RUN_STARTED"
	EventRunFinished             = "RUN_FINISHED"
	EventRunError                = "RUN_ERROR"
	EventTextMessageStart        = "TEXT_MESSAGE_START"
	EventTextMessageContent      = "TEXT_MESSAGE_CONTENT"
	EventTextMessageEnd          = "TEXT_MESSAGE_END"
	EventReasoningStart          = "REASONING_START"
	EventReasoningMessageStart   = "REASONING_MESSAGE_START"
	EventReasoningMessageContent = "REASONING_MESSAGE_CONTENT"
	EventReasoningMessageEnd     = "REASONING_MESSAGE_END"
	EventReasoningEnd            = "REASONING_END"
	EventCustom                  = "CUSTOM"
)

// CustomStatusName 图一乐状态 CUSTOM 事件名
const CustomStatusName = "tuyile.status"

// aguiEvent AG-UI 事件（字段 camelCase，可选字段 omitempty）
type aguiEvent struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp,omitempty"`
	ThreadID  string `json:"threadId,omitempty"`
	RunID     string `json:"runId,omitempty"`
	MessageID string `json:"messageId,omitempty"`
	Role      string `json:"role,omitempty"`
	Delta     string `json:"delta,omitempty"`
	Message   string `json:"message,omitempty"`
	Code      string `json:"code,omitempty"`
	Name      string `json:"name,omitempty"`  // CUSTOM
	Value     any    `json:"value,omitempty"` // CUSTOM
}

// statusValue CUSTOM tuyile.status 的 value
type statusValue struct {
	Stage string `json:"stage"`          // thinking / querying / answering
	Text  string `json:"text"`           // 前端可直接展示的状态文案
	Tool  string `json:"tool,omitempty"` // querying 阶段的工具名
	Range string `json:"range,omitempty"`
}

// sseWriter 按 AG-UI HTTP+SSE 绑定写事件：每个事件一行 `data: <json>\n\n`，LF 换行
type sseWriter struct {
	c  *gin.Context
	mu sync.Mutex
}

func newSSEWriter(c *gin.Context) *sseWriter {
	h := c.Writer.Header()
	h.Set("Content-Type", "text/event-stream; charset=utf-8")
	h.Set("Cache-Control", "no-cache, no-transform")
	h.Set("Connection", "keep-alive")
	h.Set("X-Accel-Buffering", "no") // 关闭 Nginx 代理缓冲，保证逐字推送
	c.Status(http.StatusOK)
	c.Writer.Flush()
	return &sseWriter{c: c}
}

func (w *sseWriter) send(ev aguiEvent) error {
	ev.Timestamp = time.Now().UnixMilli()
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if _, err := w.c.Writer.Write(append(append([]byte("data: "), b...), '\n', '\n')); err != nil {
		return err
	}
	w.c.Writer.Flush()
	return nil
}

// comment 发送 SSE 注释行（keep-alive），AG-UI 规定消费方必须容忍
func (w *sseWriter) comment(s string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, _ = w.c.Writer.Write([]byte(": " + s + "\n\n"))
	w.c.Writer.Flush()
}

// runEmitter 维护一次 run 的事件状态机，保证 AG-UI 配对规则：
//   - reasoning span / reasoning message / text message 都成对开闭，run 结束前全部关闭；
//   - 思考原文不下发：reasoning message 只承载「状态摘要」（AG-UI Reasoning 的 summary-only 可见性），
//     同时发 CUSTOM tuyile.status 便于前端直接渲染状态条。
type runEmitter struct {
	w           *sseWriter
	newID       func() string
	spanID      string // 当前打开的 reasoning span
	reasonMsgID string // 当前打开的 reasoning message
	textID      string // assistant 文本消息 id（整个 run 只有一条）
	textOpen    bool
	textChars   int
	stage       string
}

func (e *runEmitter) status(v statusValue) error {
	if v.Stage == e.stage && v.Stage != "querying" {
		return nil
	}
	e.stage = v.Stage
	if !e.textOpen {
		if e.spanID == "" {
			e.spanID = "rsn_" + e.newID()
			if err := e.w.send(aguiEvent{Type: EventReasoningStart, MessageID: e.spanID}); err != nil {
				return err
			}
		}
		if e.reasonMsgID != "" {
			if err := e.w.send(aguiEvent{Type: EventReasoningMessageEnd, MessageID: e.reasonMsgID}); err != nil {
				return err
			}
		}
		e.reasonMsgID = "rsm_" + e.newID()
		if err := e.w.send(aguiEvent{Type: EventReasoningMessageStart, MessageID: e.reasonMsgID, Role: "reasoning"}); err != nil {
			return err
		}
		if err := e.w.send(aguiEvent{Type: EventReasoningMessageContent, MessageID: e.reasonMsgID, Delta: v.Text}); err != nil {
			return err
		}
	}
	return e.w.send(aguiEvent{Type: EventCustom, Name: CustomStatusName, Value: v})
}

// closeReasoning 关闭 reasoning message 与 span
func (e *runEmitter) closeReasoning() error {
	if e.reasonMsgID != "" {
		if err := e.w.send(aguiEvent{Type: EventReasoningMessageEnd, MessageID: e.reasonMsgID}); err != nil {
			return err
		}
		e.reasonMsgID = ""
	}
	if e.spanID != "" {
		if err := e.w.send(aguiEvent{Type: EventReasoningEnd, MessageID: e.spanID}); err != nil {
			return err
		}
		e.spanID = ""
	}
	return nil
}

func (e *runEmitter) text(delta string, n int) error {
	if !e.textOpen {
		if err := e.closeReasoning(); err != nil {
			return err
		}
		e.textOpen = true
		if err := e.w.send(aguiEvent{Type: EventTextMessageStart, MessageID: e.textID, Role: "assistant"}); err != nil {
			return err
		}
	}
	e.textChars += n
	return e.w.send(aguiEvent{Type: EventTextMessageContent, MessageID: e.textID, Delta: delta})
}

// closeAll run 结束前关闭所有打开的消息与 span
func (e *runEmitter) closeAll() {
	_ = e.closeReasoning()
	if e.textOpen {
		_ = e.w.send(aguiEvent{Type: EventTextMessageEnd, MessageID: e.textID})
	}
}
