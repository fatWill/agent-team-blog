package tuyile

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// AG-UI 事件类型（https://docs.ag-ui.com/concepts/events）
const (
	EventRunStarted         = "RUN_STARTED"
	EventRunFinished        = "RUN_FINISHED"
	EventRunError           = "RUN_ERROR"
	EventTextMessageStart   = "TEXT_MESSAGE_START"
	EventTextMessageContent = "TEXT_MESSAGE_CONTENT"
	EventTextMessageEnd     = "TEXT_MESSAGE_END"
)

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
