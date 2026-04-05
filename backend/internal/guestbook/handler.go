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

// generateNickname 根据 deviceId 自动生成昵称：用户 + deviceId 后4位
func generateNickname(deviceID string) string {
	suffix := deviceID
	if len(suffix) > 4 {
		suffix = suffix[len(suffix)-4:]
	}
	return "用户" + suffix
}

// GetMessages GET /api/messages
func GetMessages(c *gin.Context) {
	deviceID := strings.TrimSpace(c.Query("deviceId"))

	var messages []models.Message
	db.DB.Order("created_at DESC").Find(&messages)

	list := make([]models.MessageListItem, 0, len(messages))
	for _, m := range messages {
		isOwn := deviceID != "" && m.DeviceID == deviceID
		nickname := generateNickname(m.DeviceID)

		list = append(list, models.MessageListItem{
			ID:        m.ID,
			Nickname:  nickname,
			Content:   m.Content,
			IsOwn:     isOwn,
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

	// deviceId 唯一性校验：每个设备只能留言一次
	var count int64
	db.DB.Model(&models.Message{}).Where("device_id = ?", deviceID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "每个设备只能留言一次"})
		return
	}

	// 自动生成昵称
	nickname := generateNickname(deviceID)
	clientIP := middleware.GetClientIP(c)

	msg := models.Message{
		DeviceID: deviceID,
		Nickname: &nickname,
		Content:  content,
		IP:       &clientIP,
	}

	if err := db.DB.Create(&msg).Error; err != nil {
		// 并发场景下唯一索引冲突
		if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "Duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "每个设备只能留言一次"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": fmt.Sprintf("创建留言失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        msg.ID,
		"nickname":  nickname,
		"content":   msg.Content,
		"isOwn":     true,
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
