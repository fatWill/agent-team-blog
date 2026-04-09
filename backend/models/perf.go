package models

import "time"

// PerfMetric 前端性能指标记录表
type PerfMetric struct {
	ID         uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Page       string    `json:"page" gorm:"column:page;type:text;not null"`
	FCP        *float64  `json:"fcp" gorm:"column:fcp"`
	LCP        *float64  `json:"lcp" gorm:"column:lcp"`
	TTFB       *float64  `json:"ttfb" gorm:"column:ttfb"`
	DOMReady   *float64  `json:"domReady" gorm:"column:dom_ready"`
	LoadTime   *float64  `json:"loadTime" gorm:"column:load_time"`
	JSLoad     *float64  `json:"jsLoad" gorm:"column:js_load"`
	DeviceType string    `json:"deviceType" gorm:"column:device_type;type:text"`
	UA         string    `json:"-" gorm:"column:ua;type:text"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (PerfMetric) TableName() string { return "perf_metrics" }

// PerfOverview 性能概览响应
type PerfOverview struct {
	TotalCount  int64   `json:"totalCount"`
	AvgFcp      float64 `json:"avgFcp"`
	AvgLcp      float64 `json:"avgLcp"`
	AvgTtfb     float64 `json:"avgTtfb"`
	AvgLoadTime float64 `json:"avgLoadTime"`
	P75Lcp      float64 `json:"p75Lcp"`
	P75Fcp      float64 `json:"p75Fcp"`
}

// PerfPageItem 各页面性能排行项
type PerfPageItem struct {
	Page        string  `json:"page"`
	Count       int64   `json:"count"`
	AvgFcp      float64 `json:"avgFcp"`
	AvgLcp      float64 `json:"avgLcp"`
	AvgTtfb     float64 `json:"avgTtfb"`
	AvgLoadTime float64 `json:"avgLoadTime"`
}

// PerfTrendItem 性能趋势项
type PerfTrendItem struct {
	Date    string  `json:"date"`
	AvgFcp  float64 `json:"avgFcp"`
	AvgLcp  float64 `json:"avgLcp"`
	AvgTtfb float64 `json:"avgTtfb"`
	Count   int64   `json:"count"`
}

// PerfLogItem 性能日志列表项
type PerfLogItem struct {
	ID         uint64   `json:"id"`
	Page       string   `json:"page"`
	FCP        *float64 `json:"fcp"`
	LCP        *float64 `json:"lcp"`
	TTFB       *float64 `json:"ttfb"`
	DOMReady   *float64 `json:"domReady"`
	LoadTime   *float64 `json:"loadTime"`
	JSLoad     *float64 `json:"jsLoad"`
	DeviceType string   `json:"deviceType"`
	CreatedAt  string   `json:"createdAt"`
}
