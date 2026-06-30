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

// syncLogItem 全局同步日志返回结构
type syncLogItem struct {
	ID           uint64  `json:"id"`
	ArticleID    string  `json:"articleId"`
	ArticleTitle string  `json:"articleTitle"`
	Action       string  `json:"action"`
	Status       string  `json:"status"`
	Error        *string `json:"error,omitempty"`
	CreatedAt    string  `json:"createdAt"`
}

// ListAllSyncLogs GET /api/admin/wechat/sync-logs 全局同步日志查询
func ListAllSyncLogs(c *gin.Context) {
	// 分页参数
	page := 1
	pageSize := 20
	if v := c.Query("page"); v != "" {
		if n := parseInt(v); n > 0 {
			page = n
		}
	}
	if v := c.Query("pageSize"); v != "" {
		if n := parseInt(v); n > 0 && n <= 100 {
			pageSize = n
		}
	}

	status := c.Query("status")
	action := c.Query("action")
	articleID := c.Query("articleId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// 构建查询
	q := db.DB.Table("wechat_sync_logs AS l").
		Select("l.id, l.article_id, l.action, l.status, l.error, l.created_at, COALESCE(a.title, '(已删除)') AS article_title").
		Joins("LEFT JOIN articles a ON a.id = l.article_id")

	if status != "" {
		q = q.Where("l.status = ?", status)
	}
	if action != "" {
		q = q.Where("l.action = ?", action)
	}
	if articleID != "" {
		q = q.Where("l.article_id = ?", articleID)
	}
	if startDate != "" {
		q = q.Where("l.created_at >= ?", startDate)
	}
	if endDate != "" {
		q = q.Where("l.created_at <= ?", endDate+" 23:59:59")
	}

	// 查询总数
	var total int64
	q.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	type rawRow struct {
		ID           uint64  `gorm:"column:id"`
		ArticleID    string  `gorm:"column:article_id"`
		ArticleTitle string  `gorm:"column:article_title"`
		Action       string  `gorm:"column:action"`
		Status       string  `gorm:"column:status"`
		Error        *string `gorm:"column:error"`
		CreatedAt    string  `gorm:"column:created_at"`
	}
	var rows []rawRow
	if err := q.Order("l.created_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "查询同步日志失败"})
		return
	}

	list := make([]syncLogItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, syncLogItem{
			ID:           r.ID,
			ArticleID:    r.ArticleID,
			ArticleTitle: r.ArticleTitle,
			Action:       r.Action,
			Status:       r.Status,
			Error:        r.Error,
			CreatedAt:    r.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// parseInt 安全解析整数
func parseInt(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			return 0
		}
	}
	return n
}
