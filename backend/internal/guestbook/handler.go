package guestbook

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/fatWill/agent-team-blog/backend/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// getTodayCST 获取 UTC+8 的今天日期字符串
func getTodayCST() string {
	loc := time.FixedZone("CST", 8*3600)
	return time.Now().In(loc).Format("2006-01-02")
}

// GetMessages GET /api/messages
func GetMessages(c *gin.Context) {
	deviceID := strings.TrimSpace(c.Query("deviceId"))

	var messages []models.Message
	db.DB.Order("created_at DESC").Find(&messages)

	today := getTodayCST()

	list := make([]models.MessageListItem, 0, len(messages))
	for _, m := range messages {
		isOwn := deviceID != "" && m.DeviceID == deviceID

		var lastModDate string
		if m.LastModifiedDate != nil {
			lastModDate = *m.LastModifiedDate
			if len(lastModDate) > 10 {
				lastModDate = lastModDate[:10]
			}
		}

		canEdit := isOwn && lastModDate != today

		nickname := "匿名"
		if m.Nickname != nil && *m.Nickname != "" {
			nickname = *m.Nickname
		}

		list = append(list, models.MessageListItem{
			ID:        m.ID,
			Nickname:  nickname,
			Content:   m.Content,
			IsOwn:     isOwn,
			CanEdit:   canEdit,
			CreatedAt: m.CreatedAt.Format(time.RFC3339),
			UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{"list": list})
}

// CreateMessage POST /api/messages
func CreateMessage(c *gin.Context) {
	var body struct {
		DeviceID string `json:"deviceId"`
		Nickname string `json:"nickname"`
		Content  string `json:"content"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	deviceID := strings.TrimSpace(body.DeviceID)
	if len(deviceID) < 8 || len(deviceID) > 64 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的设备 ID"})
		return
	}

	content := strings.TrimSpace(body.Content)
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "留言内容不能为空"})
		return
	}

	var nickname *string
	if n := strings.TrimSpace(body.Nickname); n != "" {
		nickname = &n
	}

	clientIP := middleware.GetClientIP(c)
	today := getTodayCST()

	var existing models.Message
	err := db.DB.Where("device_id = ?", deviceID).First(&existing).Error

	if err == nil {
		var lastModDate string
		if existing.LastModifiedDate != nil {
			lastModDate = *existing.LastModifiedDate
			if len(lastModDate) > 10 {
				lastModDate = lastModDate[:10]
			}
		}

		if lastModDate == today {
			c.JSON(http.StatusForbidden, gin.H{"error": true, "statusCode": 403, "statusMessage": "今天已经修改过了，明天再来吧"})
			return
		}

		db.DB.Model(&existing).Updates(map[string]interface{}{
			"nickname":           nickname,
			"content":            content,
			"last_modified_date": today,
			"ip":                 clientIP,
		})

		db.DB.Where("id = ?", existing.ID).First(&existing)

		displayNickname := "匿名"
		if existing.Nickname != nil && *existing.Nickname != "" {
			displayNickname = *existing.Nickname
		}

		c.JSON(http.StatusOK, gin.H{
			"id":        existing.ID,
			"nickname":  displayNickname,
			"content":   existing.Content,
			"isOwn":     true,
			"canEdit":   false,
			"createdAt": existing.CreatedAt.Format(time.RFC3339),
			"updatedAt": existing.UpdatedAt.Format(time.RFC3339),
		})
		return
	}

	msg := models.Message{
		DeviceID: deviceID,
		Nickname: nickname,
		Content:  content,
		IP:       &clientIP,
	}

	if err := db.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": fmt.Sprintf("创建留言失败: %v", err)})
		return
	}

	displayNickname := "匿名"
	if msg.Nickname != nil && *msg.Nickname != "" {
		displayNickname = *msg.Nickname
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        msg.ID,
		"nickname":  displayNickname,
		"content":   msg.Content,
		"isOwn":     true,
		"canEdit":   true,
		"createdAt": msg.CreatedAt.Format(time.RFC3339),
		"updatedAt": msg.UpdatedAt.Format(time.RFC3339),
	})
}

// UpdateMessage PUT /api/messages/:id
func UpdateMessage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的留言 ID"})
		return
	}

	var body struct {
		DeviceID string `json:"deviceId"`
		Nickname string `json:"nickname"`
		Content  string `json:"content"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	deviceID := strings.TrimSpace(body.DeviceID)
	if len(deviceID) < 8 || len(deviceID) > 64 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的设备 ID"})
		return
	}

	content := strings.TrimSpace(body.Content)
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "留言内容不能为空"})
		return
	}

	var nickname *string
	if n := strings.TrimSpace(body.Nickname); n != "" {
		nickname = &n
	}

	var msg models.Message
	if err := db.DB.Where("id = ?", id).First(&msg).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "留言不存在"})
		return
	}

	if msg.DeviceID != deviceID {
		c.JSON(http.StatusForbidden, gin.H{"error": true, "statusCode": 403, "statusMessage": "无权修改此留言"})
		return
	}

	today := getTodayCST()

	if msg.LastModifiedDate != nil {
		lastModDate := *msg.LastModifiedDate
		if len(lastModDate) > 10 {
			lastModDate = lastModDate[:10]
		}
		if lastModDate == today {
			c.JSON(http.StatusForbidden, gin.H{"error": true, "statusCode": 403, "statusMessage": "今天已经修改过了，明天再来吧"})
			return
		}
	}

	db.DB.Model(&msg).Updates(map[string]interface{}{
		"nickname":           nickname,
		"content":            content,
		"last_modified_date": today,
	})

	db.DB.Where("id = ?", id).First(&msg)

	displayNickname := "匿名"
	if msg.Nickname != nil && *msg.Nickname != "" {
		displayNickname = *msg.Nickname
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        msg.ID,
		"nickname":  displayNickname,
		"content":   msg.Content,
		"isOwn":     true,
		"canEdit":   false,
		"createdAt": msg.CreatedAt.Format(time.RFC3339),
		"updatedAt": msg.UpdatedAt.Format(time.RFC3339),
	})
}

// DeleteMessage DELETE /api/messages/:id（需鉴权）
func DeleteMessage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的留言 ID"})
		return
	}

	result := db.DB.Where("id = ?", id).Delete(&models.Message{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除留言失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "留言不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
