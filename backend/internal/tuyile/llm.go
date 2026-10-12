package tuyile

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fatWill/agent-team-blog/backend/config"
)

// chatMessage OpenAI 兼容消息
type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Stream      bool          `json:"stream"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

type chatChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
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

// llmClient OpenAI 兼容 chat/completions 流式客户端
type llmClient struct {
	cfg  config.LLMConfig
	http *http.Client
}

// 复用连接池；超时由每次请求的 context 控制（流式响应不能用 Client.Timeout 截断）
var sharedTransport = &http.Transport{
	Proxy:                 http.ProxyFromEnvironment,
	MaxIdleConns:          8,
	MaxIdleConnsPerHost:   4,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ResponseHeaderTimeout: 30 * time.Second,
}

func newLLMClient(cfg config.LLMConfig) *llmClient {
	return &llmClient{cfg: cfg, http: &http.Client{Transport: sharedTransport}}
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

// stream 发起流式补全，逐段回调 onDelta；返回 nil 表示正常结束
func (c *llmClient) stream(ctx context.Context, msgs []chatMessage, onDelta func(string) error) error {
	if err := c.configured(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.cfg.TimeoutSeconds)*time.Second)
	defer cancel()

	body, _ := json.Marshal(chatRequest{
		Model: c.cfg.Model, Messages: msgs, Stream: true, Temperature: 0.8, MaxTokens: c.cfg.MaxTokens,
	})
	url := c.cfg.BaseURL
	if !strings.HasSuffix(url, "/chat/completions") {
		url += "/chat/completions"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return &llmError{Code: "LLM_REQUEST_FAILED", Msg: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)

	resp, err := c.http.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return &llmError{Code: "LLM_TIMEOUT", Msg: "大模型响应超时"}
		}
		return &llmError{Code: "LLM_UNAVAILABLE", Msg: "大模型服务连接失败"}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		code := "LLM_UPSTREAM_ERROR"
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			code = "LLM_AUTH_FAILED"
		} else if resp.StatusCode == http.StatusTooManyRequests {
			code = "LLM_RATE_LIMITED"
		}
		return &llmError{Code: code, Msg: fmt.Sprintf("上游返回 HTTP %d: %s", resp.StatusCode, sanitize(string(snippet)))}
	}

	reader := bufio.NewReaderSize(resp.Body, 16*1024)
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			line = strings.TrimRight(line, "\r\n")
			if strings.HasPrefix(line, "data:") {
				data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				if data == "[DONE]" {
					return nil
				}
				var ch chatChunk
				if json.Unmarshal([]byte(data), &ch) == nil {
					if ch.Error != nil {
						return &llmError{Code: "LLM_UPSTREAM_ERROR", Msg: sanitize(ch.Error.Message)}
					}
					for _, c := range ch.Choices {
						if c.Delta.Content != "" {
							if e := onDelta(c.Delta.Content); e != nil {
								return e
							}
						}
					}
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return &llmError{Code: "LLM_TIMEOUT", Msg: "大模型响应超时"}
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return &llmError{Code: "LLM_STREAM_BROKEN", Msg: "大模型流式响应中断"}
		}
	}
}

// sanitize 截断并去除可能回显的密钥片段
func sanitize(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 200 {
		s = s[:200]
	}
	return strings.ReplaceAll(s, "\n", " ")
}
