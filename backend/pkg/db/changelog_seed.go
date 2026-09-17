package db

import "fmt"

// changelogSeedItem 更新日志播种条目
// logs 为 JSON 数组字符串，直接写入 changelogs.logs（TEXT）字段
type changelogSeedItem struct {
	Version string
	Date    string
	Logs    string
}

// changelogSeeds 需要保证存在于 changelogs 表中的更新日志条目。
//
// 背景：后端没有 changelog 写入接口，历史条目均由人工写库。为避免每次发版都要
// 运维手工连生产库执行 SQL（SQLite 单文件 + WAL 模式，服务运行中外部写入有风险），
// 这里改为随服务启动做一次幂等播种：靠 changelogs.version 的唯一索引 +
// INSERT OR IGNORE 保证「已存在则不动、缺失则补齐」，可安全重复执行。
//
// 新增版本时只需在此切片追加一条，无需再写一次性 SQL 脚本。
var changelogSeeds = []changelogSeedItem{
	{
		Version: "2.15.0",
		Date:    "2026-09-16",
		Logs: `["🌐 站点主域切换到 fatwill.cn",` +
			`"🖼️ 图片与静态资源域名同步切换为 assets.fatwill.cn / cdn.fatwill.cn"]`,
	},
	{
		Version: "2.15.1",
		Date:    "2026-09-17",
		Logs: `["🖼️ 13 篇早期文章的封面与正文插图从服务器本地直供迁移至腾讯云 COS，并接入 EdgeOne CDN",` +
			`"⚡ 全站图片资源统一走 assets.fatwill.cn 加速域名，首屏与图片加载体验一致提升"]`,
	},
}

// seedChangelogs 幂等补齐 changelogSeeds 中声明的更新日志条目。
//
// 依赖 uk_changelogs_version 唯一索引实现幂等：已存在同 version 的行会被忽略，
// 不会覆盖线上已有内容（避免误改人工在库里修订过的文案）。
// created_at / updated_at 显式对齐发布日期，保证与 date 字段语义一致。
func seedChangelogs() error {
	for _, item := range changelogSeeds {
		err := DB.Exec(
			`INSERT OR IGNORE INTO changelogs (version, date, logs, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?)`,
			item.Version, item.Date, item.Logs,
			item.Date+" 00:00:00", item.Date+" 00:00:00",
		).Error
		if err != nil {
			return fmt.Errorf("播种更新日志 %s 失败: %w", item.Version, err)
		}
	}
	return nil
}
