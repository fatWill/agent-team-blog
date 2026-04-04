package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fatWill/agent-team-blog/backend/utils"
	"github.com/gin-gonic/gin"
)

// GetTheme GET /api/theme
func GetTheme(c *gin.Context) {
	uid, _ := c.Cookie("blog_uid")
	if uid == "" {
		uid = c.ClientIP()
	}
	if uid == "" {
		uid = "anonymous"
	}

	key := fmt.Sprintf("theme:%s", uid)

	ctx := context.Background()
	theme, err := utils.RDB.Get(ctx, key).Result()
	if err != nil || theme != "dark" {
		theme = "light"
	}

	c.JSON(http.StatusOK, gin.H{"theme": theme})
}

// SaveTheme POST /api/theme
func SaveTheme(c *gin.Context) {
	var body struct {
		Theme string `json:"theme"`
	}
	c.ShouldBindJSON(&body)

	theme := "light"
	if body.Theme == "dark" {
		theme = "dark"
	}

	// 确保 uid cookie 存在
	uid, _ := c.Cookie("blog_uid")
	if uid == "" {
		uid = fmt.Sprintf("%d-%s", time.Now().UnixMilli(), randomString(6))
		c.SetCookie("blog_uid", uid, 60*60*24*30, "/", "", false, true)
	}

	key := fmt.Sprintf("theme:%s", uid)

	ctx := context.Background()
	// 主题偏好存 30 天
	utils.RDB.Set(ctx, key, theme, 30*24*time.Hour)

	// 同时写一个客户端可读的 cookie
	c.SetCookie("color-mode", theme, 60*60*24*30, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"ok": true, "theme": theme})
}
