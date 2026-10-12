package tuyile

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	algo "github.com/fatWill/agent-team-blog/backend/pkg/almanac"
)

// mentionedDate 对话中识别出的日期
type mentionedDate struct {
	Date     time.Time
	Raw      string
	Calendar string // solar / lunar
	IsBirth  bool
	Hour     int // 出生时辰对应的小时（-1 未知）
}

var (
	// 公历：2024-05-01 / 2024/5/1 / 2024.5.1 / 2024年5月1日(号)
	reSolar = regexp.MustCompile(`(19\d{2}|20\d{2}|2100)\s*(?:[-/.]|年)\s*(\d{1,2})\s*(?:[-/.]|月)\s*(\d{1,2})\s*[日号]?`)
	// 农历数字：农历1990年闰5月初三 / 阴历1990年五月初三 / 农历 1990-5-3
	reLunarNum = regexp.MustCompile(`(?:农历|阴历|旧历)\s*(19\d{2}|20\d{2})\s*(?:[-/.]|年)\s*(闰)?\s*([0-9]{1,2}|[正一二三四五六七八九十冬腊]{1,2})\s*(?:[-/.]|月)\s*(初[一二三四五六七八九十]|十[一二三四五六七八九]?|二十|廿[一二三四五六七八九]?|三十|[0-9]{1,2})\s*[日号]?`)
	// 预测年份：2027年运势 / 预测2026 / 看看2028年
	reYear = regexp.MustCompile(`(20\d{2})\s*年?`)
	// 时辰：子时 / 早上8点 / 晚上9点 / 21:30 / 21点
	reShiChen = regexp.MustCompile(`([子丑寅卯辰巳午未申酉戌亥])时`)
	reClock   = regexp.MustCompile(`(凌晨|早上|上午|中午|下午|傍晚|晚上)?\s*(\d{1,2})\s*(?:[:：]\d{2}|点|时)`)

	birthHints = []string{"生日", "出生", "生于", "我是", "阳历生", "农历生", "生辰"}
)

var cnNum = map[string]int{"正": 1, "一": 1, "二": 2, "三": 3, "四": 4, "五": 5, "六": 6, "七": 7, "八": 8, "九": 9, "十": 10, "冬": 11, "腊": 12}

func parseCNMonth(s string) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	switch s {
	case "十一", "冬":
		return 11
	case "十二", "腊":
		return 12
	}
	return cnNum[s]
}

func parseCNDay(s string) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	r := []rune(s)
	switch {
	case strings.HasPrefix(s, "初"):
		return cnNum[string(r[1:])]
	case s == "十":
		return 10
	case s == "二十":
		return 20
	case s == "三十":
		return 30
	case strings.HasPrefix(s, "十"):
		return 10 + cnNum[string(r[1:])]
	case strings.HasPrefix(s, "廿"):
		if len(r) == 1 {
			return 20
		}
		return 20 + cnNum[string(r[1:])]
	}
	return 0
}

// extractDates 从单条用户消息中抽取日期（公历与农历），并判断是否为生日
func extractDates(text string) []mentionedDate {
	out := make([]mentionedDate, 0, 2)
	birthCtx := hasAny(text, birthHints)
	hour := extractHour(text)

	lunarSpans := reLunarNum.FindAllStringSubmatchIndex(text, -1)
	for _, idx := range lunarSpans {
		m := reLunarNum.FindStringSubmatch(text[idx[0]:idx[1]])
		y, _ := strconv.Atoi(m[1])
		mon, day := parseCNMonth(m[3]), parseCNDay(m[4])
		if mon < 1 || mon > 12 || day < 1 || day > 30 {
			continue
		}
		t, err := algo.LunarToSolar(y, mon, day, m[2] != "")
		if err != nil {
			continue
		}
		out = append(out, mentionedDate{Date: t, Raw: m[0], Calendar: "lunar", IsBirth: birthCtx, Hour: hour})
	}

	for _, idx := range reSolar.FindAllStringSubmatchIndex(text, -1) {
		if insideAny(idx[0], lunarSpans) {
			continue
		}
		m := reSolar.FindStringSubmatch(text[idx[0]:idx[1]])
		y, _ := strconv.Atoi(m[1])
		mon, _ := strconv.Atoi(m[2])
		day, _ := strconv.Atoi(m[3])
		t := time.Date(y, time.Month(mon), day, 0, 0, 0, 0, algo.CST)
		if t.Month() != time.Month(mon) || t.Day() != day {
			continue // 2月30日之类的非法日期
		}
		cal := "solar"
		// 「农历」在日期前不远处出现但没被 reLunarNum 匹配（如农历 1990-05-03）时按农历解释
		prefix := text[max(0, idx[0]-12):idx[0]]
		if strings.Contains(prefix, "农历") || strings.Contains(prefix, "阴历") {
			if lt, err := algo.LunarToSolar(y, mon, day, false); err == nil {
				t, cal = lt, "lunar"
			}
		}
		out = append(out, mentionedDate{Date: t, Raw: m[0], Calendar: cal, IsBirth: birthCtx, Hour: hour})
	}
	return out
}

func extractHour(text string) int {
	if m := reShiChen.FindStringSubmatch(text); m != nil {
		zhi := "子丑寅卯辰巳午未申酉戌亥"
		i := strings.Index(zhi, m[1]) / len("子")
		return (i*2 + 24) % 24 // 子=0 丑=2 ...
	}
	if m := reClock.FindStringSubmatch(text); m != nil {
		h, _ := strconv.Atoi(m[2])
		if (m[1] == "下午" || m[1] == "傍晚" || m[1] == "晚上") && h < 12 {
			h += 12
		}
		if h >= 0 && h < 24 {
			return h
		}
	}
	return -1
}

// extractTargetYear 抽取预测年份：排除日期中的年份，取最后一次提及
func extractTargetYear(text string) int {
	clean := reSolar.ReplaceAllString(text, " ")
	clean = reLunarNum.ReplaceAllString(clean, " ")
	ms := reYear.FindAllStringSubmatch(clean, -1)
	for i := len(ms) - 1; i >= 0; i-- {
		y, _ := strconv.Atoi(ms[i][1])
		if y >= 2000 && y <= 2100 {
			return y
		}
	}
	return 0
}

// convoFacts 从整段对话（仅 user 消息）中提取的事实
type convoFacts struct {
	Birth      *mentionedDate
	TargetYear int
	Mentioned  []mentionedDate // 非生日的日期（按出现顺序去重，最多 3 个）
}

func collectFacts(userTexts []string) convoFacts {
	var f convoFacts
	seen := map[string]bool{}
	for i, text := range userTexts {
		dates := extractDates(text)
		for j := range dates {
			d := dates[j]
			// 生日：带生日语义，或是第一个早于 2010 年的日期（大概率在回答「出生年月日」的追问）
			if d.IsBirth || (f.Birth == nil && d.Date.Year() < 2010) {
				dd := d
				f.Birth = &dd
				continue
			}
			key := d.Date.Format(algo.DateLayout)
			if !seen[key] {
				seen[key] = true
				f.Mentioned = append(f.Mentioned, d)
			}
		}
		if y := extractTargetYear(text); y > 0 {
			// 只有在生日已给出之后（或本句明显在说年运）才认作预测年份
			if f.Birth != nil || strings.Contains(text, "运") || strings.Contains(text, "预测") || i == len(userTexts)-1 {
				f.TargetYear = y
			}
		}
	}
	if len(f.Mentioned) > 3 {
		f.Mentioned = f.Mentioned[len(f.Mentioned)-3:]
	}
	return f
}

func hasAny(s string, subs []string) bool {
	for _, x := range subs {
		if strings.Contains(s, x) {
			return true
		}
	}
	return false
}

func insideAny(pos int, spans [][]int) bool {
	for _, sp := range spans {
		if pos >= sp[0] && pos < sp[1] {
			return true
		}
	}
	return false
}

func describeHour(h int) string {
	if h < 0 {
		return "未提供"
	}
	zhi := []string{"子", "丑", "丑", "寅", "寅", "卯", "卯", "辰", "辰", "巳", "巳", "午", "午", "未", "未", "申", "申", "酉", "酉", "戌", "戌", "亥", "亥", "子"}
	return fmt.Sprintf("约 %d 点（%s时）", h, zhi[h])
}
