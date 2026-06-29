package models

import "time"

// Article 文章表
type Article struct {
	ID                  string     `json:"id" gorm:"column:id;primaryKey;type:varchar(36)"`
	Title               string     `json:"title" gorm:"column:title;type:varchar(255)"`
	Summary             string     `json:"summary" gorm:"column:summary;type:varchar(500)"`
	Content             JSON       `json:"content" gorm:"column:content;type:text"`
	CoverImage          string     `json:"coverImage" gorm:"column:cover_image;type:varchar(500)"`
	LikeCount           int        `json:"likeCount" gorm:"column:like_count;type:int;default:0"`
	Views               int        `json:"views" gorm:"column:views;type:int;default:0"`
	WechatDraftMediaID  *string    `json:"wechatDraftMediaId,omitempty" gorm:"column:wechat_draft_media_id"`
	WechatSyncedAt      *time.Time `json:"wechatSyncedAt,omitempty" gorm:"column:wechat_synced_at"`
	WechatSyncStatus    string     `json:"wechatSyncStatus" gorm:"column:wechat_sync_status;default:pending"`
	WechatSyncError     *string    `json:"wechatSyncError,omitempty" gorm:"column:wechat_sync_error"`
	WechatAutoSync      int        `json:"wechatAutoSync" gorm:"column:wechat_auto_sync;default:1"`
	CreatedAt           time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt           time.Time  `json:"updatedAt" gorm:"column:updated_at"`
}

func (Article) TableName() string { return "articles" }

// ArticleListItem 文章列表项（不含 content）
type ArticleListItem struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	CoverImage string    `json:"coverImage"`
	LikeCount  int       `json:"likeCount"`
	Views      int       `json:"views"`
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

// WechatSyncLog 微信同步日志
type WechatSyncLog struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ArticleID string    `json:"articleId" gorm:"column:article_id;type:varchar(36)"`
	Action    string    `json:"action" gorm:"column:action;type:varchar(32)"`
	Status    string    `json:"status" gorm:"column:status;type:varchar(16)"`
	Error     *string   `json:"error,omitempty" gorm:"column:error;type:text"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (WechatSyncLog) TableName() string { return "wechat_sync_logs" }