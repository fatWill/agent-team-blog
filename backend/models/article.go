package models

import "time"

// Article 文章表
type Article struct {
	ID        string    `json:"id" gorm:"column:id;primaryKey;type:varchar(36)"`
	Title     string    `json:"title" gorm:"column:title;type:varchar(255)"`
	Summary   string    `json:"summary" gorm:"column:summary;type:varchar(500)"`
	Content   JSON      `json:"content" gorm:"column:content;type:json"`
	CoverImage string   `json:"coverImage" gorm:"column:cover_image;type:varchar(500)"`
	LikeCount  int      `json:"likeCount" gorm:"column:like_count;type:int;default:0"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Article) TableName() string { return "articles" }

// ArticleListItem 文章列表项（不含 content）
type ArticleListItem struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	CoverImage string    `json:"coverImage"`
	LikeCount  int       `json:"likeCount"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// ArticleLike 文章点赞记录
type ArticleLike struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ArticleID string    `json:"articleId" gorm:"column:article_id;type:varchar(36)"`
	DeviceID  string    `json:"deviceId" gorm:"column:device_id;type:varchar(64)"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (ArticleLike) TableName() string { return "article_likes" }
