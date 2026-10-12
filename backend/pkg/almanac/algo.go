// Package almanac 老黄历基础能力：lunar-go 算法计算、huangli.com 页面解析、交叉校验合并。
//
// 该包只做纯计算与解析，不访问数据库，便于爬虫与 HTTP 服务复用及单测。
package almanac

import (
	"container/list"
	"fmt"
	"strings"
	"time"

	"github.com/6tail/lunar-go/LunarUtil"
	"github.com/6tail/lunar-go/calendar"

	"github.com/fatWill/agent-team-blog/backend/models"
)

// CST 北京时间（中国无夏令时，固定 +8，避免依赖服务器 tzdata）
var CST = time.FixedZone("CST", 8*3600)

// DateLayout 公历日期格式
const DateLayout = "2006-01-02"

// hourZhi 12 时辰地支顺序
var hourZhi = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

// hourRange 12 时辰时间范围（子时按传统口径跨日 23:00-00:59）
var hourRange = []string{
	"23:00-00:59", "01:00-02:59", "03:00-04:59", "05:00-06:59", "07:00-08:59", "09:00-10:59",
	"11:00-12:59", "13:00-14:59", "15:00-16:59", "17:00-18:59", "19:00-20:59", "21:00-22:59",
}

// algoOnlyFields 站点页面不提供、始终由算法补齐的字段
var algoOnlyFields = []string{
	"lunarYear", "lunarYearCn", "zodiacYear", "jieQi", "currentJieQi", "festivals",
	"chongGanZhi", "chongAge", "naYinYear", "naYinMonth", "zhiShenType", "xiuGong",
	"xiShen", "caiShen", "fuShen", "yangGui", "yinGui", "liuYao", "wuHou", "yueXiang",
}

// TodayCST 返回北京时间今天的日期字符串
func TodayCST() string {
	return time.Now().In(CST).Format(DateLayout)
}

// ParseDate 解析 YYYY-MM-DD，并限制在算法库可靠范围内（1900-2100）
func ParseDate(s string) (time.Time, error) {
	t, err := time.ParseInLocation(DateLayout, strings.TrimSpace(s), CST)
	if err != nil {
		return time.Time{}, fmt.Errorf("日期格式应为 YYYY-MM-DD")
	}
	if t.Year() < 1900 || t.Year() > 2100 {
		return time.Time{}, fmt.Errorf("日期超出支持范围（1900-2100）")
	}
	return t, nil
}

// Compute 用 lunar-go 计算某天的完整黄历（source=lunar-go）
func Compute(t time.Time) *models.AlmanacDay {
	solar := calendar.NewSolarFromYmd(t.Year(), int(t.Month()), t.Day())
	l := solar.GetLunar()

	month := l.GetMonth()
	monthCN := l.GetMonthInChinese() + "月"
	dayCN := l.GetDayInChinese()

	chongGZ := l.GetDayChongGan() + l.GetDayChong()
	chongAge := chongAgeOf(l.GetYear(), chongGZ)

	festivals := make([]string, 0, 4)
	festivals = appendList(festivals, solar.GetFestivals())
	festivals = appendList(festivals, l.GetFestivals())
	festivals = appendList(festivals, l.GetOtherFestivals())

	currentJieQi := l.GetJieQi()
	if currentJieQi == "" {
		if jq := l.GetPrevJieQiByWholeDay(true); jq != nil {
			currentJieQi = jq.GetName()
		}
	}

	day := &models.AlmanacDay{
		Date:         t.Format(DateLayout),
		Week:         "星期" + solar.GetWeekInChinese(),
		LunarYear:    l.GetYear(),
		LunarMonth:   abs(month),
		LunarDay:     l.GetDay(),
		IsLeapMonth:  month < 0,
		LunarMonthCN: monthCN,
		LunarDayCN:   dayCN,
		LunarDate:    monthCN + dayCN,
		LunarYearCN:  l.GetYearInChinese(),
		ZodiacYear:   l.GetYearShengXiao(),
		GanZhiYear:   l.GetYearInGanZhi(),
		GanZhiMonth:  l.GetMonthInGanZhi(),
		GanZhiDay:    l.GetDayInGanZhi(),
		JieQi:        l.GetJieQi(),
		CurrentJieQi: currentJieQi,
		Festivals:    festivals,
		Yi:           listToStrings(l.GetDayYi()),
		Ji:           listToStrings(l.GetDayJi()),
		ChongZodiac:  l.GetDayChongShengXiao(),
		ChongGanZhi:  chongGZ,
		ChongAge:     chongAge,
		Sha:          l.GetDaySha(),
		NaYinYear:    l.GetYearNaYin(),
		NaYinMonth:   l.GetMonthNaYin(),
		NaYinDay:     l.GetDayNaYin(),
		ZhiShen:      l.GetDayTianShen(),
		ZhiShenType:  l.GetDayTianShenType(),
		JianChu:      l.GetZhiXing(),
		Xiu:          l.GetXiu() + l.GetZheng() + l.GetAnimal(),
		XiuLuck:      l.GetXiuLuck(),
		XiuGong:      l.GetGong() + "方" + l.GetShou(),
		PengZu:       l.GetPengZuGan() + " " + l.GetPengZuZhi(),
		JiShen:       listToStrings(l.GetDayJiShen()),
		XiongSha:     listToStrings(l.GetDayXiongSha()),
		TaiShen:      l.GetDayPositionTai(),
		XiShen:       l.GetDayPositionXiDesc(),
		CaiShen:      l.GetDayPositionCaiDesc(),
		FuShen:       l.GetDayPositionFuDesc(),
		YangGui:      l.GetDayPositionYangGuiDesc(),
		YinGui:       l.GetDayPositionYinGuiDesc(),
		LiuYao:       l.GetLiuYao(),
		WuHou:        l.GetWuHou(),
		YueXiang:     l.GetYueXiang(),
		Hours:        computeHours(l),
		Source:       models.AlmanacSourceAlgo,
		AlgoFields:   []string{"*"},
		Mismatches:   []string{},
	}
	day.ChongDesc = formatChongDesc(day)
	return day
}

// computeHours 计算 12 时辰吉凶（取早子时，与传统黄历子时干支一致）
func computeHours(l *calendar.Lunar) []models.AlmanacHour {
	times := l.GetTimes()
	hours := make([]models.AlmanacHour, 0, 12)
	for i := 0; i < 12 && i < len(times); i++ {
		t := times[i]
		hours = append(hours, models.AlmanacHour{
			Zhi:      hourZhi[i],
			GanZhi:   t.GetGanZhi(),
			Range:    hourRange[i],
			TianShen: t.GetTianShen(),
			Luck:     t.GetTianShenLuck(),
			Chong:    t.GetChongShengXiao(),
			Sha:      t.GetSha(),
			Yi:       listToStrings(t.GetYi()),
			Ji:       listToStrings(t.GetJi()),
			JiShen:   []string{},
			XiongSha: []string{},
		})
	}
	return hours
}

// chongAgeOf 计算被冲生肖者在当年（农历年）的虚岁
func chongAgeOf(lunarYear int, chongGZ string) int {
	idx := LunarUtil.GetJiaZiIndex(chongGZ)
	if idx < 0 {
		return 0
	}
	// 公元 4 年为甲子年
	yearIdx := ((lunarYear-4)%60 + 60) % 60
	diff := ((yearIdx-idx)%60 + 60) % 60
	return diff + 1
}

func formatChongDesc(d *models.AlmanacDay) string {
	if d.ChongAge > 0 {
		return fmt.Sprintf("冲%s(%s)%d岁 煞%s", d.ChongZodiac, d.ChongGanZhi, d.ChongAge, d.Sha)
	}
	return fmt.Sprintf("冲%s(%s) 煞%s", d.ChongZodiac, d.ChongGanZhi, d.Sha)
}

// BirthInfo 出生日期的命理基础信息
type BirthInfo struct {
	Date            string `json:"date"`
	Calendar        string `json:"calendar"` // solar / lunar（用户输入口径）
	LunarDate       string `json:"lunarDate"`
	ZodiacSpring    string `json:"zodiacSpring"` // 生肖（春节口径，民间常用）
	ZodiacLiChun    string `json:"zodiacLiChun"` // 生肖（立春口径，八字常用）
	YearPillar      string `json:"yearPillar"`   // 八字年柱（立春口径）
	MonthPillar     string `json:"monthPillar"`
	DayPillar       string `json:"dayPillar"`
	TimePillar      string `json:"timePillar,omitempty"`
	YearNaYin       string `json:"yearNaYin"`
	DayNaYin        string `json:"dayNaYin"`
	DayMasterWuXing string `json:"dayMasterWuXing"` // 日主五行
}

// ComputeBirth 计算出生信息；hour<0 表示时辰未知
func ComputeBirth(t time.Time, hour int, cal string) BirthInfo {
	h := 12
	if hour >= 0 {
		h = hour
	}
	solar := calendar.NewSolar(t.Year(), int(t.Month()), t.Day(), h, 0, 0)
	l := solar.GetLunar()
	ec := l.GetEightChar()
	info := BirthInfo{
		Date:            t.Format(DateLayout),
		Calendar:        cal,
		LunarDate:       fmt.Sprintf("%s年%s月%s", l.GetYearInChinese(), l.GetMonthInChinese(), l.GetDayInChinese()),
		ZodiacSpring:    l.GetYearShengXiao(),
		ZodiacLiChun:    l.GetYearShengXiaoByLiChun(),
		YearPillar:      ec.GetYear(),
		MonthPillar:     ec.GetMonth(),
		DayPillar:       ec.GetDay(),
		YearNaYin:       ec.GetYearNaYin(),
		DayNaYin:        ec.GetDayNaYin(),
		DayMasterWuXing: string([]rune(ec.GetDayWuXing())[0:1]),
	}
	if hour >= 0 {
		info.TimePillar = ec.GetTime()
	}
	return info
}

// YearInfo 流年基础信息
type YearInfo struct {
	Year   int    `json:"year"`
	GanZhi string `json:"ganZhi"`
	Zodiac string `json:"zodiac"`
	NaYin  string `json:"naYin"`
}

// ComputeYear 计算某公历年对应农历年（春节起算）的干支、生肖、纳音
func ComputeYear(year int) YearInfo {
	l := calendar.NewLunarFromYmd(year, 1, 1)
	return YearInfo{Year: year, GanZhi: l.GetYearInGanZhi(), Zodiac: l.GetYearShengXiao(), NaYin: l.GetYearNaYin()}
}

// LunarToSolar 农历转公历；非法日期（如不存在的闰月/三十）返回 error 而不是 panic
func LunarToSolar(year, month, day int, leap bool) (t time.Time, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("农历日期不存在: %v", r)
		}
	}()
	m := month
	if leap {
		m = -month
	}
	s := calendar.NewLunarFromYmd(year, m, day).GetSolar()
	return time.Date(s.GetYear(), time.Month(s.GetMonth()), s.GetDay(), 0, 0, 0, 0, CST), nil
}

// ZodiacRelation 生肖与流年生肖的关系（本命/冲/合/刑/害/破），用于运势提示
func ZodiacRelation(birth, year string) []string {
	if birth == "" || year == "" {
		return nil
	}
	rel := make([]string, 0, 3)
	if birth == year {
		rel = append(rel, "本命年（值太岁）")
	}
	pairs := map[string][][2]string{
		"相冲（冲太岁）": {{"鼠", "马"}, {"牛", "羊"}, {"虎", "猴"}, {"兔", "鸡"}, {"龙", "狗"}, {"蛇", "猪"}},
		"六合":      {{"鼠", "牛"}, {"虎", "猪"}, {"兔", "狗"}, {"龙", "鸡"}, {"蛇", "猴"}, {"马", "羊"}},
		"相害（害太岁）": {{"鼠", "羊"}, {"牛", "马"}, {"虎", "蛇"}, {"兔", "龙"}, {"猴", "猪"}, {"鸡", "狗"}},
		"相破（破太岁）": {{"鼠", "鸡"}, {"牛", "龙"}, {"虎", "猪"}, {"兔", "马"}, {"蛇", "猴"}, {"羊", "狗"}},
		"相刑（刑太岁）": {{"鼠", "兔"}, {"虎", "蛇"}, {"蛇", "猴"}, {"虎", "猴"}, {"牛", "羊"}, {"牛", "狗"}, {"羊", "狗"}},
	}
	order := []string{"相冲（冲太岁）", "六合", "相害（害太岁）", "相破（破太岁）", "相刑（刑太岁）"}
	for _, name := range order {
		for _, p := range pairs[name] {
			if (p[0] == birth && p[1] == year) || (p[1] == birth && p[0] == year) {
				rel = append(rel, name)
				break
			}
		}
	}
	sanHe := [][]string{{"猴", "鼠", "龙"}, {"虎", "马", "狗"}, {"蛇", "鸡", "牛"}, {"猪", "兔", "羊"}}
	for _, g := range sanHe {
		if birth != year && contains(g, birth) && contains(g, year) {
			rel = append(rel, "三合")
		}
	}
	if birth == year {
		// 辰午酉亥自刑
		if contains([]string{"龙", "马", "鸡", "猪"}, birth) {
			rel = append(rel, "自刑（刑太岁）")
		}
	}
	if len(rel) == 0 {
		rel = append(rel, "与太岁无明显冲合")
	}
	return rel
}

func listToStrings(l *list.List) []string {
	out := make([]string, 0, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		if s, ok := e.Value.(string); ok && s != "" {
			out = append(out, s)
		}
	}
	return out
}

func appendList(dst []string, l *list.List) []string {
	for _, s := range listToStrings(l) {
		if !contains(dst, s) {
			dst = append(dst, s)
		}
	}
	return dst
}

func contains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
