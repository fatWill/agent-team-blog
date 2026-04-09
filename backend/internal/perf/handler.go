package perf

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fatWill/agent-team-blog/backend/models"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/fatWill/agent-team-blog/backend/pkg/middleware"
	"github.com/fatWill/agent-team-blog/backend/pkg/rds"
	"github.com/gin-gonic/gin"
)

// parseDays 解析 days 查询参数，默认 7，最大 30
func parseDays(c *gin.Context) int {
	days := 7
	if d, err := strconv.Atoi(c.Query("days")); err == nil && d > 0 && d <= 30 {
		days = d
	}
	return days
}

// startDateCST 根据天数计算 CST 起始日期字符串
func startDateCST(days int) string {
	loc := time.FixedZone("CST", 8*3600)
	return time.Now().In(loc).AddDate(0, 0, -days+1).Format("2006-01-02")
}

// Report POST /api/perf/report — 上报性能指标（无需鉴权）
func Report(c *gin.Context) {
	var body struct {
		Page       string   `json:"page"`
		FCP        *float64 `json:"fcp"`
		LCP        *float64 `json:"lcp"`
		TTFB       *float64 `json:"ttfb"`
		DOMReady   *float64 `json:"domReady"`
		LoadTime   *float64 `json:"loadTime"`
		JSLoad     *float64 `json:"jsLoad"`
		DeviceType string   `json:"deviceType"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "参数解析失败"})
		return
	}

	page := strings.TrimSpace(body.Page)
	if page == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "page 不能为空"})
		return
	}

	// 限流：同一 IP 同一页面 60 秒内只允许上报一次
	ip := middleware.GetClientIP(c)
	redisKey := fmt.Sprintf("perf:%s:%s", ip, page)
	ctx := context.Background()

	exists, err := rds.RDB.Exists(ctx, redisKey).Result()
	if err != nil {
		log.Printf("⚠️ Perf 限流 Redis 查询失败（降级写入）: %v", err)
	} else if exists > 0 {
		c.JSON(http.StatusOK, gin.H{"ok": true})
		return
	}

	if err := rds.RDB.Set(ctx, redisKey, "1", 60*time.Second).Err(); err != nil {
		log.Printf("⚠️ Perf 限流 Redis 写入失败（降级继续）: %v", err)
	}

	// 截断 UA 到 200 字符
	ua := c.GetHeader("User-Agent")
	if len(ua) > 200 {
		ua = ua[:200]
	}

	record := models.PerfMetric{
		Page:       page,
		FCP:        body.FCP,
		LCP:        body.LCP,
		TTFB:       body.TTFB,
		DOMReady:   body.DOMReady,
		LoadTime:   body.LoadTime,
		JSLoad:     body.JSLoad,
		DeviceType: strings.TrimSpace(body.DeviceType),
		UA:         ua,
	}

	if err := db.DB.Create(&record).Error; err != nil {
		log.Printf("❌ Perf 指标写入失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": "写入失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GetOverview GET /api/perf/overview?days=7 — 性能概览（需鉴权）
func GetOverview(c *gin.Context) {
	days := parseDays(c)
	start := startDateCST(days)

	var result struct {
		TotalCount  int64   `gorm:"column:total_count"`
		AvgFcp      float64 `gorm:"column:avg_fcp"`
		AvgLcp      float64 `gorm:"column:avg_lcp"`
		AvgTtfb     float64 `gorm:"column:avg_ttfb"`
		AvgLoadTime float64 `gorm:"column:avg_load_time"`
	}

	db.DB.Raw(`
		SELECT
			COUNT(*) AS total_count,
			COALESCE(AVG(fcp), 0) AS avg_fcp,
			COALESCE(AVG(lcp), 0) AS avg_lcp,
			COALESCE(AVG(ttfb), 0) AS avg_ttfb,
			COALESCE(AVG(load_time), 0) AS avg_load_time
		FROM perf_metrics
		WHERE date(created_at, '+8 hours') >= ?
	`, start).Scan(&result)

	// 计算 P75（LCP 和 FCP）
	// SQLite 不支持 PERCENTILE，使用排序 + LIMIT/OFFSET 方式
	p75Lcp := calcP75(start, "lcp")
	p75Fcp := calcP75(start, "fcp")

	c.JSON(http.StatusOK, gin.H{
		"totalCount":  result.TotalCount,
		"avgFcp":      math.Round(result.AvgFcp*100) / 100,
		"avgLcp":      math.Round(result.AvgLcp*100) / 100,
		"avgTtfb":     math.Round(result.AvgTtfb*100) / 100,
		"avgLoadTime": math.Round(result.AvgLoadTime*100) / 100,
		"p75Lcp":      p75Lcp,
		"p75Fcp":      p75Fcp,
	})
}

// calcP75 计算指定指标的第 75 百分位数
func calcP75(startDate string, column string) float64 {
	var count int64
	db.DB.Raw(fmt.Sprintf(`
		SELECT COUNT(*) FROM perf_metrics
		WHERE date(created_at, '+8 hours') >= ? AND %s IS NOT NULL
	`, column), startDate).Scan(&count)

	if count == 0 {
		return 0
	}

	// P75 位置 = ceil(count * 0.75) - 1（0-indexed offset）
	offset := int(math.Ceil(float64(count)*0.75)) - 1
	if offset < 0 {
		offset = 0
	}

	var val float64
	db.DB.Raw(fmt.Sprintf(`
		SELECT %s FROM perf_metrics
		WHERE date(created_at, '+8 hours') >= ? AND %s IS NOT NULL
		ORDER BY %s ASC
		LIMIT 1 OFFSET ?
	`, column, column, column), startDate, offset).Scan(&val)

	return math.Round(val*100) / 100
}

// GetPages GET /api/perf/pages?days=7 — 各页面性能排行（需鉴权）
func GetPages(c *gin.Context) {
	days := parseDays(c)
	start := startDateCST(days)

	var results []models.PerfPageItem
	db.DB.Raw(`
		SELECT
			page,
			COUNT(*) AS count,
			COALESCE(AVG(fcp), 0) AS avg_fcp,
			COALESCE(AVG(lcp), 0) AS avg_lcp,
			COALESCE(AVG(ttfb), 0) AS avg_ttfb,
			COALESCE(AVG(load_time), 0) AS avg_load_time
		FROM perf_metrics
		WHERE date(created_at, '+8 hours') >= ?
		GROUP BY page
		ORDER BY avg_lcp DESC
		LIMIT 20
	`, start).Scan(&results)

	// 保留两位小数
	for i := range results {
		results[i].AvgFcp = math.Round(results[i].AvgFcp*100) / 100
		results[i].AvgLcp = math.Round(results[i].AvgLcp*100) / 100
		results[i].AvgTtfb = math.Round(results[i].AvgTtfb*100) / 100
		results[i].AvgLoadTime = math.Round(results[i].AvgLoadTime*100) / 100
	}

	c.JSON(http.StatusOK, gin.H{"list": results})
}

// GetTrend GET /api/perf/trend?days=7&page= — 性能趋势（需鉴权）
func GetTrend(c *gin.Context) {
	days := parseDays(c)
	start := startDateCST(days)
	filterPage := strings.TrimSpace(c.Query("page"))

	var results []models.PerfTrendItem

	if filterPage != "" {
		db.DB.Raw(`
			SELECT
				date(created_at, '+8 hours') AS date,
				COALESCE(AVG(fcp), 0) AS avg_fcp,
				COALESCE(AVG(lcp), 0) AS avg_lcp,
				COALESCE(AVG(ttfb), 0) AS avg_ttfb,
				COUNT(*) AS count
			FROM perf_metrics
			WHERE date(created_at, '+8 hours') >= ? AND page = ?
			GROUP BY date(created_at, '+8 hours')
			ORDER BY date ASC
		`, start, filterPage).Scan(&results)
	} else {
		db.DB.Raw(`
			SELECT
				date(created_at, '+8 hours') AS date,
				COALESCE(AVG(fcp), 0) AS avg_fcp,
				COALESCE(AVG(lcp), 0) AS avg_lcp,
				COALESCE(AVG(ttfb), 0) AS avg_ttfb,
				COUNT(*) AS count
			FROM perf_metrics
			WHERE date(created_at, '+8 hours') >= ?
			GROUP BY date(created_at, '+8 hours')
			ORDER BY date ASC
		`, start).Scan(&results)
	}

	// 补全没有数据的日期 + 保留两位小数
	loc := time.FixedZone("CST", 8*3600)
	dateMap := make(map[string]models.PerfTrendItem, len(results))
	for _, r := range results {
		r.AvgFcp = math.Round(r.AvgFcp*100) / 100
		r.AvgLcp = math.Round(r.AvgLcp*100) / 100
		r.AvgTtfb = math.Round(r.AvgTtfb*100) / 100
		dateMap[r.Date] = r
	}

	fullResults := make([]models.PerfTrendItem, 0, days)
	for i := 0; i < days; i++ {
		d := time.Now().In(loc).AddDate(0, 0, -days+1+i).Format("2006-01-02")
		if item, ok := dateMap[d]; ok {
			fullResults = append(fullResults, item)
		} else {
			fullResults = append(fullResults, models.PerfTrendItem{Date: d})
		}
	}

	c.JSON(http.StatusOK, gin.H{"list": fullResults})
}

// GetLogs GET /api/perf/logs?page=1&pageSize=20&path=&days=7 — 原始日志查询（需鉴权）
func GetLogs(c *gin.Context) {
	pageNum := 1
	pageSize := 20
	days := parseDays(c)

	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		pageNum = p
	}
	if ps, err := strconv.Atoi(c.Query("pageSize")); err == nil && ps > 0 && ps <= 100 {
		pageSize = ps
	}

	start := startDateCST(days)
	filterPath := strings.TrimSpace(c.Query("path"))

	query := db.DB.Model(&models.PerfMetric{}).Where("date(created_at, '+8 hours') >= ?", start)
	if filterPath != "" {
		query = query.Where("page = ?", filterPath)
	}

	// 总数
	var total int64
	query.Count(&total)

	// 分页查询
	var records []models.PerfMetric
	offset := (pageNum - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records)

	// 格式化输出
	loc := time.FixedZone("CST", 8*3600)
	list := make([]models.PerfLogItem, 0, len(records))
	for _, r := range records {
		list = append(list, models.PerfLogItem{
			ID:         r.ID,
			Page:       r.Page,
			FCP:        r.FCP,
			LCP:        r.LCP,
			TTFB:       r.TTFB,
			DOMReady:   r.DOMReady,
			LoadTime:   r.LoadTime,
			JSLoad:     r.JSLoad,
			DeviceType: r.DeviceType,
			CreatedAt:  r.CreatedAt.In(loc).Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"list":     list,
		"total":    total,
		"page":     pageNum,
		"pageSize": pageSize,
	})
}
