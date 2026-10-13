package tuyile

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fatWill/agent-team-blog/backend/config"
)

// chatMessage OpenAI 兼容消息（含工具调用往返所需字段）
type chatMessage struct {
	Role             string     `json:"role"`
	Content          string     `json:"content"`
	ReasoningContent string     `json:"reasoning_content,omitempty"` // GLM 交错思考：工具往返时回传本轮思考
	ToolCalls        []toolCall `json:"tool_calls,omitempty"`
	ToolCallID       string     `json:"tool_call_id,omitempty"`
}

type toolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

// toolDef OpenAI function tool 定义
type toolDef struct {
	Type     string `json:"type"`
	Function struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Parameters  json.RawMessage `json:"parameters"`
	} `json:"function"`
}

type chatChunk struct {
	Choices []struct {
		Delta struct {
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
			ToolCalls        []struct {
				Index    int    `json:"index"`
				ID       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// llmError 上游错误（code 透传给 AG-UI RUN_ERROR.code）
type llmError struct {
	Code string
	Msg  string
}

func (e *llmError) Error() string { return e.Code + ": " + e.Msg }

var (
	errTimeout = func() *llmError {
		return &llmError{Code: "LLM_TIMEOUT", Msg: "图一乐掐指算得太久了（大模型响应超时），换个简单点的问法再试试～"}
	}
)

// streamHandlers 流式回调：思考增量（不下发原文，只用于状态）、正文增量
type streamHandlers struct {
	OnReasoning func(delta string) error
	OnContent   func(delta string) error
}

// streamResult 单轮结果
type streamResult struct {
	Content      string
	Reasoning    string
	ToolCalls    []toolCall
	FinishReason string
}

// llmClient OpenAI 兼容 chat/completions 流式客户端
type llmClient struct {
	cfg  config.LLMConfig
	http *http.Client
	once sync.Once
}

// newTransport 复用连接池并显式启用 HTTP/2（TLS ALPN 协商 h2，多轮工具调用复用同一连接）。
// 总超时由每轮请求的 context 控制（流式响应不能用 Client.Timeout 截断）。
func newTransport(headerTimeout time.Duration) *http.Transport {
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          8,
		MaxIdleConnsPerHost:   4,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
		ResponseHeaderTimeout: headerTimeout,
	}
}

func newLLMClient(cfg config.LLMConfig) *llmClient {
	if cfg.TimeoutSeconds <= 0 {
		cfg.TimeoutSeconds = 120
	}
	if cfg.HeaderTimeoutSeconds <= 0 {
		cfg.HeaderTimeoutSeconds = 90
	}
	if cfg.MaxToolRounds <= 0 {
		cfg.MaxToolRounds = 3
	}
	return &llmClient{cfg: cfg, http: &http.Client{Transport: newTransport(time.Duration(cfg.HeaderTimeoutSeconds) * time.Second)}}
}

// configured 是否已配置上游
func (c *llmClient) configured() error {
	switch {
	case c.cfg.APIKey == "":
		return &llmError{Code: "LLM_NOT_CONFIGURED", Msg: "图一乐还没接上大模型（服务端未配置 LLM_API_KEY），暂时算不了，稍后再来～"}
	case c.cfg.BaseURL == "" || c.cfg.Model == "":
		return &llmError{Code: "LLM_NOT_CONFIGURED", Msg: "服务端未配置 LLM_BASE_URL / LLM_MODEL"}
	case strings.Contains(strings.ToLower(c.cfg.BaseURL), "woa.com"):
		// 生产服务器无法访问公司内网，直接拒绝避免长时间挂起
		return &llmError{Code: "LLM_NOT_CONFIGURED", Msg: "LLM_BASE_URL 不可使用内网地址"}
	}
	return nil
}

// buildBody 组装请求体：标准字段 + 可选 reasoning_effort / thinking + LLM_EXTRA_BODY 合并
func (c *llmClient) buildBody(msgs []chatMessage, tools []toolDef) ([]byte, error) {
	body := map[string]any{
		"model":       c.cfg.Model,
		"messages":    msgs,
		"stream":      true,
		"temperature": 0.8,
	}
	if c.cfg.MaxTokens > 0 {
		body["max_tokens"] = c.cfg.MaxTokens
	}
	if len(tools) > 0 {
		body["tools"] = tools
		body["tool_choice"] = "auto"
	}
	if c.cfg.ReasoningEffort != "" {
		body["reasoning_effort"] = c.cfg.ReasoningEffort
	}
	if c.cfg.ThinkingType != "" {
		body["thinking"] = map[string]string{"type": c.cfg.ThinkingType}
	}
	if c.cfg.ExtraBodyJSON != "" {
		var extra map[string]any
		if err := json.Unmarshal([]byte(c.cfg.ExtraBodyJSON), &extra); err != nil {
			return nil, &llmError{Code: "LLM_NOT_CONFIGURED", Msg: "LLM_EXTRA_BODY 不是合法 JSON 对象"}
		}
		for k, v := range extra {
			body[k] = v
		}
	}
	return json.Marshal(body)
}

// isTimeout 判断是否超时：context 截止、Transport 的 ResponseHeaderTimeout（net.Error.Timeout）、读超时
func isTimeout(ctx context.Context, err error) bool {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var ne net.Error
	return errors.As(err, &ne) && ne.Timeout()
}

// classifyTransportErr 把 http.Do / 读流错误归类为 AG-UI RUN_ERROR.code
func classifyTransportErr(ctx context.Context, err error, reading bool) error {
	switch {
	case isTimeout(ctx, err):
		return errTimeout()
	case errors.Is(ctx.Err(), context.Canceled):
		return context.Canceled // 客户端断开
	case reading:
		return &llmError{Code: "LLM_STREAM_BROKEN", Msg: "大模型流式响应中断"}
	default:
		return &llmError{Code: "LLM_UNAVAILABLE", Msg: "大模型服务连接失败"}
	}
}

// stream 发起单轮流式补全；正文/思考增量走回调，工具调用累积后在结果中返回
func (c *llmClient) stream(ctx context.Context, msgs []chatMessage, tools []toolDef, h streamHandlers) (*streamResult, error) {
	if err := c.configured(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.cfg.TimeoutSeconds)*time.Second)
	defer cancel()

	body, err := c.buildBody(msgs, tools)
	if err != nil {
		return nil, err
	}
	url := c.cfg.BaseURL
	if !strings.HasSuffix(url, "/chat/completions") {
		url += "/chat/completions"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, &llmError{Code: "LLM_NOT_CONFIGURED", Msg: "LLM_BASE_URL 非法"}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, classifyTransportErr(ctx, err, false)
	}
	defer resp.Body.Close()
	c.once.Do(func() { logf("tuyile.llm.proto=%s model=%s", resp.Proto, c.cfg.Model) })
	if resp.StatusCode != http.StatusOK {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		code := "LLM_UPSTREAM_ERROR"
		switch resp.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			code = "LLM_AUTH_FAILED"
		case http.StatusTooManyRequests:
			code = "LLM_RATE_LIMITED"
		case http.StatusGatewayTimeout, http.StatusRequestTimeout:
			return nil, errTimeout()
		}
		return nil, &llmError{Code: code, Msg: fmt.Sprintf("上游返回 HTTP %d: %s", resp.StatusCode, sanitize(string(snippet)))}
	}

	res := &streamResult{}
	var content, reasoning strings.Builder
	acc := map[int]*toolCall{}
	finish := func() *streamResult {
		res.Content, res.Reasoning = content.String(), reasoning.String()
		idx := make([]int, 0, len(acc))
		for i := range acc {
			idx = append(idx, i)
		}
		sort.Ints(idx)
		for _, i := range idx {
			if acc[i].Function.Name != "" {
				res.ToolCalls = append(res.ToolCalls, *acc[i])
			}
		}
		return res
	}

	reader := bufio.NewReaderSize(resp.Body, 16*1024)
	for {
		line, rerr := reader.ReadString('\n')
		if len(line) > 0 {
			line = strings.TrimRight(line, "\r\n")
			if strings.HasPrefix(line, "data:") {
				data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				if data == "[DONE]" {
					return finish(), nil
				}
				var ch chatChunk
				if json.Unmarshal([]byte(data), &ch) == nil {
					if ch.Error != nil {
						return nil, &llmError{Code: "LLM_UPSTREAM_ERROR", Msg: sanitize(ch.Error.Message)}
					}
					for _, choice := range ch.Choices {
						d := choice.Delta
						if d.ReasoningContent != "" {
							reasoning.WriteString(d.ReasoningContent)
							if h.OnReasoning != nil {
								if e := h.OnReasoning(d.ReasoningContent); e != nil {
									return nil, e
								}
							}
						}
						if d.Content != "" {
							content.WriteString(d.Content)
							if h.OnContent != nil {
								if e := h.OnContent(d.Content); e != nil {
									return nil, e
								}
							}
						}
						for _, tc := range d.ToolCalls {
							t, ok := acc[tc.Index]
							if !ok {
								t = &toolCall{Type: "function"}
								acc[tc.Index] = t
							}
							if tc.ID != "" {
								t.ID = tc.ID
							}
							if tc.Function.Name != "" {
								t.Function.Name = tc.Function.Name
							}
							t.Function.Arguments += tc.Function.Arguments
						}
						if choice.FinishReason != nil && *choice.FinishReason != "" {
							res.FinishReason = *choice.FinishReason
						}
					}
				}
			}
		}
		if rerr != nil {
			if rerr == io.EOF {
				return finish(), nil
			}
			return nil, classifyTransportErr(ctx, rerr, true)
		}
	}
}

// sanitize 截断上游错误片段，避免日志/前端回显过长内容
func sanitize(s string) string {
	s = strings.TrimSpace(s)
	if r := []rune(s); len(r) > 200 {
		s = string(r[:200])
	}
	return strings.ReplaceAll(s, "\n", " ")
}
