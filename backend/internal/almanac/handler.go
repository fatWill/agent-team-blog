package almanac

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fatWill/agent-team-blog/backend/models"
	algo "github.com/fatWill/agent-team-blog/backend/pkg/almanac"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
)

func repo() *Repo { return NewRepo(db.DB) }

func badRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": msg})
}

// GetDay GET /api/toys/tuyile/almanac?date=YYYY-MM-DD（缺省为北京时间今天）
func GetDay(c *gin.Context) {
	date := strings.TrimSpace(c.Query("date"))
	if date == "" {
		date = algo.TodayCST()
	}
	if _, err := algo.ParseDate(date); err != nil {
		badRequest(c, err.Error())
		return
	}
	d, err := repo().GetOrCompute(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "查询黄历失败"})
		return
	}
	// 黄历按天不变，允许 CDN/浏览器缓存 10 分钟；缺省日期（今天）跨零点会变，缓存更短
	if c.Query("date") == "" {
		c.Header("Cache-Control", "public, max-age=60")
	} else {
		c.Header("Cache-Control", "public, max-age=600")
	}
	c.JSON(http.StatusOK, d)
}

// GetMonth GET /api/toys/tuyile/almanac/month?month=YYYY-MM（缺省为北京时间当月）
func GetMonth(c *gin.Context) {
	ym := strings.TrimSpace(c.Query("month"))
	if ym == "" {
		ym = time.Now().In(algo.CST).Format("2006-01")
	}
	first, err := algo.ParseDate(ym + "-01")
	if err != nil || len(ym) != 7 {
		badRequest(c, "month 格式应为 YYYY-MM")
		return
	}
	last := first.AddDate(0, 1, -1)
	rows, err := repo().ListRange(first.Format(algo.DateLayout), last.Format(algo.DateLayout))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "statusCode": 500, "statusMessage": "查询黄历失败"})
		return
	}
	byDate := make(map[string]*models.AlmanacDay, len(rows))
	for i := range rows {
		byDate[rows[i].Date] = &rows[i]
	}

	list := make([]models.AlmanacMonthItem, 0, last.Day())
	for t := first; !t.After(last); t = t.AddDate(0, 0, 1) {
		d, ok := byDate[t.Format(algo.DateLayout)]
		if !ok {
			d = algo.Compute(t)
		}
		list = append(list, toMonthItem(d))
	}
	c.Header("Cache-Control", "public, max-age=600")
	c.JSON(http.StatusOK, gin.H{
		"month": ym,
		"year":  first.Year(),
		"days":  len(list),
		"list":  list,
	})
}

func toMonthItem(d *models.AlmanacDay) models.AlmanacMonthItem {
	cell := d.LunarDayCN
	switch {
	case d.JieQi != "":
		cell = d.JieQi
	case d.LunarDay == 1:
		cell = d.LunarMonthCN
	}
	return models.AlmanacMonthItem{
		Date:        d.Date,
		Week:        d.Week,
		LunarDate:   d.LunarDate,
		CellText:    cell,
		GanZhiDay:   d.GanZhiDay,
		JieQi:       d.JieQi,
		Festivals:   nonNil(d.Festivals),
		ZhiShenType: d.ZhiShenType,
		YiTop:       top(d.Yi, 3),
		JiTop:       top(d.Ji, 3),
		ChongZodiac: d.ChongZodiac,
		Source:      d.Source,
	}
}

func top(arr []string, n int) []string {
	if len(arr) > n {
		return arr[:n]
	}
	return nonNil(arr)
}

func nonNil(arr []string) []string {
	if arr == nil {
		return []string{}
	}
	return arr
}

// FormatForPrompt 把单日黄历压缩成适合注入 LLM 上下文的一段文本
func FormatForPrompt(label string, d *models.AlmanacDay) string {
	var b strings.Builder
	fmt.Fprintf(&b, "【%s %s %s】农历%s年%s（生肖%s年）\n", label, d.Date, d.Week, d.LunarYearCN, d.LunarDate, d.ZodiacYear)
	fmt.Fprintf(&b, "干支：%s年 %s月 %s日；纳音：年%s 月%s 日%s\n", d.GanZhiYear, d.GanZhiMonth, d.GanZhiDay, d.NaYinYear, d.NaYinMonth, d.NaYinDay)
	if d.JieQi != "" || len(d.Festivals) > 0 {
		fmt.Fprintf(&b, "节气：%s；节日：%s\n", d.JieQi, strings.Join(d.Festivals, "、"))
	}
	fmt.Fprintf(&b, "宜：%s\n忌：%s\n", strings.Join(d.Yi, " "), strings.Join(d.Ji, " "))
	fmt.Fprintf(&b, "冲煞：%s；值神：%s(%s)；建除：%s；星宿：%s(%s)\n", d.ChongDesc, d.ZhiShen, d.ZhiShenType, d.JianChu, d.Xiu, d.XiuLuck)
	fmt.Fprintf(&b, "彭祖百忌：%s；胎神：%s；喜神%s 财神%s 福神%s\n", d.PengZu, d.TaiShen, d.XiShen, d.CaiShen, d.FuShen)
	fmt.Fprintf(&b, "吉神宜趋：%s；凶神宜忌：%s\n", strings.Join(d.JiShen, " "), strings.Join(d.XiongSha, " "))
	ji := make([]string, 0, 6)
	for _, h := range d.Hours {
		if h.Luck == "吉" {
			ji = append(ji, h.Zhi+"时")
		}
	}
	fmt.Fprintf(&b, "吉时：%s；数据来源：%s\n", strings.Join(ji, " "), d.Source)
	return b.String()
}
