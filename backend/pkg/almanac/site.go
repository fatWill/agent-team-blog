package almanac

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/fatWill/agent-team-blog/backend/models"
)

// SiteDay huangli.com 日页解析结果（仅站点原始字段）
type SiteDay struct {
	Date        string
	Week        string
	LunarMonth  string // 冬月 / 闰六月
	LunarDay    string // 二十
	Yi          []string
	Ji          []string
	JiShen      []string
	XiongSha    []string
	ZhiShen     string
	NaYinDay    string
	JianChu     string
	PengZu      string
	ChongZodiac string
	Sha         string
	Xiu         string // 毕月乌
	XiuLuck     string // 吉 / 凶
	TaiShen     string
	BaZi        string // 年 月 日（取自子时四柱前三柱）
	Hours       []SiteHour
}

// SiteHour 时辰表一列
type SiteHour struct {
	GanZhi   string
	Range    string
	TianShen string
	Chong    string
	Luck     string
	JiShen   []string
	XiongSha []string
	Yi       []string
	Ji       []string
	Sha      string
}

var (
	reH3Date = regexp.MustCompile(`(\d{4})年(\d{2})月(\d{2})日\s*(星期.)`)
	reChong  = regexp.MustCompile(`.日冲(.)`)
	reXiu    = regexp.MustCompile(`^[东南西北]方(.{3})\((.)\)$`)
)

// SiteURL 返回 huangli.com 某天的日页 URL
func SiteURL(date string) string {
	// date: YYYY-MM-DD → /huangli/YYYY/MM_DD.html
	return fmt.Sprintf("https://www.huangli.com/huangli/%s/%s_%s.html", date[0:4], date[5:7], date[8:10])
}

// ParseSite 解析 huangli.com 日页。页面结构异常时返回 error，由调用方回退算法。
func ParseSite(r io.Reader, wantDate string) (*SiteDay, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 失败: %w", err)
	}
	info := doc.Find(".calendar-info .lunar-info").First()
	if info.Length() == 0 {
		return nil, fmt.Errorf("页面缺少 .lunar-info 主信息区")
	}

	d := &SiteDay{}
	h3 := strings.TrimSpace(info.Find(".lunar-top h3").First().Text())
	m := reH3Date.FindStringSubmatch(h3)
	if m == nil {
		return nil, fmt.Errorf("无法解析公历日期标题: %q", h3)
	}
	d.Date = m[1] + "-" + m[2] + "-" + m[3]
	d.Week = m[4]
	if d.Date != wantDate {
		return nil, fmt.Errorf("页面日期 %s 与请求日期 %s 不一致", d.Date, wantDate)
	}

	d.LunarMonth = clean(info.Find(".lunar-month").First().Text())
	d.LunarDay = clean(info.Find(".lunar-day").First().Text())

	// 宜/忌两列，每列依次为 [宜|忌] 与 [吉神宜趋|凶煞宜忌]
	info.Find(".yi-ji > .yi").Each(func(_ int, col *goquery.Selection) {
		title := clean(col.Find("h3").First().Text())
		lists := col.Find(".day-yi")
		first := spans(lists.Eq(0))
		second := spans(lists.Eq(1))
		switch title {
		case "宜":
			d.Yi, d.JiShen = first, second
		case "忌":
			d.Ji, d.XiongSha = first, second
		}
	})

	info.Find(".lunar-bottom .item").Each(func(_ int, it *goquery.Selection) {
		sp := it.ChildrenFiltered("span")
		if sp.Length() < 2 {
			return
		}
		k, v := clean(sp.Eq(0).Text()), clean(sp.Eq(1).Text())
		switch {
		case k == "值神":
			d.ZhiShen = v
		case k == "五行":
			d.NaYinDay = v
		case k == "十二建除":
			d.JianChu = v
		case k == "二十八宿":
			if mm := reXiu.FindStringSubmatch(v); mm != nil {
				d.Xiu, d.XiuLuck = mm[1], mm[2]
			} else {
				d.Xiu = v
			}
		case k == "今日胎神":
			d.TaiShen = v
		case strings.Contains(k, "日冲"):
			if mm := reChong.FindStringSubmatch(k); mm != nil {
				d.ChongZodiac = mm[1]
			}
			d.Sha = strings.TrimPrefix(v, "煞")
		}
	})
	info.Find(".lunar-bottom p").Each(func(_ int, p *goquery.Selection) {
		t := clean(p.Text())
		if strings.HasPrefix(t, "彭祖百忌") {
			d.PengZu = strings.TrimSpace(strings.TrimLeft(strings.TrimPrefix(t, "彭祖百忌"), "：:"))
		}
	})

	d.Hours, d.BaZi = parseHourTable(doc)
	return d, nil
}

// parseHourTable 解析 .w.times 下按行标签组织的 13 列时辰表
func parseHourTable(doc *goquery.Document) ([]SiteHour, string) {
	rows := map[string][]string{}
	doc.Find(".w.times table tr").Each(func(_ int, tr *goquery.Selection) {
		tds := tr.ChildrenFiltered("td")
		if tds.Length() != 13 {
			return
		}
		label := clean(tds.Eq(0).Text())
		if _, seen := rows[label]; seen {
			// 第二张表（时辰凶吉 6 列）结构不同，13 列只取首次出现
			return
		}
		vals := make([]string, 0, 12)
		tds.Slice(1, 13).Each(func(_ int, td *goquery.Selection) {
			vals = append(vals, clean(td.Text()))
		})
		rows[label] = vals
	})

	bazi := ""
	doc.Find(".w.times table tr").EachWithBreak(func(_ int, tr *goquery.Selection) bool {
		tds := tr.ChildrenFiltered("td")
		if tds.Length() > 1 && clean(tds.Eq(0).Text()) == "四柱八字" {
			parts := strings.Fields(clean(tds.Eq(1).Text()))
			if len(parts) >= 3 {
				bazi = strings.Join(parts[:3], " ")
			}
			return false
		}
		return true
	})

	gz := rows["时辰"]
	if len(gz) != 12 {
		return nil, bazi
	}
	hours := make([]SiteHour, 12)
	for i := 0; i < 12; i++ {
		hours[i] = SiteHour{
			GanZhi:   gz[i],
			Range:    at(rows["时刻"], i),
			TianShen: at(rows["星神"], i),
			Chong:    strings.TrimPrefix(at(rows["正冲"], i), "冲"),
			Luck:     at(rows["吉凶"], i),
			JiShen:   fields(at(rows["吉神"], i), " "),
			XiongSha: fields(at(rows["凶煞"], i), " "),
			Yi:       fields(at(rows["时宜"], i), ","),
			Ji:       fields(at(rows["时忌"], i), ","),
			Sha:      strings.TrimPrefix(at(rows["煞方"], i), "煞"),
		}
	}
	return hours, bazi
}

// Merge 以站点数据为主、算法为基准合并，并记录交叉校验不一致项。
//
// 规则：
//   - 干支、农历日期、冲煞属于确定性历法字段：以算法为准写入，不一致记 mismatch（站点值写入 mismatch 便于排查）
//   - 宜忌、吉神凶煞、时辰宜忌等流派差异字段：站点有值用站点，否则算法补齐并登记 algoFields
func Merge(algo *models.AlmanacDay, site *SiteDay) (*models.AlmanacDay, []string) {
	if site == nil {
		return algo, nil
	}
	out := *algo
	mismatches := make([]string, 0)
	filled := append([]string{}, algoOnlyFields...)

	check := func(field, siteVal, algoVal string) {
		if siteVal != "" && siteVal != algoVal {
			mismatches = append(mismatches, fmt.Sprintf("%s: site=%s algo=%s", field, siteVal, algoVal))
		}
	}

	// ---- 确定性字段：交叉校验 ----
	if site.BaZi != "" {
		check("ganZhi", site.BaZi, strings.Join([]string{algo.GanZhiYear, algo.GanZhiMonth, algo.GanZhiDay}, " "))
	} else {
		filled = append(filled, "ganZhi")
	}
	check("lunarDate", site.LunarMonth+site.LunarDay, algo.LunarDate)
	check("chongZodiac", site.ChongZodiac, algo.ChongZodiac)
	check("sha", site.Sha, algo.Sha)
	check("week", site.Week, algo.Week)
	// 二十八宿也是确定性字段；站点写法为「宿+动物+七曜」(毕乌月)，标准为「宿+七曜+动物」(毕月乌)，按字集比较
	if site.Xiu != "" && sortRunes(site.Xiu) != sortRunes(algo.Xiu) {
		mismatches = append(mismatches, fmt.Sprintf("xiu: site=%s algo=%s", site.Xiu, algo.Xiu))
	}
	if site.Xiu == "" {
		filled = append(filled, "xiu")
	}

	// ---- 站点优先字段 ----
	pickList := func(field string, siteVal []string, dst *[]string) {
		if len(siteVal) > 0 {
			*dst = siteVal
		} else {
			filled = append(filled, field)
		}
	}
	pickStr := func(field, siteVal string, dst *string) {
		if siteVal != "" {
			*dst = siteVal
		} else {
			filled = append(filled, field)
		}
	}
	pickList("yi", site.Yi, &out.Yi)
	pickList("ji", site.Ji, &out.Ji)
	pickList("jiShen", site.JiShen, &out.JiShen)
	pickList("xiongSha", site.XiongSha, &out.XiongSha)
	pickStr("zhiShen", site.ZhiShen, &out.ZhiShen)
	pickStr("naYinDay", site.NaYinDay, &out.NaYinDay)
	pickStr("jianChu", site.JianChu, &out.JianChu)
	pickStr("pengZu", site.PengZu, &out.PengZu)
	// 站点只给短句「甲不开仓 子不问卜」，算法为完整句「甲不开仓财物耗散 子不问卜自惹祸殃」；两半前缀一致时取完整句
	if sp, ap := strings.Fields(site.PengZu), strings.Fields(algo.PengZu); len(sp) == 2 && len(ap) == 2 &&
		strings.HasPrefix(ap[0], sp[0]) && strings.HasPrefix(ap[1], sp[1]) {
		out.PengZu = algo.PengZu
	}
	pickStr("taiShen", site.TaiShen, &out.TaiShen)
	pickStr("xiuLuck", site.XiuLuck, &out.XiuLuck)

	// 站点值神与算法值神不一致时，黄道/黑道类型以站点值神重新判定
	if site.ZhiShen != "" && site.ZhiShen != algo.ZhiShen {
		out.ZhiShenType = tianShenType(site.ZhiShen)
	}

	// ---- 时辰：逐列合并 ----
	if len(site.Hours) == 12 && len(algo.Hours) == 12 {
		hours := make([]models.AlmanacHour, 12)
		for i := 0; i < 12; i++ {
			a, s := algo.Hours[i], site.Hours[i]
			h := a
			if s.GanZhi != "" && s.GanZhi != a.GanZhi {
				mismatches = append(mismatches, fmt.Sprintf("hour[%s].ganZhi: site=%s algo=%s", a.Zhi, s.GanZhi, a.GanZhi))
			}
			if s.TianShen != "" {
				h.TianShen = s.TianShen
			}
			if s.Luck != "" {
				h.Luck = s.Luck
			}
			if len(s.Yi) > 0 {
				h.Yi = dropNone(s.Yi)
			}
			if len(s.Ji) > 0 {
				h.Ji = dropNone(s.Ji)
			}
			if len(s.JiShen) > 0 {
				h.JiShen = s.JiShen
			}
			if len(s.XiongSha) > 0 {
				h.XiongSha = s.XiongSha
			}
			hours[i] = h
		}
		out.Hours = hours
	} else {
		filled = append(filled, "hours")
	}

	out.ChongDesc = formatChongDesc(&out)
	out.Source = models.AlmanacSourceMerged
	out.AlgoFields = filled
	out.Mismatches = mismatches
	return &out, mismatches
}

func tianShenType(name string) string {
	switch name {
	case "青龙", "明堂", "金匮", "天德", "玉堂", "司命":
		return "黄道"
	case "天刑", "朱雀", "白虎", "天牢", "玄武", "勾陈":
		return "黑道"
	}
	return ""
}

func spans(sel *goquery.Selection) []string {
	out := make([]string, 0)
	sel.Find("span").Each(func(_ int, s *goquery.Selection) {
		if v := clean(s.Text()); v != "" {
			out = append(out, v)
		}
	})
	return out
}

func fields(s, sep string) []string {
	out := make([]string, 0)
	for _, p := range strings.Split(s, sep) {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func dropNone(arr []string) []string {
	out := make([]string, 0, len(arr))
	for _, v := range arr {
		if v != "无" {
			out = append(out, v)
		}
	}
	return out
}

func at(arr []string, i int) string {
	if i < len(arr) {
		return arr[i]
	}
	return ""
}

func sortRunes(s string) string {
	r := []rune(s)
	sort.Slice(r, func(i, j int) bool { return r[i] < r[j] })
	return string(r)
}

func clean(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
