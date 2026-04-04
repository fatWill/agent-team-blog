package models

import "time"

// Album 相册表
type Album struct {
	ID           uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string     `json:"name" gorm:"column:name;type:varchar(100)"`
	Description  *string    `json:"description" gorm:"column:description;type:varchar(200)"`
	CoverURL     *string    `json:"coverUrl" gorm:"column:cover_url;type:varchar(500)"`
	PasswordHash *string    `json:"-" gorm:"column:password_hash;type:varchar(255)"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"column:updated_at"`
}

func (Album) TableName() string { return "albums" }

// AlbumListItem 相册列表项（含 photo_count 和 hasPassword）
type AlbumListItem struct {
	ID          uint64  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CoverURL    *string `json:"coverUrl"`
	PhotoCount  int     `json:"photoCount"`
	HasPassword bool    `json:"hasPassword"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

// Photo 照片表
type Photo struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	AlbumID      uint64    `json:"albumId" gorm:"column:album_id"`
	URL          string    `json:"url" gorm:"column:url;type:varchar(500)"`
	Caption      *string   `json:"caption" gorm:"column:caption;type:varchar(200)"`
	PasswordHash *string   `json:"-" gorm:"column:password_hash;type:varchar(255)"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Photo) TableName() string { return "photos" }

// PhotoListItem 照片列表项
type PhotoListItem struct {
	ID          uint64  `json:"id"`
	AlbumID     uint64  `json:"albumId"`
	URL         string  `json:"url"`
	Caption     *string `json:"caption"`
	HasPassword bool    `json:"hasPassword"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

// PhotoLike 照片点赞记录
type PhotoLike struct {
	ID       uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	PhotoID  uint64    `json:"photoId" gorm:"column:photo_id"`
	DeviceID string    `json:"deviceId" gorm:"column:device_id;type:varchar(64)"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (PhotoLike) TableName() string { return "photo_likes" }

// PhotoDislike 照片踩记录
type PhotoDislike struct {
	ID       uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	PhotoID  uint64    `json:"photoId" gorm:"column:photo_id"`
	DeviceID string    `json:"deviceId" gorm:"column:device_id;type:varchar(64)"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (PhotoDislike) TableName() string { return "photo_dislikes" }
