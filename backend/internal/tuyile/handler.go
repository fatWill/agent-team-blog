package tuyile

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/fatWill/agent-team-blog/backend/config"
	"github.com/fatWill/agent-team-blog/backend/internal/almanac"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/fatWill/agent-team-blog/backend/pkg/middleware"
)

const (
	// MaxMessageRunes 单条消息最大字符数
	MaxMessageRunes = 500
	// MaxHistoryMessages 最多保留的历史消息条数（user+assistant，约 10 轮）
	MaxHistoryMessages = 20
	// MaxBodyBytes 请求体上限
	MaxBodyBytes = 64 << 10
	// MaxAssistantRunes 历史 assistant 消息截断长度（控制上游 token 成本）
	MaxAssistantRunes = 1200
	// FutureDaysWithTools 工具可用时注入的未来天数精简表（控制提示词长度与首字延迟）
	FutureDaysWithTools = 7
	// FutureDaysFallback 工具不可用（LLM_DISABLE_TOOLS=1）时的兜底注入天数
	FutureDaysFallback = 30
	keepAliveInterval  = 15 * time.Second
)

var (
	llm     *llmClient
	limiter *ipLimiter
	sem     chan struct{}
	nowFn   = time.Now
	logf    = log.Printf
)

// Init 注入配置（main 启动时调用一次）
func Init(llmCfg config.LLMConfig, cfg config.TuyileConfig) {
	llm = newLLMClient(llmCfg)
	limiter = newIPLimiter(cfg.RatePerMinute, cfg.RatePerDay)
	if cfg.MaxConcurrency <= 0 {
		cfg.MaxConcurrency = 8
	}
	sem = make(chan struct{}, cfg.MaxConcurrency)
	initShare(cfg)
	if err := llm.configured(); err != nil {
		logf("⚠️ tuyile: %v（chat 接口将返回 RUN_ERROR）", err)
	}
}

// runInput AG-UI RunAgentInput（只取本服务用到的字段，其余如 tools/state/context 忽略）
type runInput struct {
	ThreadID string         `json:"threadId"`
	RunID    string         `json:"runId"`
	Messages []inputMessage `json:"messages"`
}

type inputMessage struct {
	ID      string          `json:"id"`
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

// text 兼容 content 为 string 或 ContentPart[]（取 type=text 的片段）
func (m inputMessage) text() string {
	if len(m.Content) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(m.Content, &s) == nil {
		return s
	}
	var parts []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if json.Unmarshal(m.Content, &parts) == nil {
		var b strings.Builder
		for _, p := range parts {
			if p.Type == "text" {
				b.WriteString(p.Text)
			}
		}
		return b.String()
	}
	return ""
}

func reject(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": true, "statusCode": status, "statusMessage": msg})
}

// Chat POST /api/toys/tuyile/chat — AG-UI 协议（HTTP + SSE）运势聊天
//
// 运行开始前的输入错误（JSON 非法、消息超长、限流）按 AG-UI 绑定约定返回 HTTP 错误、无事件流；
// 运行开始后的错误（未配置 key、上游失败）以 RUN_ERROR 事件在流内返回。
// query_almanac 工具调用在服务端闭环完成，对前端只透出 CUSTOM tuyile.status 状态。
func Chat(c *gin.Context) {
	if llm == nil {
		reject(c, http.StatusServiceUnavailable, "服务未初始化")
		return
	}
	if ok, reason := limiter.allow(middleware.GetClientIP(c)); !ok {
		msg := "算得太勤快啦，歇一分钟再来～"
		if reason == "day" {
			msg = "今天的卦已经算满啦，明天再来图一乐～"
		}
		c.Header("Retry-After", ternary(reason == "day", "3600", "60"))
		reject(c, http.StatusTooManyRequests, msg)
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	var in runInput
	if err := json.NewDecoder(c.Request.Body).Decode(&in); err != nil {
		reject(c, http.StatusBadRequest, "请求体不是合法的 RunAgentInput JSON")
		return
	}
	msgs, userTexts, opening, err := normalizeMessages(in.Messages)
	if err != nil {
		reject(c, http.StatusBadRequest, err.Error())
		return
	}
	if in.ThreadID = strings.TrimSpace(in.ThreadID); in.ThreadID == "" || len(in.ThreadID) > 128 {
		in.ThreadID = uuid.NewString()
	}
	if in.RunID = strings.TrimSpace(in.RunID); in.RunID == "" || len(in.RunID) > 128 {
		in.RunID = uuid.NewString()
	}

	select {
	case sem <- struct{}{}:
		defer func() { <-sem }()
	default:
		c.Header("Retry-After", "5")
		reject(c, http.StatusServiceUnavailable, "图一乐小摊排队中，稍后再试～")
		return
	}

	start := nowFn()
	w := newSSEWriter(c)
	_ = w.send(aguiEvent{Type: EventRunStarted, ThreadID: in.ThreadID, RunID: in.RunID})

	ctx := c.Request.Context()
	if err := llm.configured(); err != nil {
		sendRunError(w, err)
		return
	}
	em := &runEmitter{w: w, newID: uuid.NewString, textID: "msg_" + uuid.NewString()}
	// 先发状态：前端立刻显示「掐指一算中…」，不必等上游首包
	_ = em.status(statusValue{Stage: "thinking", Text: "掐指一算中…"})

	// 空闲 keep-alive，避免思考/工具阶段长时间无字节被 Nginx/浏览器断开
	stopKA := make(chan struct{})
	go func() {
		t := time.NewTicker(keepAliveInterval)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				w.comment("keep-alive")
			case <-stopKA:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
	defer close(stopKA)

	repo := almanac.NewRepo(db.DB)
	facts := collectFacts(userTexts)
	var ranges []dateRange
	if len(userTexts) > 0 && !opening {
		ranges = parseRelativeDates(userTexts[len(userTexts)-1], nowFn())
	}
	useTools := !llm.cfg.DisableTools && !opening
	futureDays := FutureDaysWithTools
	if !useTools {
		futureDays = FutureDaysFallback
	}
	if opening {
		futureDays = 0
	}
	knowledge := buildKnowledge(repo, facts, ranges, nowFn(), futureDays)

	upstream := make([]chatMessage, 0, len(msgs)+4)
	upstream = append(upstream, chatMessage{Role: "system", Content: systemPrompt}, chatMessage{Role: "system", Content: knowledge})
	if opening {
		upstream = append(upstream, chatMessage{Role: "system", Content: openingPrompt})
	}
	upstream = append(upstream, msgs...)

	stats, streamErr := runAgent(ctx, repo, em, upstream, useTools)
	em.closeAll()
	switch {
	case streamErr == nil && !em.textOpen:
		sendRunError(w, &llmError{Code: "LLM_EMPTY_RESPONSE", Msg: "大模型没有返回内容，换个问法试试～"})
	case streamErr == nil:
		_ = w.send(aguiEvent{Type: EventRunFinished, ThreadID: in.ThreadID, RunID: in.RunID})
	case errors.Is(streamErr, context.Canceled):
		// 客户端主动断开，无需再写
	default:
		sendRunError(w, streamErr)
	}
	// 日志不含姓名、生日、消息原文
	logf("tuyile.chat.done run=%s msgs=%d opening=%v out_chars=%d reasoning_chars=%d tool_calls=%d rounds=%d first_byte=%s birth=%v year=%d ranges=%d cost=%s err=%s",
		in.RunID, len(msgs), opening, em.textChars, stats.reasoningChars, stats.toolCalls, stats.rounds,
		stats.firstContent.Truncate(time.Millisecond), facts.Birth != nil, facts.TargetYear, len(ranges),
		time.Since(start).Truncate(time.Millisecond), errCode(streamErr))
}

type agentStats struct {
	rounds         int
	toolCalls      int
	reasoningChars int
	firstContent   time.Duration
}

// runAgent 服务端闭环的工具调用循环：模型请求 query_almanac → 服务端查库 → 回填 tool 消息 → 继续，最多 MaxToolRounds 轮
func runAgent(ctx context.Context, repo *almanac.Repo, em *runEmitter, msgs []chatMessage, useTools bool) (agentStats, error) {
	var st agentStats
	begin := time.Now()
	var tools []toolDef
	if useTools {
		tools = []toolDef{almanacTool}
	}
	for round := 0; ; round++ {
		st.rounds = round + 1
		roundTools := tools
		if round >= llm.cfg.MaxToolRounds {
			roundTools = nil // 轮数用尽：不再给工具，强制模型直接作答
		}
		res, err := llm.stream(ctx, msgs, roundTools, streamHandlers{
			OnReasoning: func(delta string) error {
				st.reasoningChars += utf8.RuneCountInString(delta)
				return em.status(statusValue{Stage: "thinking", Text: "掐指一算中…"})
			},
			OnContent: func(delta string) error {
				if st.firstContent == 0 {
					st.firstContent = time.Since(begin)
				}
				return em.text(delta, utf8.RuneCountInString(delta))
			},
		})
		if err != nil {
			return st, err
		}
		if len(res.ToolCalls) == 0 || roundTools == nil {
			return st, nil
		}
		// 回填 assistant(tool_calls) + tool 结果；保留本轮 reasoning_content（GLM 交错思考要求）
		msgs = append(msgs, chatMessage{Role: "assistant", Content: res.Content, ReasoningContent: res.Reasoning, ToolCalls: res.ToolCalls})
		for _, tc := range res.ToolCalls {
			st.toolCalls++
			result := "未知工具：" + tc.Function.Name
			if tc.Function.Name == toolQueryAlmanac {
				var a almanacArgs
				_ = json.Unmarshal([]byte(tc.Function.Arguments), &a)
				rng := a.Start
				if a.End != "" && a.End != a.Start {
					rng += "~" + a.End
				}
				if err := em.status(statusValue{Stage: "querying", Text: "翻黄历中…", Tool: toolQueryAlmanac, Range: rng}); err != nil {
					return st, err
				}
				result = runAlmanacTool(repo, tc.Function.Arguments)
			}
			msgs = append(msgs, chatMessage{Role: "tool", ToolCallID: tc.ID, Content: result})
		}
		if err := em.status(statusValue{Stage: "thinking", Text: "掐指一算中…"}); err != nil {
			return st, err
		}
	}
}

func sendRunError(w *sseWriter, err error) {
	var le *llmError
	if errors.As(err, &le) {
		_ = w.send(aguiEvent{Type: EventRunError, Message: le.Msg, Code: le.Code})
		return
	}
	_ = w.send(aguiEvent{Type: EventRunError, Message: "图一乐卦摊出了点小状况，请稍后再试", Code: "INTERNAL_ERROR"})
}

func errCode(err error) string {
	var le *llmError
	if errors.As(err, &le) {
		return le.Code
	}
	if errors.Is(err, context.Canceled) {
		return "client_canceled"
	}
	if err != nil {
		return "other"
	}
	return ""
}

// normalizeMessages 只保留 user/assistant（丢弃客户端注入的 system/developer/tool），截断历史并校验长度。
// messages 为空时返回 opening=true，并以占位 user 发言触发开场白。
func normalizeMessages(in []inputMessage) ([]chatMessage, []string, bool, error) {
	out := make([]chatMessage, 0, len(in))
	for _, m := range in {
		role := strings.ToLower(strings.TrimSpace(m.Role))
		if role != "user" && role != "assistant" {
			continue
		}
		text := strings.TrimSpace(m.text())
		if text == "" {
			continue
		}
		out = append(out, chatMessage{Role: role, Content: text})
	}
	if len(out) > MaxHistoryMessages {
		out = out[len(out)-MaxHistoryMessages:]
	}
	if len(out) == 0 {
		return []chatMessage{{Role: "user", Content: openingUserTurn}}, nil, true, nil
	}
	if out[len(out)-1].Role != "user" {
		return nil, nil, false, errors.New("最后一条消息必须是 user")
	}
	if utf8.RuneCountInString(out[len(out)-1].Content) > MaxMessageRunes {
		return nil, nil, false, errors.New("单条消息最多 500 字")
	}
	userTexts := make([]string, 0, len(out))
	for i := range out {
		limit := MaxMessageRunes
		if out[i].Role == "assistant" {
			limit = MaxAssistantRunes
		}
		out[i].Content = truncateRunes(out[i].Content, limit)
		if out[i].Role == "user" {
			userTexts = append(userTexts, out[i].Content)
		}
	}
	return out, userTexts, false, nil
}

func truncateRunes(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	return string([]rune(s)[:n])
}
