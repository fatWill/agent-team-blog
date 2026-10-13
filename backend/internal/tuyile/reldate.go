package tuyile

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	algo "github.com/fatWill/agent-team-blog/backend/pkg/almanac"
)

// dateRange 相对日期解析结果（闭区间，北京时间零点）
type dateRange struct {
	Label string
	Start time.Time
	End   time.Time
}

// MaxRangeDays 单次黄历范围查询上限（工具与兜底注入共用）
const MaxRangeDays = 62

var (
	reWeekDay   = regexp.MustCompile(`(这|本|下下|下|上)?(?:周|星期|礼拜)([一二三四五六日天1-7])`)
	reYearMonth = regexp.MustCompile(`(?:(20\d{2})\s*年|(今年|明年|后年))?\s*(\d{1,2}|[一二三四五六七八九十]|十一|十二)\s*月份?(?:([上中下])旬)?`)
	cnDigit     = map[string]int{"一": 1, "二": 2, "三": 3, "四": 4, "五": 5, "六": 6, "七": 7, "八": 8, "九": 9, "十": 10, "十一": 11, "十二": 12, "日": 7, "天": 7}
)

func dayStart(t time.Time) time.Time {
	t = t.In(algo.CST)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, algo.CST)
}

// mondayOf 返回 t 所在周的周一（中国习惯：周一为一周第一天）
func mondayOf(t time.Time) time.Time {
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7
	}
	return t.AddDate(0, 0, 1-wd)
}

func monthRange(y, m int) (time.Time, time.Time) {
	first := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, algo.CST)
	return first, first.AddDate(0, 1, -1)
}

func parseMonthNum(s string) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return cnDigit[s]
}

// parseRelativeDates 从用户消息中解析相对日期 / 日期范围（今天、明天、后天、本周、下周X、下个月、明年三月、X月上旬…）。
// 规则偏保守：解析不出就不注入，交给模型通过工具或追问处理。结果去重，单个范围截断到 MaxRangeDays。
func parseRelativeDates(text string, now time.Time) []dateRange {
	today := dayStart(now)
	out := make([]dateRange, 0, 3)
	seen := map[string]bool{}
	add := func(label string, s, e time.Time) {
		if e.Before(s) {
			return
		}
		if e.Sub(s) > time.Duration(MaxRangeDays-1)*24*time.Hour {
			e = s.AddDate(0, 0, MaxRangeDays-1)
		}
		k := s.Format(algo.DateLayout) + e.Format(algo.DateLayout)
		if seen[k] {
			return
		}
		seen[k] = true
		out = append(out, dateRange{Label: label, Start: s, End: e})
	}

	// 单日：大后天要先于后天匹配
	single := []struct {
		kw     string
		offset int
	}{{"大后天", 3}, {"后天", 2}, {"明天", 1}, {"明儿", 1}, {"今天", 0}, {"今日", 0}, {"今儿", 0}}
	rest := text
	for _, s := range single {
		if strings.Contains(rest, s.kw) {
			d := today.AddDate(0, 0, s.offset)
			add(s.kw, d, d)
			rest = strings.ReplaceAll(rest, s.kw, " ")
		}
	}

	// 周X / 下周X
	thisMon := mondayOf(today)
	for _, m := range reWeekDay.FindAllStringSubmatch(text, -1) {
		n := cnDigit[m[2]]
		if n == 0 {
			n, _ = strconv.Atoi(m[2])
		}
		base := thisMon
		switch m[1] {
		case "下":
			base = thisMon.AddDate(0, 0, 7)
		case "下下":
			base = thisMon.AddDate(0, 0, 14)
		case "上":
			base = thisMon.AddDate(0, 0, -7)
		}
		d := base.AddDate(0, 0, n-1)
		if m[1] == "" && d.Before(today) {
			d = d.AddDate(0, 0, 7) // 只说「周五」且本周已过 → 指下一个周五
		}
		add(m[0], d, d)
	}
	textNoWD := reWeekDay.ReplaceAllString(text, " ")
	textNoWD = strings.NewReplacer("下周末", " ", "这周末", " ", "本周末", " ").Replace(textNoWD)

	// 周范围
	switch {
	case hasAny(textNoWD, []string{"下下周", "下下星期", "下下礼拜"}):
		add("下下周", thisMon.AddDate(0, 0, 14), thisMon.AddDate(0, 0, 20))
	case hasAny(textNoWD, []string{"下周", "下星期", "下礼拜", "下个星期", "下个礼拜"}):
		add("下周", thisMon.AddDate(0, 0, 7), thisMon.AddDate(0, 0, 13))
	}
	if hasAny(textNoWD, []string{"本周", "这周", "这星期", "这礼拜", "这个星期", "这个礼拜"}) {
		add("本周", thisMon, thisMon.AddDate(0, 0, 6))
	}
	if hasAny(text, []string{"周末", "这周末", "本周末"}) && !strings.Contains(text, "下周末") {
		add("本周末", thisMon.AddDate(0, 0, 5), thisMon.AddDate(0, 0, 6))
	}
	if strings.Contains(text, "下周末") {
		add("下周末", thisMon.AddDate(0, 0, 12), thisMon.AddDate(0, 0, 13))
	}
	if hasAny(text, []string{"最近几天", "这几天", "近几天", "未来几天", "最近一周", "未来一周", "七天内", "7天内"}) {
		add("未来 7 天", today, today.AddDate(0, 0, 6))
	}
	if hasAny(text, []string{"半个月", "两周内", "未来两周"}) {
		add("未来 14 天", today, today.AddDate(0, 0, 13))
	}

	// 月范围：本月 / 下个月 / 下下个月
	y, mo := today.Year(), int(today.Month())
	if hasAny(text, []string{"本月", "这个月", "这月"}) {
		s, e := monthRange(y, mo)
		if s.Before(today) {
			s = today
		}
		add("本月剩余", s, e)
	}
	if hasAny(text, []string{"月底", "月末"}) {
		_, e := monthRange(y, mo)
		add("本月底", e.AddDate(0, 0, -6), e)
	}
	if hasAny(text, []string{"下下个月", "下下月"}) {
		s, e := monthRange(y, mo+2)
		add("下下个月", s, e)
	} else if hasAny(text, []string{"下个月", "下月"}) {
		s, e := monthRange(y, mo+1)
		add("下个月", s, e)
	}

	// 显式月份：2027年3月 / 明年三月 / 3月上旬 / 十二月（未给年份时：本月及以后取今年，否则取明年）
	clean := reSolar.ReplaceAllString(text, " ")
	clean = reLunarNum.ReplaceAllString(clean, " ")
	for _, m := range reYearMonth.FindAllStringSubmatch(clean, -1) {
		mon := parseMonthNum(m[3])
		if mon < 1 || mon > 12 {
			continue
		}
		yy := y
		switch {
		case m[1] != "":
			yy, _ = strconv.Atoi(m[1])
		case m[2] == "明年":
			yy = y + 1
		case m[2] == "后年":
			yy = y + 2
		case m[2] == "今年":
		default:
			if mon < mo {
				yy = y + 1
			}
		}
		if yy < 1901 || yy > 2099 {
			continue
		}
		s, e := monthRange(yy, mon)
		label := strings.TrimSpace(m[0])
		switch m[4] {
		case "上":
			e = s.AddDate(0, 0, 9)
		case "中":
			s, e = s.AddDate(0, 0, 10), s.AddDate(0, 0, 19)
		case "下":
			s = s.AddDate(0, 0, 20)
		}
		add(label, s, e)
	}
	if len(out) > 3 {
		out = out[:3]
	}
	return out
}
