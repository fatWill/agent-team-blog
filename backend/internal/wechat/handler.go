package wechat

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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

	ticket, err := getJSAPITicket()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         true,
			"statusCode":    500,
			"statusMessage": "获取微信签名失败",
		})
		return
	}

	nonceStr := randomString(16)
	timestamp := time.Now().Unix()

	signStr := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ticket, nonceStr, timestamp, url)

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

// GetAccessTokenStatus GET /api/admin/wechat/access-token-status
func GetAccessTokenStatus(c *gin.Context) {
	info := GetTokenInfo()
	c.JSON(http.StatusOK, gin.H{
		"hasToken":     info.Token != "",
		"remainingSec": int(info.Remaining.Seconds()),
		"lastRefresh":  info.LastRefresh,
	})
}

// GetServerIP GET /api/admin/wechat/server-ip
func GetServerIP(c *gin.Context) {
	// 获取本机出口 IP
	resp, err := http.Get("https://ifconfig.me/ip")
	var outboundIP string
	if err == nil {
		defer resp.Body.Close()
		buf := make([]byte, 64)
		n, _ := resp.Body.Read(buf)
		outboundIP = string(buf[:n])
	} else {
		outboundIP = "获取失败: " + err.Error()
	}

	// 获取微信 API 服务器 IP（可选，需 access_token）
	var wxIPs []string
	token, tokenErr := GetAccessToken()
	if tokenErr == nil && token != "" {
		wxIPs, _ = getWxAPIIPs(token)
	}

	c.JSON(http.StatusOK, gin.H{
		"outboundIP":   outboundIP,
		"wxAPIIPs":     wxIPs,
		"hint":         "请将 outboundIP 添加到微信公众号后台 → 基本配置 → IP 白名单",
		"mpConsoleURL": "https://mp.weixin.qq.com/",
	})
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
