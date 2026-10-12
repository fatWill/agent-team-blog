package almanac

import (
	"os"
	"strings"
	"testing"
	"time"
)

func mustDate(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := ParseDate(s)
	if err != nil {
		t.Fatalf("ParseDate(%s): %v", s, err)
	}
	return d
}

func TestComputeKnownDays(t *testing.T) {
	cases := []struct {
		date, gz, lunar, chong, sha string
		leap                        bool
		age                         int
	}{
		// 2024-01-01：癸卯年 甲子月 甲子日，冬月二十，冲马(戊午) 煞南；1978 戊午生人在 2023 癸卯年虚岁 46
		{"2024-01-01", "癸卯 甲子 甲子", "冬月二十", "马", "南", false, 46},
		// 2024-02-10：春节，甲辰年正月初一
		{"2024-02-10", "甲辰 丙寅 甲辰", "正月初一", "狗", "南", false, 0},
		// 2025-07-25：闰六月初一
		{"2025-07-25", "乙巳 癸未 乙未", "闰六月初一", "牛", "西", true, 0},
	}
	for _, c := range cases {
		d := Compute(mustDate(t, c.date))
		got := d.GanZhiYear + " " + d.GanZhiMonth + " " + d.GanZhiDay
		if got != c.gz || d.LunarDate != c.lunar || d.ChongZodiac != c.chong || d.Sha != c.sha || d.IsLeapMonth != c.leap {
			t.Errorf("%s: got gz=%s lunar=%s chong=%s sha=%s leap=%v", c.date, got, d.LunarDate, d.ChongZodiac, d.Sha, d.IsLeapMonth)
		}
		if c.age > 0 && d.ChongAge != c.age {
			t.Errorf("%s: chongAge=%d want %d", c.date, d.ChongAge, c.age)
		}
		if len(d.Hours) != 12 || d.Hours[0].Zhi != "子" || d.Hours[11].Zhi != "亥" {
			t.Errorf("%s: hours 异常 %d", c.date, len(d.Hours))
		}
		if d.Source != "lunar-go" || len(d.Yi) == 0 || d.PengZu == "" || d.TaiShen == "" {
			t.Errorf("%s: 字段缺失 %+v", c.date, d)
		}
	}
	if fs := Compute(mustDate(t, "2024-02-10")).Festivals; !contains(fs, "春节") {
		t.Errorf("2024-02-10 节日应含春节: %v", fs)
	}
	if jq := Compute(mustDate(t, "2024-02-04")).JieQi; jq != "立春" {
		t.Errorf("2024-02-04 应为立春, got %q", jq)
	}
}

func TestParseAndMergeSite(t *testing.T) {
	f, err := os.Open("testdata/huangli_2024-01-01.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	site, err := ParseSite(f, "2024-01-01")
	if err != nil {
		t.Fatalf("ParseSite: %v", err)
	}
	if site.BaZi != "癸卯 甲子 甲子" || site.LunarMonth != "冬月" || site.LunarDay != "二十" {
		t.Errorf("基础字段: %+v", site)
	}
	if strings.Join(site.Yi, ",") != "交易,进人口,祭祀,沐浴,捕捉,入殓,除服,成服,安葬" {
		t.Errorf("宜: %v", site.Yi)
	}
	if len(site.Ji) != 5 || len(site.JiShen) != 7 || len(site.XiongSha) != 5 {
		t.Errorf("忌/吉神/凶煞数量: %v %v %v", site.Ji, site.JiShen, site.XiongSha)
	}
	if site.ZhiShen != "金匮" || site.NaYinDay != "海中金" || site.JianChu != "建" || site.ChongZodiac != "马" || site.Sha != "南" {
		t.Errorf("底部字段: %+v", site)
	}
	if site.Xiu != "毕乌月" && site.Xiu != "毕月乌" {
		t.Errorf("星宿: %q", site.Xiu)
	}
	if site.PengZu != "甲不开仓 子不问卜" || site.TaiShen != "占门碓 外东南" {
		t.Errorf("彭祖/胎神: %q %q", site.PengZu, site.TaiShen)
	}
	if len(site.Hours) != 12 || site.Hours[0].GanZhi != "甲子" || site.Hours[6].Luck != "吉" || site.Hours[2].TianShen != "白虎" {
		t.Fatalf("时辰表: %+v", site.Hours)
	}

	merged, mm := Merge(Compute(mustDate(t, "2024-01-01")), site)
	if len(mm) != 0 {
		t.Errorf("确定性字段不应有不一致: %v", mm)
	}
	if merged.Source != "huangli.com+lunar-go" || len(merged.Yi) != 9 || merged.PengZu != "甲不开仓财物耗散 子不问卜自惹祸殃" || merged.Xiu != "毕月乌" {
		t.Errorf("合并结果: %+v", merged)
	}
	if contains(merged.Hours[6].Yi, "无") {
		t.Errorf("时宜应剔除“无”: %v", merged.Hours[6].Yi)
	}
}

func TestParseSiteLeapMonthAndDateGuard(t *testing.T) {
	b, err := os.ReadFile("testdata/huangli_2025-07-25.html")
	if err != nil {
		t.Fatal(err)
	}
	site, err := ParseSite(strings.NewReader(string(b)), "2025-07-25")
	if err != nil {
		t.Fatal(err)
	}
	_, mm := Merge(Compute(mustDate(t, "2025-07-25")), site)
	if len(mm) != 0 {
		t.Errorf("闰月日不应有不一致: %v", mm)
	}
	if _, err := ParseSite(strings.NewReader(string(b)), "2025-07-26"); err == nil {
		t.Error("页面日期与请求日期不一致时应报错")
	}
	if _, err := ParseSite(strings.NewReader("<html>blocked</html>"), "2025-07-25"); err == nil {
		t.Error("非黄历页面应报错")
	}
}

func TestMergeRecordsMismatch(t *testing.T) {
	calc := Compute(mustDate(t, "2024-01-01"))
	site := &SiteDay{BaZi: "癸卯 甲子 乙丑", ChongZodiac: "羊", Yi: []string{"祭祀"}}
	out, mm := Merge(calc, site)
	if len(mm) != 2 {
		t.Fatalf("应记录 2 项不一致, got %v", mm)
	}
	if out.GanZhiDay != "甲子" {
		t.Errorf("确定性字段应以算法为准, got %s", out.GanZhiDay)
	}
	if strings.Join(out.Yi, "") != "祭祀" || !contains(out.AlgoFields, "hours") || !contains(out.AlgoFields, "pengZu") {
		t.Errorf("站点优先 / 算法补齐登记异常: yi=%v algo=%v", out.Yi, out.AlgoFields)
	}
}

func TestBirthAndHelpers(t *testing.T) {
	b := ComputeBirth(mustDate(t, "1990-01-20"), -1, "solar")
	// 1990-01-20 在立春前、春节前：生肖（春节）蛇，八字年柱己巳
	if b.ZodiacSpring != "蛇" || b.YearPillar != "己巳" || b.TimePillar != "" {
		t.Errorf("birth: %+v", b)
	}
	if b2 := ComputeBirth(mustDate(t, "1990-01-20"), 9, "solar"); b2.TimePillar == "" {
		t.Error("给出时辰应返回时柱")
	}
	if y := ComputeYear(2027); y.GanZhi != "丁未" || y.Zodiac != "羊" {
		t.Errorf("year: %+v", y)
	}
	s, err := LunarToSolar(2025, 6, 1, true)
	if err != nil || s.Format(DateLayout) != "2025-07-25" {
		t.Errorf("闰六月初一应为 2025-07-25: %v %v", s, err)
	}
	if _, err := LunarToSolar(2024, 6, 1, true); err == nil {
		t.Error("2024 无闰六月应报错")
	}
	if r := ZodiacRelation("鼠", "马"); r[0] != "相冲（冲太岁）" {
		t.Errorf("鼠马应相冲: %v", r)
	}
	if _, err := ParseDate("2024-13-01"); err == nil {
		t.Error("非法日期应报错")
	}
}
