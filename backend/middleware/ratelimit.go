package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// rateLimitEntry IP 限频记录
type rateLimitEntry struct {
	count   int
	resetAt time.Time
}

var (
	ipMap   = make(map[string]*rateLimitEntry)
	ipMutex sync.Mutex
)

func init() {
	// 每 5 分钟清理过期条目
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			ipMutex.Lock()
			now := time.Now()
			for ip, entry := range ipMap {
				if now.After(entry.resetAt) {
					delete(ipMap, ip)
				}
			}
			ipMutex.Unlock()
		}
	}()
}

// checkRateLimit 检查 IP 是否超出限频
func checkRateLimit(ip string, maxRequests int) bool {
	ipMutex.Lock()
	defer ipMutex.Unlock()

	now := time.Now()
	entry, exists := ipMap[ip]

	if !exists || now.After(entry.resetAt) {
		ipMap[ip] = &rateLimitEntry{count: 1, resetAt: now.Add(time.Minute)}
		return true
	}

	if entry.count < maxRequests {
		entry.count++
		return true
	}

	return false
}

// GetClientIP 从请求中提取客户端 IP
func GetClientIP(c *gin.Context) string {
	// 优先读 x-forwarded-for
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	// 其次读 x-real-ip
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}
	return c.ClientIP()
}

// RateLimit 限流中间件
func RateLimit(maxRequests int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := GetClientIP(c)
		if !checkRateLimit(ip, maxRequests) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":         true,
				"statusCode":    429,
				"statusMessage": "操作太频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
