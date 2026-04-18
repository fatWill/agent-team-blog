package renovation

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/gin-gonic/gin"
)

// GetArticles GET /api/renovation/articles — 获取所有装修文章列表（不含 content）
func GetArticles(c *gin.Context) {
	var articles []models.RenovationArticleListItem

	err := db.DB.Model(&models.RenovationArticle{}).
		Select("id, title, created_at, updated_at").
		Order("id DESC").
		Find(&articles).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "查询装修文章列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list": articles})
}

// GetArticle GET /api/renovation/articles/:id — 获取装修文章详情（含 content）
func GetArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的文章 ID"})
		return
	}

	var article models.RenovationArticle
	if err := db.DB.Where("id = ?", id).First(&article).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "装修文章不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        article.ID,
		"title":     article.Title,
		"content":   article.Content,
		"createdAt": article.CreatedAt,
		"updatedAt": article.UpdatedAt,
	})
}

// CreateArticle POST /api/renovation/articles — 创建装修文章（需鉴权）
func CreateArticle(c *gin.Context) {
	var body struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	if strings.TrimSpace(body.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "文章标题不能为空"})
		return
	}
	if strings.TrimSpace(body.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "文章内容不能为空"})
		return
	}

	article := models.RenovationArticle{
		Title:   strings.TrimSpace(body.Title),
		Content: strings.TrimSpace(body.Content),
	}

	if err := db.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "创建装修文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        article.ID,
		"title":     article.Title,
		"content":   article.Content,
		"createdAt": article.CreatedAt,
		"updatedAt": article.UpdatedAt,
	})
}

// UpdateArticle PUT /api/renovation/articles/:id — 更新装修文章（需鉴权）
func UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的文章 ID"})
		return
	}

	var body struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	if body.Title == nil && body.Content == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "至少需要提供 title 或 content"})
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
	if body.Content != nil {
		updates["content"] = strings.TrimSpace(*body.Content)
	}

	result := db.DB.Model(&models.RenovationArticle{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "更新装修文章失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "装修文章不存在"})
		return
	}

	// 返回更新后的文章
	var article models.RenovationArticle
	db.DB.Where("id = ?", id).First(&article)

	c.JSON(http.StatusOK, gin.H{
		"id":        article.ID,
		"title":     article.Title,
		"content":   article.Content,
		"createdAt": article.CreatedAt,
		"updatedAt": article.UpdatedAt,
	})
}

// DeleteArticle DELETE /api/renovation/articles/:id — 删除装修文章（需鉴权）
func DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的文章 ID"})
		return
	}

	result := db.DB.Where("id = ?", id).Delete(&models.RenovationArticle{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "删除装修文章失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "装修文章不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
