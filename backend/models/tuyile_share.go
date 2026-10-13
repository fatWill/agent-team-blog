package models

import "time"

// TuyileShareMessage 分享快照中的单条消息（仅 user / assistant 纯文本）
type TuyileShareMessage struct {
	Role    string `json:"role"`    // user / assistant
	Content string `json:"content"` // 纯文本（前端以安全模式渲染 Markdown）
}

// TuyileShare 图一乐对话分享快照（只读，不可基于分享继续对话）
type TuyileShare struct {
	ID        uint64               `json:"-" gorm:"primaryKey;autoIncrement"`
	Hash      string               `json:"hash" gorm:"column:hash"`                         // 12 位 base62 随机串（crypto/rand）
	Messages  []TuyileShareMessage `json:"messages" gorm:"column:messages;serializer:json"` // 对话快照
	CreatedAt time.Time            `json:"createdAt" gorm:"column:created_at"`
}

// TableName 表名
func (TuyileShare) TableName() string { return "tuyile_shares" }
