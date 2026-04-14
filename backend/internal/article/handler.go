package article

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/fatWill/agent-team-blog/backend/pkg/nginx"
	"github.com/fatWill/agent-team-blog/backend/pkg/rds"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetArticles GET /api/articles
func GetArticles(c *gin.Context) {
	title := c.Query("title")
	search := c.Query("search")

	var articles []struct {
		models.Article
	}

	q := db.DB.Select("id, title, summary, cover_image, like_count, views, created_at, updated_at")
	if title != "" {
		q = q.Where("title LIKE ?", "%"+title+"%")
	}
	if search != "" {
		q = q.Where("title LIKE ?", "%"+search+"%")
	}
	q = q.Order("created_at DESC")

	if err := q.Table("articles").Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "查询文章列表失败"})
		return
	}

	list := make([]models.ArticleListItem, 0, len(articles))
	for _, a := range articles {
		list = append(list, models.ArticleListItem{
			ID:         a.ID,
			Title:      a.Title,
			Summary:    a.Summary,
			CoverImage: a.CoverImage,
			LikeCount:  a.LikeCount,
			Views:      a.Views,
			CreatedAt:  a.CreatedAt,
			UpdatedAt:  a.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"list": list})
}

// GetArticle GET /api/articles/:id
func GetArticle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "缺少文章 ID"})
		return
	}

	var article models.Article
	if err := db.DB.Where("id = ?", id).First(&article).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         article.ID,
		"title":      article.Title,
		"summary":    article.Summary,
		"coverImage": article.CoverImage,
		"content":    article.Content,
		"likeCount":  article.LikeCount,
		"views":      article.Views,
		"createdAt":  article.CreatedAt,
		"updatedAt":  article.UpdatedAt,
	})
}

// CreateArticle POST /api/articles
func CreateArticle(c *gin.Context) {
	var body struct {
		Title      string          `json:"title"`
		Summary    string          `json:"summary"`
		CoverImage string          `json:"coverImage"`
		Content    json.RawMessage `json:"content"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	if strings.TrimSpace(body.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "文章标题不能为空"})
		return
	}

	if len(body.Content) == 0 || string(body.Content) == "null" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "文章内容不能为空，且必须为 Tiptap JSON 对象"})
		return
	}

	a := models.Article{
		ID:         uuid.New().String(),
		Title:      strings.TrimSpace(body.Title),
		Summary:    strings.TrimSpace(body.Summary),
		CoverImage: strings.TrimSpace(body.CoverImage),
		Content:    models.JSON(body.Content),
	}

	if err := db.DB.Create(&a).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "创建文章失败"})
		return
	}

	// 异步清除 Nginx 首页缓存
	nginx.PurgeCache([]string{"/"})

	c.JSON(http.StatusOK, gin.H{
		"id":         a.ID,
		"title":      a.Title,
		"summary":    a.Summary,
		"coverImage": a.CoverImage,
		"content":    a.Content,
		"likeCount":  a.LikeCount,
		"views":      a.Views,
		"createdAt":  a.CreatedAt,
		"updatedAt":  a.UpdatedAt,
	})
}

// UpdateArticle PUT /api/articles/:id
func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "缺少文章 ID"})
		return
	}

	var body struct {
		Title      *string          `json:"title"`
		Summary    *string          `json:"summary"`
		CoverImage *string          `json:"coverImage"`
		Content    *json.RawMessage `json:"content"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	if body.Title == nil && body.Summary == nil && body.CoverImage == nil && body.Content == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "至少需要提供 title、summary、coverImage 或 content 中的一个字段"})
		return
	}

	updates := map[string]interface{}{}
	if body.Title != nil {
		t := strings.TrimSpace(*body.Title)
		if t != "" {
			updates["title"] = t
		}
	}
	if body.Summary != nil {
		updates["summary"] = strings.TrimSpace(*body.Summary)
	}
	if body.CoverImage != nil {
		updates["cover_image"] = strings.TrimSpace(*body.CoverImage)
	}
	if body.Content != nil {
		updates["content"] = string(*body.Content)
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "没有有效的更新字段"})
		return
	}

	result := db.DB.Model(&models.Article{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "更新文章失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "文章不存在"})
		return
	}

	// 返回更新后的文章
	var a models.Article
	db.DB.Where("id = ?", id).First(&a)

	// 异步清除 Nginx 缓存：文章详情页 + 首页
	nginx.PurgeCache([]string{"/", "/articles/" + id})

	c.JSON(http.StatusOK, gin.H{
		"id":         a.ID,
		"title":      a.Title,
		"summary":    a.Summary,
		"coverImage": a.CoverImage,
		"content":    a.Content,
		"likeCount":  a.LikeCount,
		"views":      a.Views,
		"createdAt":  a.CreatedAt,
		"updatedAt":  a.UpdatedAt,
	})
}

// DeleteArticle DELETE /api/articles/:id
func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "缺少文章 ID"})
		return
	}

	result := db.DB.Where("id = ?", id).Delete(&models.Article{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "删除文章失败"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "文章不存在"})
		return
	}

	// 异步清除 Nginx 缓存：文章详情页 + 首页
	nginx.PurgeCache([]string{"/", "/articles/" + id})

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// LikeArticle POST /api/articles/:id/like
func LikeArticle(c *gin.Context) {
	articleID := c.Param("id")
	if articleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的文章 ID"})
		return
	}

	var body struct {
		DeviceID string `json:"deviceId"`
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

	// 检查文章是否存在
	var count int64
	db.DB.Model(&models.Article{}).Where("id = ?", articleID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "文章不存在"})
		return
	}

	// 检查是否已点赞
	var likeCount int64
	db.DB.Model(&models.ArticleLike{}).Where("article_id = ? AND device_id = ?", articleID, deviceID).Count(&likeCount)

	var liked bool
	if likeCount > 0 {
		// 已点赞 → 取消
		db.DB.Where("article_id = ? AND device_id = ?", articleID, deviceID).Delete(&models.ArticleLike{})
		db.DB.Model(&models.Article{}).Where("id = ?", articleID).
			Update("like_count", gorm.Expr("MAX(like_count - 1, 0)"))
		liked = false
	} else {
		// 未点赞 → 点赞
		db.DB.Create(&models.ArticleLike{ArticleID: articleID, DeviceID: deviceID})
		db.DB.Model(&models.Article{}).Where("id = ?", articleID).
			Update("like_count", gorm.Expr("like_count + 1"))
		liked = true
	}

	// 查询最新计数
	var a models.Article
	db.DB.Select("like_count").Where("id = ?", articleID).First(&a)

	c.JSON(http.StatusOK, gin.H{"liked": liked, "likeCount": a.LikeCount})
}

// GetArticleLikeStatus GET /api/articles/:id/like-status
func GetArticleLikeStatus(c *gin.Context) {
	articleID := c.Param("id")
	if articleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的文章 ID"})
		return
	}

	deviceID := strings.TrimSpace(c.Query("deviceId"))

	var a models.Article
	if err := db.DB.Select("like_count").Where("id = ?", articleID).First(&a).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "文章不存在"})
		return
	}

	liked := false
	if deviceID != "" {
		var cnt int64
		db.DB.Model(&models.ArticleLike{}).Where("article_id = ? AND device_id = ?", articleID, deviceID).Count(&cnt)
		liked = cnt > 0
	}

	c.JSON(http.StatusOK, gin.H{"liked": liked, "likeCount": a.LikeCount})
}

// GetArticleLikeStatusBatch GET /api/articles/like-status-batch
func GetArticleLikeStatusBatch(c *gin.Context) {
	deviceID := strings.TrimSpace(c.Query("deviceId"))
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "deviceId 不能为空"})
		return
	}

	idsParam := strings.TrimSpace(c.Query("ids"))
	if idsParam == "" {
		c.JSON(http.StatusOK, gin.H{"likedIds": []string{}})
		return
	}

	// 解析并过滤有效的文章 ID，最多 100 个
	parts := strings.Split(idsParam, ",")
	ids := make([]string, 0, len(parts))
	for _, s := range parts {
		s = strings.TrimSpace(s)
		if s != "" {
			ids = append(ids, s)
		}
		if len(ids) >= 100 {
			break
		}
	}

	if len(ids) == 0 {
		c.JSON(http.StatusOK, gin.H{"likedIds": []string{}})
		return
	}

	var likes []models.ArticleLike
	db.DB.Select("article_id").Where("article_id IN ? AND device_id = ?", ids, deviceID).Find(&likes)

	likedIds := make([]string, 0, len(likes))
	for _, l := range likes {
		likedIds = append(likedIds, l.ArticleID)
	}

	c.JSON(http.StatusOK, gin.H{"likedIds": likedIds})
}

// ViewArticle POST /api/articles/:id/view — 记录文章阅读量
func ViewArticle(c *gin.Context) {
	articleID := c.Param("id")
	if articleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的文章 ID"})
		return
	}

	var body struct {
		DeviceID string `json:"deviceId"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	deviceID := strings.TrimSpace(body.DeviceID)
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "deviceId 不能为空"})
		return
	}

	// 检查文章是否存在
	var article models.Article
	if err := db.DB.Select("id, views").Where("id = ?", articleID).First(&article).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "文章不存在"})
		return
	}

	// 防刷：同一设备 1 小时内重复访问同一篇文章不重复计数
	redisKey := fmt.Sprintf("view:%s:%s", articleID, deviceID)
	ctx := context.Background()
	exists, err := rds.RDB.Exists(ctx, redisKey).Result()
	if err != nil {
		// Redis 不可用时降级，直接计数
	} else if exists > 0 {
		// 1h 内重复访问，返回当前阅读量但不 +1
		c.JSON(http.StatusOK, gin.H{"views": article.Views})
		return
	}

	// 设置防刷 key，1h 过期
	_ = rds.RDB.Set(ctx, redisKey, "1", time.Hour).Err()

	// 阅读量 +1
	db.DB.Model(&models.Article{}).Where("id = ?", articleID).
		Update("views", gorm.Expr("views + 1"))

	// 查询最新阅读量
	db.DB.Select("views").Where("id = ?", articleID).First(&article)

	c.JSON(http.StatusOK, gin.H{"views": article.Views})
}

// GetRandomArticle GET /api/articles/random — 随机返回一篇文章
func GetRandomArticle(c *gin.Context) {
	var article models.Article
	if err := db.DB.Select("id").Order("RANDOM()").First(&article).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "statusCode": 404, "statusMessage": "暂无文章"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": article.ID})
}