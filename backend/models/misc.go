package models

import "time"

// Profile 个人资料表
type Profile struct {
	ID     uint64 `json:"-" gorm:"primaryKey"`
	Avatar string `json:"avatar" gorm:"column:avatar;type:varchar(500)"`
	Bio    string `json:"bio" gorm:"column:bio;type:text"`
}

func (Profile) TableName() string { return "profile" }

// Message 留言表
type Message struct {
	ID               uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	DeviceID         string     `json:"-" gorm:"column:device_id;type:varchar(64);uniqueIndex"`
	Nickname         *string    `json:"nickname" gorm:"column:nickname;type:varchar(100)"`
	Content          string     `json:"content" gorm:"column:content;type:text"`
	LastModifiedDate *string    `json:"-" gorm:"column:last_modified_date;type:date"`
	IP               *string    `json:"-" gorm:"column:ip;type:varchar(45)"`
	CreatedAt        time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt        time.Time  `json:"updatedAt" gorm:"column:updated_at"`
}

func (Message) TableName() string { return "messages" }

// MessageListItem 留言列表项
type MessageListItem struct {
	ID        uint64 `json:"id"`
	Nickname  string `json:"nickname"`
	Content   string `json:"content"`
	IsOwn     bool   `json:"isOwn"`
	CanEdit   bool   `json:"canEdit"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// Changelog 更新日志表
type Changelog struct {
	ID        uint64    `json:"-" gorm:"primaryKey;autoIncrement"`
	Version   string    `json:"version" gorm:"column:version;type:varchar(20);uniqueIndex"`
	Date      string    `json:"date" gorm:"column:date;type:varchar(20)"`
	Logs      JSON      `json:"logs" gorm:"column:logs;type:json"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Changelog) TableName() string { return "changelogs" }
