package pitfall

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/fatWill/agent-team-blog/backend/internal/backup"
	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 分类白名单
var allowedCategories = map[string]bool{
	"局改":   true,
	"硬装":   true,
	"家私":   true,
	"家电":   true,
	"全屋智能": true,
	"其他":   true,
}

// GetItems GET /api/renovation/pitfall/items — 获取所有踩坑日记条目
func GetItems(c *gin.Context) {
	var items []models.PitfallItem
	if err := db.DB.Order("sort_order ASC, id ASC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "查询踩坑日记失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// batchCreateItem 批量新增请求项
type batchCreateItem struct {
	Category  string `json:"category"`
	Pitfall   string `json:"pitfall"`
	Solution  string `json:"solution"`
	Remark    string `json:"remark"`
	SortOrder int    `json:"sort_order"`
}

// batchUpdateItem 批量更新请求项
type batchUpdateItem struct {
	ID        int64  `json:"id"`
	Category  string `json:"category"`
	Pitfall   string `json:"pitfall"`
	Solution  string `json:"solution"`
	Remark    string `json:"remark"`
	SortOrder int    `json:"sort_order"`
}

// batchDeleteItem 批量删除请求项
type batchDeleteItem struct {
	ID int64 `json:"id"`
}

// batchRequest 批量编辑请求体
type batchRequest struct {
	Creates []batchCreateItem `json:"creates"`
	Updates []batchUpdateItem `json:"updates"`
	Deletes []batchDeleteItem `json:"deletes"`
}

// BatchEditItems POST /api/renovation/pitfall/items/batch — 批量编辑踩坑日记（事务）
func BatchEditItems(c *gin.Context) {
	var body batchRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	// 校验 creates
	for i, item := range body.Creates {
		if msg := validateItem(item.Category, item.Pitfall, item.Solution, item.Remark, fmt.Sprintf("creates[%d]", i)); msg != "" {
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
		if msg := validateItem(item.Category, item.Pitfall, item.Solution, item.Remark, fmt.Sprintf("updates[%d]", i)); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": msg})
			return
		}
	}

	// 校验 deletes
	for i, item := range body.Deletes {
		if item.ID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": fmt.Sprintf("deletes[%d]: id 无效", i)})
			return
		}
	}

	// 事务前备份（失败不阻塞业务）
	if err := backup.BackupNow("prepitfall"); err != nil {
		log.Printf("⚠️ 事务前备份失败: %v", err)
	}

	// 事务执行
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除
		for _, d := range body.Deletes {
			result := tx.Where("id = ?", d.ID).Delete(&models.PitfallItem{})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("删除失败: id=%d 不存在", d.ID)
			}
		}

		// 2. 更新
		for _, u := range body.Updates {
			result := tx.Model(&models.PitfallItem{}).Where("id = ?", u.ID).Updates(map[string]interface{}{
				"category":   strings.TrimSpace(u.Category),
				"pitfall":    strings.TrimSpace(u.Pitfall),
				"solution":   strings.TrimSpace(u.Solution),
				"remark":     strings.TrimSpace(u.Remark),
				"sort_order": u.SortOrder,
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
			item := models.PitfallItem{
				Category:  strings.TrimSpace(cr.Category),
				Pitfall:   strings.TrimSpace(cr.Pitfall),
				Solution:  strings.TrimSpace(cr.Solution),
				Remark:    strings.TrimSpace(cr.Remark),
				SortOrder: cr.SortOrder,
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
	var items []models.PitfallItem
	db.DB.Order("sort_order ASC, id ASC").Find(&items)

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": gin.H{"items": items}})
}

// validateItem 校验单条记录
func validateItem(category, pitfall, solution, remark, prefix string) string {
	cat := strings.TrimSpace(category)
	if cat == "" {
		return fmt.Sprintf("%s: category 不能为空", prefix)
	}
	if !allowedCategories[cat] {
		return fmt.Sprintf("%s: category 不在允许范围内，允许值: 局改/硬装/家私/家电/全屋智能/其他", prefix)
	}

	pit := strings.TrimSpace(pitfall)
	sol := strings.TrimSpace(solution)
	if pit == "" && sol == "" {
		return fmt.Sprintf("%s: pitfall 和 solution 不能同时为空", prefix)
	}

	// 长度校验
	if utf8.RuneCountInString(cat) > 50 {
		return fmt.Sprintf("%s: category 超过 50 字符限制", prefix)
	}
	if utf8.RuneCountInString(pit) > 5000 {
		return fmt.Sprintf("%s: pitfall 超过 5000 字符限制", prefix)
	}
	if utf8.RuneCountInString(sol) > 5000 {
		return fmt.Sprintf("%s: solution 超过 5000 字符限制", prefix)
	}
	if utf8.RuneCountInString(strings.TrimSpace(remark)) > 5000 {
		return fmt.Sprintf("%s: remark 超过 5000 字符限制", prefix)
	}

	return ""
}
