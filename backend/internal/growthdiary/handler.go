package growthdiary

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatWill/agent-team-blog/backend/internal/backup"
	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetItems GET /api/growth-diary/items — 获取所有成长日记条目（按 happened_at DESC）
func GetItems(c *gin.Context) {
	var items []models.GrowthDiaryItem
	if err := db.DB.Order("happened_at DESC, id DESC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "查询成长日记失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// GetMonths GET /api/growth-diary/months — 获取有动态的年月列表（倒序）
func GetMonths(c *gin.Context) {
	var months []string
	err := db.DB.Model(&models.GrowthDiaryItem{}).
		Select("DISTINCT strftime('%Y-%m', happened_at) as month").
		Order("month DESC").
		Pluck("month", &months).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "查询月份列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"months": months})
}

// batchCreateItem 批量新增请求项
type batchCreateItem struct {
	Content    string   `json:"content"`
	Images     []string `json:"images"`
	Videos     []string `json:"videos"`
	HappenedAt string   `json:"happened_at"`
}

// batchUpdateItem 批量更新请求项
type batchUpdateItem struct {
	ID         int64    `json:"id"`
	Content    string   `json:"content"`
	Images     []string `json:"images"`
	Videos     []string `json:"videos"`
	HappenedAt string   `json:"happened_at"`
}

// batchRequest 批量编辑请求体
type batchRequest struct {
	Creates []batchCreateItem `json:"creates"`
	Updates []batchUpdateItem `json:"updates"`
	Deletes []int64           `json:"deletes"`
}

// BatchEditItems POST /api/growth-diary/items/batch — 批量编辑成长日记（事务）
func BatchEditItems(c *gin.Context) {
	var body batchRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	// 校验 creates
	for i, item := range body.Creates {
		if msg := validateItem(item.Content, item.Images, item.Videos, fmt.Sprintf("creates[%d]", i)); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": msg})
			return
		}
	}

	// 校验 updates
	for i, item := range body.Updates {
		if item.ID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": fmt.Sprintf("updates[%d]: id 无效", i)})
			return
		}
		if msg := validateItem(item.Content, item.Images, item.Videos, fmt.Sprintf("updates[%d]", i)); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": msg})
			return
		}
	}

	// 校验 deletes
	for i, id := range body.Deletes {
		if id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": fmt.Sprintf("deletes[%d]: id 无效", i)})
			return
		}
	}

	// 事务前备份
	if err := backup.BackupNow("pregrowth"); err != nil {
		log.Printf("⚠️ 事务前备份失败: %v", err)
	}

	// 事务执行
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除
		for _, id := range body.Deletes {
			result := tx.Where("id = ?", id).Delete(&models.GrowthDiaryItem{})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("删除失败: id=%d 不存在", id)
			}
		}

		// 2. 更新
		for _, u := range body.Updates {
			imagesJSON, _ := json.Marshal(u.Images)
			videosJSON, _ := json.Marshal(u.Videos)
			happenedAt := parseHappenedAt(u.HappenedAt)

			result := tx.Model(&models.GrowthDiaryItem{}).Where("id = ?", u.ID).Updates(map[string]interface{}{
				"content":     strings.TrimSpace(u.Content),
				"images":      string(imagesJSON),
				"videos":      string(videosJSON),
				"happened_at": happenedAt,
			})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("更新失败: id=%d 不存在", u.ID)
			}
		}

		// 3. 新增
		for _, cr := range body.Creates {
			imagesJSON, _ := json.Marshal(cr.Images)
			videosJSON, _ := json.Marshal(cr.Videos)
			happenedAt := parseHappenedAt(cr.HappenedAt)

			item := models.GrowthDiaryItem{
				Content:    strings.TrimSpace(cr.Content),
				Images:     string(imagesJSON),
				Videos:     string(videosJSON),
				HappenedAt: happenedAt,
			}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "不存在") {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": errMsg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "批量操作失败: " + errMsg})
		return
	}

	// 返回最新完整列表
	var items []models.GrowthDiaryItem
	db.DB.Order("happened_at DESC, id DESC").Find(&items)

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": gin.H{"items": items}})
}

// validateItem 校验单条记录
func validateItem(content string, images []string, videos []string, prefix string) string {
	ct := strings.TrimSpace(content)
	hasContent := ct != ""
	hasImages := len(images) > 0
	hasVideos := len(videos) > 0

	if !hasContent && !hasImages && !hasVideos {
		return fmt.Sprintf("%s: content、images、videos 至少一项非空", prefix)
	}

	// content 长度校验
	if utf8.RuneCountInString(ct) > 10000 {
		return fmt.Sprintf("%s: content 超过 10000 字符限制", prefix)
	}

	// images URL 校验
	for j, img := range images {
		if strings.TrimSpace(img) == "" {
			return fmt.Sprintf("%s: images[%d] 不能为空字符串", prefix, j)
		}
		if !isValidURL(img) {
			return fmt.Sprintf("%s: images[%d] 不是有效的 URL", prefix, j)
		}
	}

	// videos URL 校验
	for j, vid := range videos {
		if strings.TrimSpace(vid) == "" {
			return fmt.Sprintf("%s: videos[%d] 不能为空字符串", prefix, j)
		}
		if !isValidURL(vid) {
			return fmt.Sprintf("%s: videos[%d] 不是有效的 URL", prefix, j)
		}
	}

	return ""
}

// isValidURL 简单校验 URL 格式
func isValidURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

// parseHappenedAt 解析 happened_at 时间字符串，失败则返回当前时间
func parseHappenedAt(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Now()
	}
	// 尝试多种格式
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Now()
}
