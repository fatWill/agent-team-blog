package auth

import (
	"net/http"
	"os"

	"github.com/fatWill/agent-team-blog/backend/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// 管理后台账号信息，从环境变量读取
var (
	validUsername = getEnvDefault("ADMIN_USERNAME", "admin")
	validPassword = os.Getenv("ADMIN_PASSWORD")
)

// getEnvDefault 读取环境变量，未设置时返回默认值
func getEnvDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Login POST /api/auth/login
func Login(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.Username == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "请提供用户名和密码"})
		return
	}

	if body.Username != validUsername || body.Password != validPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "statusCode": 401, "statusMessage": "用户名或密码错误"})
		return
	}

	token, err := middleware.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "生成 Token 失败"})
		return
	}

	if err := middleware.SaveToken(token, body.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "保存 Token 失败"})
		return
	}

	// 写入 httpOnly cookie
	c.SetCookie(middleware.CookieName, token, int(middleware.TokenTTL.Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// AuthCheck GET /api/auth/check
func AuthCheck(c *gin.Context) {
	username := c.GetString("username")
	c.JSON(http.StatusOK, gin.H{"ok": true, "username": username})
}

// Logout POST /api/auth/logout
func Logout(c *gin.Context) {
	token, err := c.Cookie(middleware.CookieName)
	if err == nil && token != "" {
		_ = middleware.DeleteToken(token)
	}

	// 清除 cookie（设置 maxAge = -1）
	c.SetCookie(middleware.CookieName, "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
