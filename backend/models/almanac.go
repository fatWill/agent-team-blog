package models

import "time"

// AlmanacSource 黄历数据来源标识
const (
	// AlmanacSourceAlgo 纯算法生成（lunar-go），未抓取到站点数据
	AlmanacSourceAlgo = "lunar-go"
	// AlmanacSourceMerged 站点抓取为主，缺失字段由 lunar-go 补齐
	AlmanacSourceMerged = "huangli.com+lunar-go"
)

// AlmanacHour 单个时辰的吉凶信息（12 时辰，子时为 23:00-00:59）
type AlmanacHour struct {
	Zhi      string   `json:"zhi"`      // 地支：子/丑/...
	GanZhi   string   `json:"ganZhi"`   // 时柱干支：甲子
	Range    string   `json:"range"`    // 时间范围：23:00-00:59
	TianShen string   `json:"tianShen"` // 星神（值神）：金匮
	Luck     string   `json:"luck"`     // 吉 / 凶
	Chong    string   `json:"chong"`    // 冲生肖：马
	Sha      string   `json:"sha"`      // 煞方：南
	Yi       []string `json:"yi"`       // 时宜
	Ji       []string `json:"ji"`       // 时忌
	JiShen   []string `json:"jiShen"`   // 时辰吉神
	XiongSha []string `json:"xiongSha"` // 时辰凶煞
}

// AlmanacDay 老黄历日记录（每天一条，date 唯一）
type AlmanacDay struct {
	ID   uint64 `json:"-" gorm:"primaryKey;autoIncrement"`
	Date string `json:"date" gorm:"column:date;uniqueIndex"` // 公历日期 YYYY-MM-DD

	Week         string `json:"week" gorm:"column:week"`                   // 星期一
	LunarYear    int    `json:"lunarYear" gorm:"column:lunar_year"`        // 农历年（数字）
	LunarMonth   int    `json:"lunarMonth" gorm:"column:lunar_month"`      // 农历月（数字，闰月为正数，配合 isLeapMonth）
	LunarDay     int    `json:"lunarDay" gorm:"column:lunar_day"`          // 农历日（数字）
	IsLeapMonth  bool   `json:"isLeapMonth" gorm:"column:is_leap_month"`   // 是否闰月
	LunarMonthCN string `json:"lunarMonthCn" gorm:"column:lunar_month_cn"` // 冬月 / 闰六月
	LunarDayCN   string `json:"lunarDayCn" gorm:"column:lunar_day_cn"`     // 二十 / 初一
	LunarDate    string `json:"lunarDate" gorm:"column:lunar_date"`        // 冬月二十 / 闰六月初一
	LunarYearCN  string `json:"lunarYearCn" gorm:"column:lunar_year_cn"`   // 二〇二三
	ZodiacYear   string `json:"zodiacYear" gorm:"column:zodiac_year"`      // 生肖年：兔
	GanZhiYear   string `json:"ganZhiYear" gorm:"column:gz_year"`          // 年柱：癸卯
	GanZhiMonth  string `json:"ganZhiMonth" gorm:"column:gz_month"`        // 月柱：甲子
	GanZhiDay    string `json:"ganZhiDay" gorm:"column:gz_day"`            // 日柱：甲子
	JieQi        string `json:"jieQi" gorm:"column:jieqi"`                 // 当天交节气（无则空）
	CurrentJieQi string `json:"currentJieQi" gorm:"column:current_jieqi"`  // 所处节气区间：冬至

	Festivals []string `json:"festivals" gorm:"column:festivals;serializer:json"` // 节日（公历+农历+其他）
	Yi        []string `json:"yi" gorm:"column:yi;serializer:json"`               // 宜
	Ji        []string `json:"ji" gorm:"column:ji;serializer:json"`               // 忌

	ChongZodiac string `json:"chongZodiac" gorm:"column:chong_zodiac"` // 冲生肖：马
	ChongGanZhi string `json:"chongGanZhi" gorm:"column:chong_gz"`     // 冲的干支：戊午
	ChongAge    int    `json:"chongAge" gorm:"column:chong_age"`       // 被冲者虚岁（按当年农历年计算）
	ChongDesc   string `json:"chongDesc" gorm:"column:chong_desc"`     // 冲马(戊午)47岁
	Sha         string `json:"sha" gorm:"column:sha"`                  // 煞方：南

	NaYinYear   string `json:"naYinYear" gorm:"column:nayin_year"`     // 年纳音
	NaYinMonth  string `json:"naYinMonth" gorm:"column:nayin_month"`   // 月纳音
	NaYinDay    string `json:"naYinDay" gorm:"column:nayin_day"`       // 日纳音（五行）：海中金
	ZhiShen     string `json:"zhiShen" gorm:"column:zhishen"`          // 值神：金匮
	ZhiShenType string `json:"zhiShenType" gorm:"column:zhishen_type"` // 黄道 / 黑道
	JianChu     string `json:"jianChu" gorm:"column:jianchu"`          // 建除十二神：建
	Xiu         string `json:"xiu" gorm:"column:xiu"`                  // 二十八星宿：毕月乌
	XiuLuck     string `json:"xiuLuck" gorm:"column:xiu_luck"`         // 星宿吉凶
	XiuGong     string `json:"xiuGong" gorm:"column:xiu_gong"`         // 星宿所属方位：西方白虎
	PengZu      string `json:"pengZu" gorm:"column:pengzu"`            // 彭祖百忌

	JiShen   []string `json:"jiShen" gorm:"column:jishen;serializer:json"`     // 吉神宜趋
	XiongSha []string `json:"xiongSha" gorm:"column:xiongsha;serializer:json"` // 凶神宜忌

	TaiShen  string `json:"taiShen" gorm:"column:taishen"`   // 胎神占方：占门碓 外东南
	XiShen   string `json:"xiShen" gorm:"column:xishen"`     // 喜神方位：东北
	CaiShen  string `json:"caiShen" gorm:"column:caishen"`   // 财神方位：东北
	FuShen   string `json:"fuShen" gorm:"column:fushen"`     // 福神方位：正北
	YangGui  string `json:"yangGui" gorm:"column:yanggui"`   // 阳贵神方位
	YinGui   string `json:"yinGui" gorm:"column:yingui"`     // 阴贵神方位
	LiuYao   string `json:"liuYao" gorm:"column:liuyao"`     // 六曜
	WuHou    string `json:"wuHou" gorm:"column:wuhou"`       // 物候
	YueXiang string `json:"yueXiang" gorm:"column:yuexiang"` // 月相

	Hours []AlmanacHour `json:"hours" gorm:"column:hours;serializer:json"` // 12 时辰吉凶

	Source     string   `json:"source" gorm:"column:source"`                          // 数据来源
	AlgoFields []string `json:"algoFields" gorm:"column:algo_fields;serializer:json"` // 由算法补齐的字段
	Mismatches []string `json:"mismatches" gorm:"column:mismatches;serializer:json"`  // 站点与算法不一致的字段（以站点为准）

	CreatedAt time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

// TableName 表名
func (AlmanacDay) TableName() string { return "almanac_days" }

// AlmanacMonthItem 月视图精简字段
type AlmanacMonthItem struct {
	Date        string   `json:"date"`
	Week        string   `json:"week"`
	LunarDate   string   `json:"lunarDate"`
	CellText    string   `json:"cellText"` // 日格显示用：节气 > 初一显示月名 > 农历日名
	GanZhiDay   string   `json:"ganZhiDay"`
	JieQi       string   `json:"jieQi"`
	Festivals   []string `json:"festivals"`
	ZhiShenType string   `json:"zhiShenType"`
	YiTop       []string `json:"yiTop"` // 宜（前 3 项）
	JiTop       []string `json:"jiTop"` // 忌（前 3 项）
	ChongZodiac string   `json:"chongZodiac"`
	Source      string   `json:"source"`
}
