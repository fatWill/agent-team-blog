package models

import "time"

// PitfallItem 装修踩坑日记条目
type PitfallItem struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Category  string    `json:"category" gorm:"column:category;type:varchar(50);not null;default:''"`
	Pitfall   string    `json:"pitfall" gorm:"column:pitfall;type:text;not null;default:''"`
	Solution  string    `json:"solution" gorm:"column:solution;type:text;not null;default:''"`
	Remark    string    `json:"remark" gorm:"column:remark;type:text;not null;default:''"`
	SortOrder int       `json:"sortOrder" gorm:"column:sort_order;type:integer;not null;default:0"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (PitfallItem) TableName() string { return "pitfall_items" }
