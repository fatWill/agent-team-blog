package models

import "time"

// RenovationArticle 装修文章表
type RenovationArticle struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Title     string    `json:"title" gorm:"column:title;type:text;not null;default:''"`
	Content   string    `json:"content" gorm:"column:content;type:text;not null;default:''"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (RenovationArticle) TableName() string { return "renovation_articles" }

// RenovationArticleListItem 装修文章列表项（不含 content）
type RenovationArticleListItem struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
