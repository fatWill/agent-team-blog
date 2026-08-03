package models

import "time"

// GrowthDiaryItem 阳阳的成长日记条目
type GrowthDiaryItem struct {
	ID         int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Content    string    `json:"content" gorm:"column:content;type:text;not null;default:''"`
	Images     string    `json:"images" gorm:"column:images;type:text;not null;default:'[]'"`
	Videos     string    `json:"videos" gorm:"column:videos;type:text;not null;default:'[]'"`
	HappenedAt time.Time `json:"happenedAt" gorm:"column:happened_at;not null"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (GrowthDiaryItem) TableName() string { return "growth_diary_items" }
