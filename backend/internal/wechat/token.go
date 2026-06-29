package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fatWill/agent-team-blog/backend/pkg/rds"
)

// 微信公众号配置，从环境变量读取
var (
	appID     = os.Getenv("WECHAT_APP_ID")
	appSecret = os.Getenv("WECHAT_APP_SECRET")
)

// Redis 缓存 Key
const (
	keyAccessToken = "wechat:access_token"
	keyJSAPITicket = "wechat:jsapi_ticket"
	cacheTTL       = 7000 * time.Second // 略短于微信的 7200s，留 200s 安全余量
)

// tokenMu 防止并发获取 access_token
var tokenMu sync.Mutex

// TokenInfo access_token 状态信息
type TokenInfo struct {
	Token       string
	Remaining   time.Duration
	LastRefresh time.Time
}

// 微信接口响应结构
type wxTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type wxTicketResp struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
}

type wxIPResp struct {
	IPList  []string `json:"ip_list"`
	ErrCode int      `json:"errcode"`
	ErrMsg  string   `json:"errmsg"`
}

// GetAccessToken 获取 access_token（带缓存、加锁、重试）
func GetAccessToken() (string, error) {
	ctx := context.Background()

	// 先尝试从缓存获取（无锁快速路径）
	token, err := rds.RDB.Get(ctx, keyAccessToken).Result()
	if err == nil && token != "" {
		return token, nil
	}

	// 加锁，防止并发请求微信 API
	tokenMu.Lock()
	defer tokenMu.Unlock()

	// 双重检查：拿到锁后再查一次缓存
	token, err = rds.RDB.Get(ctx, keyAccessToken).Result()
	if err == nil && token != "" {
		return token, nil
	}

	// 重试 3 次
	var lastErr error
	for i := 0; i < 3; i++ {
		token, lastErr = fetchAccessToken()
		if lastErr == nil {
			// 写入缓存
			rds.RDB.Set(ctx, keyAccessToken, token, cacheTTL)
			return token, nil
		}
		log.Printf("[wechat] 获取 access_token 第 %d 次失败: %v", i+1, lastErr)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return "", fmt.Errorf("获取 access_token 失败（已重试3次）: %w", lastErr)
}

// ForceRefreshToken 强制刷新 access_token（用于 40001 错误后重试）
func ForceRefreshToken() (string, error) {
	ctx := context.Background()
	rds.RDB.Del(ctx, keyAccessToken)

	tokenMu.Lock()
	defer tokenMu.Unlock()

	token, err := fetchAccessToken()
	if err != nil {
		return "", err
	}
	rds.RDB.Set(ctx, keyAccessToken, token, cacheTTL)
	return token, nil
}

// GetTokenInfo 获取 token 状态信息
func GetTokenInfo() TokenInfo {
	ctx := context.Background()
	token, _ := rds.RDB.Get(ctx, keyAccessToken).Result()
	ttl, _ := rds.RDB.TTL(ctx, keyAccessToken).Result()

	return TokenInfo{
		Token:       token,
		Remaining:   ttl,
		LastRefresh: time.Now().Add(-(cacheTTL - ttl)),
	}
}

// fetchAccessToken 直接调用微信接口获取 access_token
func fetchAccessToken() (string, error) {
	reqURL := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		appID, appSecret,
	)

	resp, err := http.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("请求微信 access_token 接口失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取微信响应失败: %w", err)
	}

	var result wxTokenResp
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析微信响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信接口错误: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	return result.AccessToken, nil
}

// getJSAPITicket 获取 jsapi_ticket
func getJSAPITicket() (string, error) {
	ctx := context.Background()

	ticket, err := rds.RDB.Get(ctx, keyJSAPITicket).Result()
	if err == nil && ticket != "" {
		return ticket, nil
	}

	accessToken, err := GetAccessToken()
	if err != nil {
		return "", err
	}

	reqURL := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi",
		accessToken,
	)

	resp, err := http.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("请求微信 jsapi_ticket 接口失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取微信响应失败: %w", err)
	}

	var result wxTicketResp
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析微信响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		if result.ErrCode == 40001 || result.ErrCode == 42001 {
			rds.RDB.Del(ctx, keyAccessToken)
		}
		return "", fmt.Errorf("微信接口错误: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	rds.RDB.Set(ctx, keyJSAPITicket, result.Ticket, cacheTTL)
	return result.Ticket, nil
}

// getWxAPIIPs 获取微信 API 服务器 IP 列表
func getWxAPIIPs(token string) ([]string, error) {
	reqURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/get_api_domain_ip?access_token=%s", token)
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result wxIPResp
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.ErrCode != 0 {
		return nil, fmt.Errorf("errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}
	return result.IPList, nil
}
