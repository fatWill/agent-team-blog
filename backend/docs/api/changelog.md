# 更新日志接口文档

## [GET] /api/changelog - 获取更新日志

### 描述
获取所有版本更新日志，按 id 倒序（最新版本在前）。

### 成功响应 (200)

```json
{
  "changelog": [
    {
      "version": "2.0.3",
      "date": "2026-04-03",
      "logs": [
        "🗑️ 删除顶部重复Tab导航栏",
        "📝 保留各Tab大字号标题",
        "📐 优化头部与内容间距"
      ],
      "createdAt": "2026-04-03T12:00:00+08:00",
      "updatedAt": "2026-04-03T12:00:00+08:00"
    },
    {
      "version": "2.0.2",
      "date": "2026-04-02",
      "logs": [
        "🏷️ Tab导航添加emoji图标",
        "💻 PC端激活态底部指示线",
        "📱 移动端激活态左侧竖线"
      ],
      "createdAt": "2026-04-02T12:00:00+08:00",
      "updatedAt": "2026-04-02T12:00:00+08:00"
    }
  ]
}
```

### 数据结构

| 字段 | 类型 | 说明 |
|------|------|------|
| version | string | 版本号（无 v 前缀，如 `2.0.3`） |
| date | string | 发版日期（`YYYY-MM-DD`） |
| logs | string[] | 更新内容数组（每条 ≤20 字，最多 5 条，可用 Emoji 前缀） |
| createdAt | string | 记录创建时间 |
| updatedAt | string | 记录更新时间 |

### 说明

- 此接口为**只读**接口，后端未提供 changelog 写入 API
- `logs` 字段在 SQLite 中存储为 `TEXT`，内容是 JSON 数组字符串，通过自定义 `models.JSON` 类型完成序列化/反序列化
- 版本号格式规范：`X.Y.Z`，无 `v` 前缀
- 排序为 `id DESC`，即**按插入顺序倒序**而非按 `date` 排序；补录历史版本时需注意展示顺序

### 更新日志的写入方式（2026-09-16 起）

新版本条目通过**启动幂等播种**写入，不再手工连生产库执行 SQL：

1. 在 `pkg/db/changelog_seed.go` 的 `changelogSeeds` 切片中追加一条条目
2. 构建并部署新二进制，服务启动时 `autoMigrate()` 会调用 `seedChangelogs()`
3. 依赖 `uk_changelogs_version` 唯一索引 + `INSERT OR IGNORE` 保证幂等：已存在同 `version` 的行不会被覆盖，缺失则补齐

> 兜底方案：`scripts/seed_changelog_2.15.0.sql` 可在不发布二进制时人工补数。因 SQLite 为 WAL 模式且服务连接池 `MaxOpenConns=1`，人工执行前须先停服务并备份数据库文件。

> 注意：`frontend/server/data/changelog.json` 不是线上数据源，仅作前端本地兜底/内容底稿；线上 `/api/changelog` 只读 SQLite `changelogs` 表。两处内容需人工保持一致。

