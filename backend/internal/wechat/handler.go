package wechat

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/fatWill/agent-team-blog/backend/pkg/rds"
	"github.com/gin-gonic/gin"
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

// GetJSSDKConfig GET /api/wechat/jssdk-config
func GetJSSDKConfig(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":         true,
			"statusCode":    400,
			"statusMessage": "缺少 url 参数",
		})
		return
	}

	if appID == "" || appSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         true,
			"statusCode":    500,
			"statusMessage": "微信公众号未配置",
		})
		return
	}

	// 获取 jsapi_ticket（内部会自动获取/缓存 access_token）
	ticket, err := getJSAPITicket()
	if err != nil {
		log.Printf("获取 jsapi_ticket 失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         true,
			"statusCode":    500,
			"statusMessage": "获取微信签名失败",
		})
		return
	}

	// 生成签名参数
	nonceStr := randomString(16)
	timestamp := time.Now().Unix()

	// 按字典序拼接签名字符串
	signStr := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ticket, nonceStr, timestamp, url)

	// SHA1 签名
	h := sha1.New()
	h.Write([]byte(signStr))
	signature := fmt.Sprintf("%x", h.Sum(nil))

	c.JSON(http.StatusOK, gin.H{
		"appId":     appID,
		"timestamp": timestamp,
		"nonceStr":  nonceStr,
		"signature": signature,
	})
}

// getAccessToken 获取 access_token，优先从 Redis 缓存读取
func getAccessToken() (string, error) {
	ctx := context.Background()

	// 尝试从缓存获取
	token, err := rds.RDB.Get(ctx, keyAccessToken).Result()
	if err == nil && token != "" {
		return token, nil
	}

	// 缓存未命中，调用微信接口
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

	// 写入 Redis 缓存
	rds.RDB.Set(ctx, keyAccessToken, result.AccessToken, cacheTTL)

	return result.AccessToken, nil
}

// getJSAPITicket 获取 jsapi_ticket，优先从 Redis 缓存读取
func getJSAPITicket() (string, error) {
	ctx := context.Background()

	// 尝试从缓存获取
	ticket, err := rds.RDB.Get(ctx, keyJSAPITicket).Result()
	if err == nil && ticket != "" {
		return ticket, nil
	}

	// 缓存未命中，先获取 access_token
	accessToken, err := getAccessToken()
	if err != nil {
		return "", err
	}

	// 调用微信接口获取 jsapi_ticket
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
		// 如果 ticket 获取失败且错误码是 token 过期，清除 access_token 缓存
		if result.ErrCode == 40001 || result.ErrCode == 42001 {
			rds.RDB.Del(ctx, keyAccessToken)
		}
		return "", fmt.Errorf("微信接口错误: errcode=%d, errmsg=%s", result.ErrCode, result.ErrMsg)
	}

	// 写入 Redis 缓存
	rds.RDB.Set(ctx, keyJSAPITicket, result.Ticket, cacheTTL)

	return result.Ticket, nil
}

// randomString 生成指定长度的随机字符串
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
