package models

import "time"

// BudgetItem 装修成本预算项
type BudgetItem struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Category  string    `json:"category" gorm:"column:category;type:varchar(100);not null;default:''"`
	ItemName  string    `json:"itemName" gorm:"column:item_name;type:varchar(255);not null;default:''"`
	Amount    float64   `json:"amount" gorm:"column:amount;type:real;not null;default:0"`
	Remark    string    `json:"remark" gorm:"column:remark;type:varchar(255);not null;default:''"`
	SortOrder int       `json:"sortOrder" gorm:"column:sort_order;type:integer;not null;default:0"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (BudgetItem) TableName() string { return "budget_items" }
