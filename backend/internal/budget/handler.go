package budget

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/fatWill/agent-team-blog/backend/internal/backup"
	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetItems GET /api/renovation/budget/items — 获取所有预算项
func GetItems(c *gin.Context) {
	var items []models.BudgetItem
	if err := db.DB.Order("sort_order ASC, id ASC").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "查询预算项失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// batchCreateItem 批量新增请求项
type batchCreateItem struct {
	Category  string  `json:"category"`
	ItemName  string  `json:"item_name"`
	Amount    float64 `json:"amount"`
	Actual    float64 `json:"actual"`
	Remark    string  `json:"remark"`
	SortOrder int     `json:"sort_order"`
}

// batchUpdateItem 批量更新请求项
type batchUpdateItem struct {
	ID        int64   `json:"id"`
	Category  string  `json:"category"`
	ItemName  string  `json:"item_name"`
	Amount    float64 `json:"amount"`
	Actual    float64 `json:"actual"`
	Remark    string  `json:"remark"`
	SortOrder int     `json:"sort_order"`
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

// BatchEditItems POST /api/renovation/budget/items/batch — 批量编辑预算项（事务）
func BatchEditItems(c *gin.Context) {
	var body batchRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
		return
	}

	// 校验 creates
	for i, item := range body.Creates {
		if strings.TrimSpace(item.Category) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": fmt.Sprintf("creates[%d]: category 不能为空", i)})
			return
		}
		if strings.TrimSpace(item.ItemName) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": fmt.Sprintf("creates[%d]: item_name 不能为空", i)})
			return
		}
		if item.Amount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": fmt.Sprintf("creates[%d]: amount 不能为负数", i)})
			return
		}
		if item.Actual < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": fmt.Sprintf("creates[%d]: actual 不能为负数", i)})
			return
		}
	}

	// 校验 updates
	for i, item := range body.Updates {
		if item.ID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": fmt.Sprintf("updates[%d]: id 无效", i)})
			return
		}
		if strings.TrimSpace(item.Category) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": fmt.Sprintf("updates[%d]: category 不能为空", i)})
			return
		}
		if strings.TrimSpace(item.ItemName) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": fmt.Sprintf("updates[%d]: item_name 不能为空", i)})
			return
		}
		if item.Amount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": fmt.Sprintf("updates[%d]: amount 不能为负数", i)})
			return
		}
		if item.Actual < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": fmt.Sprintf("updates[%d]: actual 不能为负数", i)})
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
	if err := backup.BackupNow("prebudget"); err != nil {
		log.Printf("⚠️ 事务前备份失败: %v", err)
	}

	// 事务执行
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除
		for _, d := range body.Deletes {
			result := tx.Where("id = ?", d.ID).Delete(&models.BudgetItem{})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("删除失败: id=%d 不存在", d.ID)
			}
		}

		// 2. 更新
		for _, u := range body.Updates {
			result := tx.Model(&models.BudgetItem{}).Where("id = ?", u.ID).Updates(map[string]interface{}{
				"category":   strings.TrimSpace(u.Category),
				"item_name":  strings.TrimSpace(u.ItemName),
				"amount":     roundAmount(u.Amount),
				"actual":     roundAmount(u.Actual),
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
			item := models.BudgetItem{
				Category:  strings.TrimSpace(cr.Category),
				ItemName:  strings.TrimSpace(cr.ItemName),
				Amount:    roundAmount(cr.Amount),
				Actual:    roundAmount(cr.Actual),
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
	var items []models.BudgetItem
	db.DB.Order("sort_order ASC, id ASC").Find(&items)

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": gin.H{"items": items}})
}

// roundAmount 保留两位小数
func roundAmount(v float64) float64 {
	return math.Round(v*100) / 100
}
