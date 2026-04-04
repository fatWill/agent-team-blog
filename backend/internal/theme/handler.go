package theme

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fatWill/agent-team-blog/backend/internal/upload"
	"github.com/fatWill/agent-team-blog/backend/pkg/rds"
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
	t, err := rds.RDB.Get(ctx, key).Result()
	if err != nil || t != "dark" {
		t = "light"
	}

	c.JSON(http.StatusOK, gin.H{"theme": t})
}

// SaveTheme POST /api/theme
func SaveTheme(c *gin.Context) {
	var body struct {
		Theme string `json:"theme"`
	}
	c.ShouldBindJSON(&body)

	t := "light"
	if body.Theme == "dark" {
		t = "dark"
	}

	uid, _ := c.Cookie("blog_uid")
	if uid == "" {
		uid = fmt.Sprintf("%d-%s", time.Now().UnixMilli(), upload.RandomString(6))
		c.SetCookie("blog_uid", uid, 60*60*24*30, "/", "", false, true)
	}

	key := fmt.Sprintf("theme:%s", uid)

	ctx := context.Background()
	rds.RDB.Set(ctx, key, t, 30*24*time.Hour)

	c.SetCookie("color-mode", t, 60*60*24*30, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{"ok": true, "theme": t})
}
