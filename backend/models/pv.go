package models

import "time"

// PageView 页面访问记录表
type PageView struct {
	ID         uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Path       string    `json:"path" gorm:"column:path;type:text;not null"`
	DeviceID   string    `json:"device_id" gorm:"column:device_id;type:text;not null"`
	IP         string    `json:"ip" gorm:"column:ip;type:text;not null"`
	UserAgent  string    `json:"-" gorm:"column:user_agent;type:text"`
	DeviceType string    `json:"device_type" gorm:"column:device_type;type:text"`
	Browser    string    `json:"browser" gorm:"column:browser;type:text"`
	OS         string    `json:"os" gorm:"column:os;type:text"`
	Referer    string    `json:"referer" gorm:"column:referer;type:text"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
}

func (PageView) TableName() string { return "page_views" }

// PVTrendItem PV/UV 趋势数据项
type PVTrendItem struct {
	Date string `json:"date"`
	PV   int64  `json:"pv"`
	UV   int64  `json:"uv"`
}

// TopPageItem Top 页面数据项
type TopPageItem struct {
	Path string `json:"path"`
	PV   int64  `json:"pv"`
	UV   int64  `json:"uv"`
}

// PVOverview 统计概览
type PVOverview struct {
	TodayPV int64 `json:"today_pv"`
	TodayUV int64 `json:"today_uv"`
	TotalPV int64 `json:"total_pv"`
	TotalUV int64 `json:"total_uv"`
}
