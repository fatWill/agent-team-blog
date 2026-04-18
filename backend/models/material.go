package models

import "time"

// MaterialAttachment 单个附件（序列化为 JSON 存储）
type MaterialAttachment struct {
	ID           string `json:"id"`                      // UUID v4
	Type         string `json:"type"`                    // "pdf" | "image" | "video"
	URL          string `json:"url"`                     // COS 文件 URL
	Name         string `json:"name"`                    // 原始文件名
	Size         int64  `json:"size,omitempty"`          // 文件大小（字节）
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`  // 视频封面图 URL
	Duration     int    `json:"duration,omitempty"`      // 视频时长（秒）
}

// MaterialItem 材料清单条目
type MaterialItem struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" gorm:"column:title;type:text;not null;default:''"`
	Tags        JSON      `json:"tags" gorm:"column:tags;type:text"`        // []string JSON 数组
	Attachments JSON      `json:"attachments" gorm:"column:attachments;type:text"` // []MaterialAttachment JSON 数组
	SortOrder   int       `json:"sortOrder" gorm:"column:sort_order;type:integer;not null;default:0"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (MaterialItem) TableName() string { return "material_items" }

// MaterialItemListItem 列表项
type MaterialItemListItem struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Tags        JSON      `json:"tags"`
	Attachments JSON      `json:"attachments"`
	SortOrder   int       `json:"sortOrder"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
