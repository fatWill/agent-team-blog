package wechat

import (
	"net/http"

	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/gin-gonic/gin"
)

// SyncWechatHandler POST /api/admin/articles/:id/sync-wechat 手动触发同步
func SyncWechatHandler(c *gin.Context) {
	articleID := c.Param("id")
	if articleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "缺少文章 ID"})
		return
	}

	// 检查文章是否存在
	var count int64
	db.DB.Model(&models.Article{}).Where("id = ?", articleID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "文章不存在"})
		return
	}

	// 同步执行（管理 API 手动触发不走异步队列，直接返回结果）
	if err := SyncArticle(articleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":         true,
			"statusCode":    500,
			"statusMessage": "同步失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "message": "同步成功"})
}

// GetSyncLogsHandler GET /api/admin/articles/:id/wechat-sync-logs 查看同步日志
func GetSyncLogsHandler(c *gin.Context) {
	articleID := c.Param("id")
	if articleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "缺少文章 ID"})
		return
	}

	var logs []models.WechatSyncLog
	db.DB.Where("article_id = ?", articleID).Order("created_at DESC").Limit(50).Find(&logs)

	c.JSON(http.StatusOK, gin.H{"list": logs})
}
