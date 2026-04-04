package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/fatWill/agent-team-blog/backend/utils"
	"github.com/gin-gonic/gin"
)

const (
	// TokenTTL Token 有效期 72 小时
	TokenTTL = 72 * time.Hour
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

// SaveToken 保存 Token 到 Redis（72h 过期）
func SaveToken(token, username string) error {
	return utils.RDB.Set(context.Background(), TokenKeyPrefix+token, username, TokenTTL).Err()
}

// DeleteToken 删除 Token
func DeleteToken(token string) error {
	return utils.RDB.Del(context.Background(), TokenKeyPrefix+token).Err()
}

// VerifyToken 验证 Token 并滚动续期
func VerifyToken(token string) (string, error) {
	ctx := context.Background()
	username, err := utils.RDB.Get(ctx, TokenKeyPrefix+token).Result()
	if err != nil {
		return "", err
	}
	// 滚动续期
	utils.RDB.Expire(ctx, TokenKeyPrefix+token, TokenTTL)
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
