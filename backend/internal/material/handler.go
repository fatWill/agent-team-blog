package material

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/internal/upload"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/gin-gonic/gin"
)

// GetMaterials GET /api/materials — 获取所有材料清单条目
func GetMaterials(c *gin.Context) {
	var items []models.MaterialItemListItem

	err := db.DB.Model(&models.MaterialItem{}).
		Select("id, title, tags, attachments, sort_order, created_at, updated_at").
		Order("sort_order ASC, id DESC").
		Find(&items).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "查询材料清单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list": items})
}

// GetMaterial GET /api/materials/:id — 获取单个条目详情
func GetMaterial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的 ID"})
		return
	}

	var item models.MaterialItem
	if err := db.DB.Where("id = ?", id).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "材料条目不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          item.ID,
		"title":       item.Title,
		"tags":        item.Tags,
		"attachments": item.Attachments,
		"sortOrder":   item.SortOrder,
		"createdAt":   item.CreatedAt,
		"updatedAt":   item.UpdatedAt,
	})
}

// CreateMaterial POST /api/materials — 新建条目
func CreateMaterial(c *gin.Context) {
	var body struct {
		Title       string                     `json:"title"`
		Tags        []string                   `json:"tags"`
		Attachments []models.MaterialAttachment `json:"attachments"`
		SortOrder   int                        `json:"sortOrder"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	if strings.TrimSpace(body.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "标题不能为空"})
		return
	}

	tagsJSON, _ := json.Marshal(body.Tags)
	attachmentsJSON, _ := json.Marshal(body.Attachments)

	item := models.MaterialItem{
		Title:       strings.TrimSpace(body.Title),
		Tags:        models.JSON(tagsJSON),
		Attachments: models.JSON(attachmentsJSON),
		SortOrder:   body.SortOrder,
	}

	if err := db.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "创建材料条目失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          item.ID,
		"title":       item.Title,
		"tags":        item.Tags,
		"attachments": item.Attachments,
		"sortOrder":   item.SortOrder,
		"createdAt":   item.CreatedAt,
		"updatedAt":   item.UpdatedAt,
	})
}

// UpdateMaterial PUT /api/materials/:id — 更新条目（支持部分更新）
func UpdateMaterial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的 ID"})
		return
	}

	var body struct {
		Title       *string                     `json:"title,omitempty"`
		Tags        *[]string                   `json:"tags,omitempty"`
		Attachments *[]models.MaterialAttachment `json:"attachments,omitempty"`
		SortOrder   *int                        `json:"sortOrder,omitempty"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	if body.Title == nil && body.Tags == nil && body.Attachments == nil && body.SortOrder == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "至少需要提供一个更新字段"})
		return
	}

	updates := map[string]interface{}{}
	if body.Title != nil {
		t := strings.TrimSpace(*body.Title)
		if t == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "标题不能为空"})
			return
		}
		updates["title"] = t
	}
	if body.Tags != nil {
		tagsJSON, _ := json.Marshal(*body.Tags)
		updates["tags"] = string(tagsJSON)
	}
	if body.Attachments != nil {
		attachmentsJSON, _ := json.Marshal(*body.Attachments)
		updates["attachments"] = string(attachmentsJSON)
	}
	if body.SortOrder != nil {
		updates["sort_order"] = *body.SortOrder
	}

	result := db.DB.Model(&models.MaterialItem{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "更新材料条目失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "材料条目不存在"})
		return
	}

	var item models.MaterialItem
	db.DB.Where("id = ?", id).First(&item)

	c.JSON(http.StatusOK, gin.H{
		"id":          item.ID,
		"title":       item.Title,
		"tags":        item.Tags,
		"attachments": item.Attachments,
		"sortOrder":   item.SortOrder,
		"createdAt":   item.CreatedAt,
		"updatedAt":   item.UpdatedAt,
	})
}

// DeleteMaterial DELETE /api/materials/:id — 删除条目（同时异步删除 COS 附件）
func DeleteMaterial(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的 ID"})
		return
	}

	// 先查出附件列表
	var item models.MaterialItem
	if err := db.DB.Where("id = ?", id).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "材料条目不存在"})
		return
	}

	// 删除数据库记录
	result := db.DB.Where("id = ?", id).Delete(&models.MaterialItem{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "删除材料条目失败"})
		return
	}

	// 异步删除 COS 附件
	go func() {
		var attachments []models.MaterialAttachment
		if err := json.Unmarshal([]byte(item.Attachments), &attachments); err != nil {
			log.Printf("⚠️  解析附件 JSON 失败: %v", err)
			return
		}
		urls := make([]string, 0, len(attachments))
		for _, a := range attachments {
			if a.URL != "" {
				urls = append(urls, a.URL)
			}
			if a.ThumbnailURL != "" {
				urls = append(urls, a.ThumbnailURL)
			}
		}
		if len(urls) > 0 {
			upload.BatchDeleteFromCOS(urls)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// UpdateMaterialSort PUT /api/materials/:id/sort — 更新排序
func UpdateMaterialSort(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的 ID"})
		return
	}

	var body struct {
		SortOrder int `json:"sortOrder"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	result := db.DB.Model(&models.MaterialItem{}).Where("id = ?", id).Update("sort_order", body.SortOrder)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "更新排序失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "材料条目不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
