package pv

import (
	"context"
	"fmt"
	"log"
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

// ==================== User-Agent 解析 ====================

// parseDeviceType 从 UA 解析设备类型
func parseDeviceType(ua string) string {
	if strings.Contains(ua, "Mobile") || strings.Contains(ua, "iPhone") || strings.Contains(ua, "Android") {
		return "mobile"
	}
	if strings.Contains(ua, "iPad") || strings.Contains(ua, "Tablet") {
		return "tablet"
	}
	return "desktop"
}

// parseBrowser 从 UA 解析浏览器
func parseBrowser(ua string) string {
	switch {
	case strings.Contains(ua, "Edg/"):
		return "Edge"
	case strings.Contains(ua, "OPR/") || strings.Contains(ua, "Opera"):
		return "Opera"
	case strings.Contains(ua, "Chrome"):
		return "Chrome"
	case strings.Contains(ua, "Firefox"):
		return "Firefox"
	case strings.Contains(ua, "Safari"):
		return "Safari"
	default:
		return "Other"
	}
}

// parseOS 从 UA 解析操作系统
func parseOS(ua string) string {
	switch {
	case strings.Contains(ua, "iPhone"):
		return "iOS"
	case strings.Contains(ua, "iPad"):
		return "iPadOS"
	case strings.Contains(ua, "Android"):
		return "Android"
	case strings.Contains(ua, "Windows"):
		return "Windows"
	case strings.Contains(ua, "Mac OS X"):
		return "macOS"
	case strings.Contains(ua, "Linux"):
		return "Linux"
	default:
		return "Other"
	}
}

// getTodayCST 获取 UTC+8 的今天日期字符串
func getTodayCST() string {
	loc := time.FixedZone("CST", 8*3600)
	return time.Now().In(loc).Format("2006-01-02")
}

// ==================== Handlers ====================

// RecordPV POST /api/pv/record — 上报访问记录（无需鉴权）
func RecordPV(c *gin.Context) {
	var body struct {
		Path     string `json:"path"`
		DeviceID string `json:"device_id"`
		Referer  string `json:"referer"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "参数解析失败"})
		return
	}

	path := strings.TrimSpace(body.Path)
	deviceID := strings.TrimSpace(body.DeviceID)

	if path == "" || deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "path 和 device_id 不能为空"})
		return
	}

	// 防刷：同一 device_id + path 60秒内只记录一次
	redisKey := fmt.Sprintf("pv:%s:%s", deviceID, path)
	ctx := context.Background()
	exists, err := rds.RDB.Exists(ctx, redisKey).Result()
	if err != nil {
		// Redis 不可用时降级，直接写入
		log.Printf("⚠️ PV 防刷 Redis 查询失败（降级写入）: %v", err)
	} else if exists > 0 {
		// 60s 内重复访问，静默返回成功
		c.JSON(http.StatusOK, gin.H{"ok": true})
		return
	}

	// 设置防刷 key，60s 过期
	if err := rds.RDB.Set(ctx, redisKey, "1", 60*time.Second).Err(); err != nil {
		log.Printf("⚠️ PV 防刷 Redis 写入失败（降级继续）: %v", err)
	}

	// 解析 User-Agent
	ua := c.GetHeader("User-Agent")
	ip := middleware.GetClientIP(c)

	record := models.PageView{
		Path:       path,
		DeviceID:   deviceID,
		IP:         ip,
		UserAgent:  ua,
		DeviceType: parseDeviceType(ua),
		Browser:    parseBrowser(ua),
		OS:         parseOS(ua),
		Referer:    strings.TrimSpace(body.Referer),
	}

	if err := db.DB.Create(&record).Error; err != nil {
		log.Printf("❌ PV 写入失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": "写入失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// GetTrend GET /api/pv/trend?days=7 — PV/UV 趋势（需鉴权）
func GetTrend(c *gin.Context) {
	days := 7
	if d, err := strconv.Atoi(c.Query("days")); err == nil && (d == 7 || d == 30) {
		days = d
	}

	// 计算起始日期（UTC+8）
	loc := time.FixedZone("CST", 8*3600)
	startDate := time.Now().In(loc).AddDate(0, 0, -days+1).Format("2006-01-02")

	var results []models.PVTrendItem

	// SQLite 日期函数：date(created_at, '+8 hours') 将 UTC 转为 CST
	db.DB.Raw(`
		SELECT date(created_at, '+8 hours') AS date,
		       COUNT(*) AS pv,
		       COUNT(DISTINCT device_id) AS uv
		FROM page_views
		WHERE date(created_at, '+8 hours') >= ?
		GROUP BY date(created_at, '+8 hours')
		ORDER BY date ASC
	`, startDate).Scan(&results)

	// 补全没有数据的日期
	dateMap := make(map[string]models.PVTrendItem, len(results))
	for _, r := range results {
		dateMap[r.Date] = r
	}

	fullResults := make([]models.PVTrendItem, 0, days)
	for i := 0; i < days; i++ {
		d := time.Now().In(loc).AddDate(0, 0, -days+1+i).Format("2006-01-02")
		if item, ok := dateMap[d]; ok {
			fullResults = append(fullResults, item)
		} else {
			fullResults = append(fullResults, models.PVTrendItem{Date: d, PV: 0, UV: 0})
		}
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": fullResults})
}

// GetTopPages GET /api/pv/top-pages?days=7 — Top5 页面（需鉴权）
func GetTopPages(c *gin.Context) {
	days := 7
	if d, err := strconv.Atoi(c.Query("days")); err == nil && (d == 7 || d == 30) {
		days = d
	}

	loc := time.FixedZone("CST", 8*3600)
	startDate := time.Now().In(loc).AddDate(0, 0, -days+1).Format("2006-01-02")

	var results []models.TopPageItem
	db.DB.Raw(`
		SELECT path,
		       COUNT(*) AS pv,
		       COUNT(DISTINCT device_id) AS uv
		FROM page_views
		WHERE date(created_at, '+8 hours') >= ?
		GROUP BY path
		ORDER BY pv DESC
		LIMIT 5
	`, startDate).Scan(&results)

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": results})
}

// GetLogs GET /api/pv/logs?page=1&page_size=20&path=&date= — 访问日志列表（需鉴权）
func GetLogs(c *gin.Context) {
	page := 1
	pageSize := 20

	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(c.Query("page_size")); err == nil && ps > 0 && ps <= 100 {
		pageSize = ps
	}

	filterPath := strings.TrimSpace(c.Query("path"))
	filterDate := strings.TrimSpace(c.Query("date"))

	query := db.DB.Model(&models.PageView{})

	if filterPath != "" {
		query = query.Where("path = ?", filterPath)
	}
	if filterDate != "" {
		// 按 CST 日期筛选
		query = query.Where("date(created_at, '+8 hours') = ?", filterDate)
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 分页查询
	var records []models.PageView
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records)

	// 格式化输出
	loc := time.FixedZone("CST", 8*3600)
	list := make([]gin.H, 0, len(records))
	for _, r := range records {
		list = append(list, gin.H{
			"id":          r.ID,
			"path":        r.Path,
			"device_id":   r.DeviceID,
			"ip":          r.IP,
			"device_type": r.DeviceType,
			"browser":     r.Browser,
			"os":          r.OS,
			"referer":     r.Referer,
			"created_at":  r.CreatedAt.In(loc).Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetOverview GET /api/pv/overview — 统计概览（需鉴权）
func GetOverview(c *gin.Context) {
	today := getTodayCST()

	var overview models.PVOverview

	// 今日 PV
	db.DB.Raw(`SELECT COUNT(*) FROM page_views WHERE date(created_at, '+8 hours') = ?`, today).Scan(&overview.TodayPV)

	// 今日 UV
	db.DB.Raw(`SELECT COUNT(DISTINCT device_id) FROM page_views WHERE date(created_at, '+8 hours') = ?`, today).Scan(&overview.TodayUV)

	// 总 PV
	db.DB.Raw(`SELECT COUNT(*) FROM page_views`).Scan(&overview.TotalPV)

	// 总 UV
	db.DB.Raw(`SELECT COUNT(DISTINCT device_id) FROM page_views`).Scan(&overview.TotalUV)

	c.JSON(http.StatusOK, gin.H{"ok": true, "data": overview})
}
