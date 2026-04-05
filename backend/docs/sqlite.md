# SQLite 使用规范

## 连接配置

- **驱动**：`github.com/glebarez/sqlite`（基于 `modernc.org/sqlite` 纯 Go 实现，无 CGO）
- **ORM**：GORM v1.31.1
- **数据库文件**：通过 `DB_PATH` 环境变量配置，默认 `/root/blog-data/blog.db`
- **连接池**：MaxOpenConns=1, MaxIdleConns=1（SQLite 单文件数据库，并发写入受限）
- **日志级别**：Warn（仅记录慢查询和错误）

## PRAGMA 优化

```sql
PRAGMA journal_mode=WAL;     -- WAL 模式，支持并发读
PRAGMA synchronous=NORMAL;   -- 平衡性能与安全
PRAGMA cache_size=-8000;     -- 8MB 缓存
PRAGMA busy_timeout=5000;    -- 锁等待 5 秒
PRAGMA foreign_keys=ON;      -- 启用外键约束
```

## 自动建表

首次连接时自动创建所有表（`CREATE TABLE IF NOT EXISTS`），无需手动执行 DDL。建表逻辑在 `pkg/db/db.go` 的 `autoMigrate()` 函数中。

## 全量表结构

### articles（文章表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | TEXT | PK | UUID v4，由应用层生成 |
| title | TEXT | NOT NULL DEFAULT '' | 文章标题 |
| summary | TEXT | NOT NULL DEFAULT '' | 文章摘要 |
| content | TEXT | | Tiptap 编辑器 JSON 内容 |
| cover_image | TEXT | NOT NULL DEFAULT '' | 封面图 URL |
| like_count | INTEGER | NOT NULL DEFAULT 0 | 点赞计数（冗余，与 article_likes 同步） |
| created_at | DATETIME | NOT NULL DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | DATETIME | NOT NULL DEFAULT CURRENT_TIMESTAMP | 更新时间 |

### article_likes（文章点赞记录）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | INTEGER | PK AUTOINCREMENT | |
| article_id | TEXT | NOT NULL DEFAULT '' | 关联文章 ID |
| device_id | TEXT | NOT NULL DEFAULT '' | 设备指纹 |
| created_at | DATETIME | NOT NULL DEFAULT CURRENT_TIMESTAMP | 点赞时间 |

**索引**：`UNIQUE (article_id, device_id)` — 每设备每文章只能点赞一次

### albums（相册表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | INTEGER | PK AUTOINCREMENT | |
| name | TEXT | NOT NULL DEFAULT '' | 相册名称 |
| description | TEXT | NULLABLE | 相册描述 |
| cover_url | TEXT | NULLABLE | 封面 URL（自动取最新照片） |
| password_hash | TEXT | NULLABLE | bcrypt 密码哈希 |
| created_at | DATETIME | NOT NULL DEFAULT CURRENT_TIMESTAMP | |
| updated_at | DATETIME | NOT NULL DEFAULT CURRENT_TIMESTAMP | |

### photos（照片表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | INTEGER | PK AUTOINCREMENT | |
| album_id | INTEGER | NOT NULL DEFAULT 0 | 所属相册 ID |
| url | TEXT | NOT NULL DEFAULT '' | 图片 URL |
| caption | TEXT | NULLABLE | 图片说明 |
| password_hash | TEXT | NULLABLE | bcrypt 密码哈希 |
| created_at | DATETIME | NOT NULL DEFAULT CURRENT_TIMESTAMP | |
| updated_at | DATETIME | NOT NULL DEFAULT CURRENT_TIMESTAMP | |

**索引**：`idx_photos_album_id (album_id)`

### photo_likes（照片点赞记录）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | INTEGER | PK AUTOINCREMENT | |
| photo_id | INTEGER | NOT NULL DEFAULT 0 | 照片 ID |
| device_id | TEXT | NOT NULL DEFAULT '' | 设备指纹 |
| created_at | DATETIME | NOT NULL DEFAULT CURRENT_TIMESTAMP | |

**索引**：`UNIQUE (photo_id, device_id)` — 幂等，使用 `INSERT OR IGNORE`

### photo_dislikes（照片踩记录）

结构与 `photo_likes` 完全一致。

**索引**：`UNIQUE (photo_id, device_id)`

### profile（博主资料表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | INTEGER | PK | 固定为 1（单行表） |
| avatar | TEXT | NOT NULL DEFAULT '' | 头像 URL |
| bio | TEXT | NOT NULL DEFAULT '' | 个人简介 |

> 无 created_at / updated_at 字段。初始化时自动插入默认行：`INSERT OR IGNORE INTO profile (id, avatar, bio) VALUES (1, '', '')`

### messages（留言板表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | INTEGER | PK AUTOINCREMENT | |
| device_id | TEXT | NOT NULL DEFAULT '' | 设备指纹（每设备一条留言） |
| nickname | TEXT | NULLABLE | 昵称（默认显示"匿名"） |
| content | TEXT | NOT NULL DEFAULT '' | 留言内容 |
| last_modified_date | TEXT | NULLABLE | 最后修改日期（YYYY-MM-DD，每日限改一次） |
| ip | TEXT | NULLABLE | 客户端 IP |
| created_at | DATETIME | NOT NULL DEFAULT CURRENT_TIMESTAMP | |
| updated_at | DATETIME | NOT NULL DEFAULT CURRENT_TIMESTAMP | |

**索引**：`UNIQUE (device_id)`

### changelogs（更新日志表）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | INTEGER | PK AUTOINCREMENT | |
| version | TEXT | NOT NULL DEFAULT '' | 版本号（如 `1.9.1`，无 v 前缀） |
| date | TEXT | NOT NULL DEFAULT '' | 发版日期（`YYYY-MM-DD`） |
| logs | TEXT | | 更新内容 JSON 数组 |
| created_at | DATETIME | NOT NULL DEFAULT CURRENT_TIMESTAMP | |
| updated_at | DATETIME | NOT NULL DEFAULT CURRENT_TIMESTAMP | |

**索引**：`UNIQUE (version)`

## SQL 兼容性注意事项

### 与 MySQL 的差异

| MySQL 语法 | SQLite 替代 | 说明 |
|------------|-------------|------|
| `AUTO_INCREMENT` | `AUTOINCREMENT` | 自增主键 |
| `INSERT IGNORE` | `INSERT OR IGNORE` | 忽略唯一冲突 |
| `GREATEST(a, b)` | `MAX(a, b)` | 取最大值 |
| `TINYINT(1)` | `INTEGER` | SQLite 无 TINYINT |
| `type:json` | `type:text` | SQLite 无原生 JSON 类型，用 TEXT 存储 |
| `type:varchar(N)` | `TEXT` | SQLite 文本统一为 TEXT |
| `type:date` | `TEXT` | 日期用 TEXT 存储 |

### 参数化查询（强制）

```go
// ✅ 正确
db.DB.Where("id = ?", id).First(&article)
db.DB.Where("title LIKE ?", "%"+keyword+"%").Find(&articles)

// ❌ 禁止拼接
db.DB.Where(fmt.Sprintf("id = '%s'", id))
```

### 幂等插入

```go
// SQLite 使用 INSERT OR IGNORE
db.DB.Exec("INSERT OR IGNORE INTO photo_likes (photo_id, device_id) VALUES (?, ?)", photoID, deviceID)
```

### 计数安全递减

```go
// SQLite 使用 MAX 替代 MySQL 的 GREATEST
db.DB.Model(&models.Article{}).Where("id = ?", articleID).
    Update("like_count", gorm.Expr("MAX(like_count - 1, 0)"))
```

## 自定义 JSON 类型

`models/json.go` 定义了 `JSON` 类型，用于 SQLite TEXT 字段的 JSON 序列化/反序列化：

```go
type JSON json.RawMessage
```

实现了以下接口：
- `driver.Valuer` — 写入数据库（转为 string）
- `sql.Scanner` — 从数据库读取（支持 []byte 和 string）
- `json.Marshaler` — JSON 序列化
- `json.Unmarshaler` — JSON 反序列化

使用场景：`Article.Content`（Tiptap JSON）、`Changelog.Logs`（更新日志数组）
