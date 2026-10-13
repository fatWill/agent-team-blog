package tuyile

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fatWill/agent-team-blog/backend/internal/almanac"
	"github.com/fatWill/agent-team-blog/backend/models"
	algo "github.com/fatWill/agent-team-blog/backend/pkg/almanac"
)

const toolQueryAlmanac = "query_almanac"

// almanacTool query_almanac 工具定义（OpenAI function calling）
var almanacTool = func() toolDef {
	var t toolDef
	t.Type = "function"
	t.Function.Name = toolQueryAlmanac
	t.Function.Description = "查询老黄历（权威知识库）。凡是涉及具体日期/日期范围的宜忌、吉日、冲煞、时辰吉凶、喜神财神福神方位，都必须先调用本工具，不得凭记忆编造。" +
		"支持单日（只传 start）或日期范围（start~end，闭区间，最多 62 天）。可用 keywords 按「宜」过滤吉日，如 [\"搬家\"]、[\"开业\"]、[\"结婚\"]。" +
		"单日查询返回完整信息（含 12 时辰吉凶）；范围查询返回每日精简信息，需要某天时辰时再单日查询。"
	t.Function.Parameters = json.RawMessage(`{
  "type": "object",
  "properties": {
    "start": {"type": "string", "description": "开始日期（公历）YYYY-MM-DD"},
    "end": {"type": "string", "description": "结束日期（公历）YYYY-MM-DD，可省略=单日；与 start 间隔不超过 62 天"},
    "keywords": {"type": "array", "items": {"type": "string"}, "description": "按宜过滤的事项，如 搬家/入宅/开业/结婚/出行/签约/动土/求职；省略则不过滤"}
  },
  "required": ["start"]
}`)
	return t
}()

type almanacArgs struct {
	Start    string   `json:"start"`
	End      string   `json:"end"`
	Keywords []string `json:"keywords"`
}

// yiAliases 口语 → 黄历宜项（任一命中即视为宜）
var yiAliases = map[string][]string{
	"搬家": {"移徙", "入宅"}, "乔迁": {"移徙", "入宅"}, "入宅": {"入宅"}, "入住": {"入宅", "移徙"}, "移徙": {"移徙"},
	"开业": {"开市"}, "开张": {"开市"}, "开店": {"开市"}, "开市": {"开市"}, "开公司": {"开市", "立券"}, "创业": {"开市", "立券"},
	"结婚": {"嫁娶"}, "嫁娶": {"嫁娶"}, "领证": {"嫁娶", "订盟"}, "婚礼": {"嫁娶"}, "订婚": {"订盟", "纳采"}, "提亲": {"纳采"},
	"出行": {"出行"}, "旅游": {"出行"}, "出差": {"出行"}, "出门": {"出行"},
	"签约": {"立券", "交易"}, "签合同": {"立券", "交易"}, "交易": {"交易"}, "买房": {"置产", "交易", "立券"}, "买车": {"交易", "纳财"}, "提车": {"交易", "纳财"},
	"装修": {"修造", "动土"}, "动土": {"动土"}, "修造": {"修造"}, "安床": {"安床"},
	"入职": {"出行", "开市", "赴任"}, "上任": {"赴任", "出行"}, "求职": {"出行", "开市"}, "面试": {"出行"},
	"理发": {"理发"}, "剪头": {"理发"}, "祭祀": {"祭祀"}, "祈福": {"祈福"}, "求财": {"纳财", "开市"}, "纳财": {"纳财"},
	"求医": {"求医", "治病"}, "入学": {"入学"}, "考试": {"入学", "祈福"}, "安葬": {"安葬"}, "栽种": {"栽种"},
}

// expandKeywords 把口语关键词映射为黄历宜项；未登记的词原样作为宜项
func expandKeywords(kws []string) []string {
	out := make([]string, 0, len(kws)*2)
	seen := map[string]bool{}
	for _, k := range kws {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		list, ok := yiAliases[k]
		if !ok {
			list = []string{k}
		}
		for _, x := range list {
			if !seen[x] {
				seen[x] = true
				out = append(out, x)
			}
		}
	}
	return out
}

func matchYi(d *models.AlmanacDay, want []string) []string {
	hit := make([]string, 0, 2)
	for _, w := range want {
		for _, y := range d.Yi {
			if y == w {
				hit = append(hit, w)
				break
			}
		}
	}
	return hit
}

// intentKeywords 从用户问题中识别择日意图（命中 yiAliases 的口语词），用于范围注入时预筛选，按出现顺序去重
func intentKeywords(text string) []string {
	type hit struct {
		pos int
		kw  string
	}
	hits := make([]hit, 0, 2)
	for k := range yiAliases {
		if i := strings.Index(text, k); i >= 0 {
			hits = append(hits, hit{i, k})
		}
	}
	// 位置升序；同位置优先更长的词（如「开公司」优先于「开」开头的短词）
	for i := 1; i < len(hits); i++ {
		for j := i; j > 0 && (hits[j].pos < hits[j-1].pos || (hits[j].pos == hits[j-1].pos && len(hits[j].kw) > len(hits[j-1].kw))); j-- {
			hits[j], hits[j-1] = hits[j-1], hits[j]
		}
	}
	out := make([]string, 0, len(hits))
	for _, h := range hits {
		out = append(out, h.kw)
	}
	return out
}

// renderRange 范围黄历文本：有择日关键词时只列命中的日子（与 query_almanac 同口径），否则列每日精简
func renderRange(repo *almanac.Repo, label string, s, e time.Time, keywords []string) string {
	args, _ := json.Marshal(almanacArgs{Start: s.Format(algo.DateLayout), End: e.Format(algo.DateLayout), Keywords: keywords})
	return "\n【用户提到的「" + label + "」】" + strings.TrimPrefix(runAlmanacTool(repo, string(args)), "【黄历查询")
}

// runAlmanacTool 执行 query_almanac，返回给模型的纯文本结果（同时用于兜底注入）
func runAlmanacTool(repo *almanac.Repo, rawArgs string) string {
	var a almanacArgs
	if err := json.Unmarshal([]byte(rawArgs), &a); err != nil {
		return "参数错误：需要 JSON {\"start\":\"YYYY-MM-DD\",\"end\":\"YYYY-MM-DD\",\"keywords\":[...]}"
	}
	s, err := algo.ParseDate(strings.TrimSpace(a.Start))
	if err != nil {
		return "参数错误：start " + err.Error()
	}
	e := s
	if strings.TrimSpace(a.End) != "" {
		if e, err = algo.ParseDate(strings.TrimSpace(a.End)); err != nil {
			return "参数错误：end " + err.Error()
		}
	}
	if e.Before(s) {
		s, e = e, s
	}
	truncated := false
	if e.Sub(s) > time.Duration(MaxRangeDays-1)*24*time.Hour {
		e, truncated = s.AddDate(0, 0, MaxRangeDays-1), true
	}
	days, err := repo.RangeOrCompute(s, e)
	if err != nil {
		return "查询失败：黄历库暂不可用"
	}
	want := expandKeywords(a.Keywords)

	var b strings.Builder
	if len(days) == 1 && len(want) == 0 {
		b.WriteString(almanac.FormatForPrompt("查询结果", days[0]))
		b.WriteString(formatHours(days[0]))
		return b.String()
	}
	fmt.Fprintf(&b, "【黄历查询 %s ~ %s，共 %d 天】", s.Format(algo.DateLayout), e.Format(algo.DateLayout), len(days))
	if truncated {
		fmt.Fprintf(&b, "（超过 %d 天已截断，如需后续日期请再次查询）", MaxRangeDays)
	}
	if len(want) > 0 {
		fmt.Fprintf(&b, "\n筛选：宜包含「%s」任一项", strings.Join(want, "/"))
	}
	b.WriteString("\n")
	matched := 0
	for _, d := range days {
		hit := matchYi(d, want)
		if len(want) > 0 && len(hit) == 0 {
			continue
		}
		matched++
		b.WriteString(formatCompact(d))
		if len(hit) > 0 {
			fmt.Fprintf(&b, " ✅命中:%s", strings.Join(hit, "、"))
		}
		b.WriteString("\n")
	}
	if len(want) > 0 {
		if matched == 0 {
			b.WriteString("该范围内没有宜上述事项的日子（请如实告诉用户，可建议换个时间段或看看宜忌相近的日子）。\n")
		} else {
			fmt.Fprintf(&b, "共 %d 天符合。\n", matched)
		}
	}
	return b.String()
}

// formatCompact 单日一行精简信息（范围查询 / 未来 30 天表）
func formatCompact(d *models.AlmanacDay) string {
	ji := make([]string, 0, 6)
	for _, h := range d.Hours {
		if h.Luck == "吉" {
			ji = append(ji, h.Zhi)
		}
	}
	extra := ""
	if d.JieQi != "" {
		extra += " 节气:" + d.JieQi
	}
	if len(d.Festivals) > 0 {
		extra += " 节日:" + strings.Join(d.Festivals, "、")
	}
	return fmt.Sprintf("%s %s 农历%s %s日 %s(%s)%s｜宜:%s｜忌:%s｜%s｜吉时:%s｜喜神%s 财神%s 福神%s",
		d.Date, strings.TrimPrefix(d.Week, "星期"), d.LunarDate, d.GanZhiDay, d.ZhiShen, d.ZhiShenType, extra,
		strings.Join(d.Yi, " "), strings.Join(d.Ji, " "), d.ChongDesc, strings.Join(ji, ""), d.XiShen, d.CaiShen, d.FuShen)
}

// formatHours 12 时辰明细（单日查询用）
func formatHours(d *models.AlmanacDay) string {
	var b strings.Builder
	b.WriteString("十二时辰：\n")
	for _, h := range d.Hours {
		fmt.Fprintf(&b, "- %s时 %s %s %s(%s) 冲%s 煞%s；宜:%s；忌:%s\n", h.Zhi, h.Range, h.GanZhi, h.TianShen, h.Luck, h.Chong, h.Sha,
			strings.Join(h.Yi, " "), strings.Join(h.Ji, " "))
	}
	return b.String()
}
