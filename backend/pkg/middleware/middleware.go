package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fatWill/agent-team-blog/backend/pkg/rds"
	"github.com/gin-gonic/gin"
)

// ==================== 鉴权中间件 ====================

const (
	// TokenTTL Token 有效期 30 天
	TokenTTL = 30 * 24 * time.Hour
	// TokenKeyPrefix Redis Key 前缀
	TokenKeyPrefix = "auth_token:"
	// CookieName cookie 名称
	CookieName = "auth_token"
)

// GenerateToken 生成 64 位随机 Token
func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SaveToken 保存 Token 到 Redis
func SaveToken(token, username string) error {
	return rds.RDB.Set(context.Background(), TokenKeyPrefix+token, username, TokenTTL).Err()
}

// DeleteToken 删除 Token
func DeleteToken(token string) error {
	return rds.RDB.Del(context.Background(), TokenKeyPrefix+token).Err()
}

// VerifyToken 验证 Token 并滚动续期
func VerifyToken(token string) (string, error) {
	ctx := context.Background()
	username, err := rds.RDB.Get(ctx, TokenKeyPrefix+token).Result()
	if err != nil {
		return "", err
	}
	// 滚动续期
	rds.RDB.Expire(ctx, TokenKeyPrefix+token, TokenTTL)
	return username, nil
}

// AuthRequired 鉴权中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(CookieName)
		if err != nil || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": true, "statusCode": 401, "statusMessage": "未登录，请先登录"})
			c.Abort()
			return
		}

		username, err := VerifyToken(token)
		if err != nil || username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": true, "statusCode": 401, "statusMessage": "Token 已过期或无效，请重新登录"})
			c.Abort()
			return
		}

		// 续期 cookie
		c.SetCookie(CookieName, token, int(TokenTTL.Seconds()), "/", "", false, true)
		c.Set("username", username)
		c.Next()
	}
}

// ==================== IP 限频中间件 ====================

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
