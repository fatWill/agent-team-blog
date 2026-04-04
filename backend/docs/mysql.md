# MySQL 使用规范

## 连接配置

- **数据库**：`blog`
- **字符集**：`utf8mb4`
- **ORM**：GORM v1.31.1（gorm.io/driver/mysql）
- **连接池**：MaxOpenConns=10, MaxIdleConns=5（针对 2 核 2G 服务器优化）
- **日志级别**：Warn（仅记录慢查询和错误）

## 全量表结构

### articles（文章表）

```go
type Article struct {
    ID         string    `gorm:"column:id;primaryKey;type:varchar(36)"`       // UUID v4
    Title      string    `gorm:"column:title;type:varchar(255)"`              // 标题
    Summary    string    `gorm:"column:summary;type:varchar(500)"`            // 摘要
    Content    JSON      `gorm:"column:content;type:json"`                    // Tiptap JSON DSL
    CoverImage string    `gorm:"column:cover_image;type:varchar(500)"`        // 封面图 URL
    LikeCount  int       `gorm:"column:like_count;type:int;default:0"`        // 点赞计数（冗余字段）
    CreatedAt  time.Time `gorm:"column:created_at"`
    UpdatedAt  time.Time `gorm:"column:updated_at"`
}
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | varchar(36) | PK | UUID v4，由应用层生成 |
| title | varchar(255) | NOT NULL | 文章标题 |
| summary | varchar(500) | | 文章摘要 |
| content | json | | Tiptap 编辑器 JSON 内容 |
| cover_image | varchar(500) | | 封面图 URL |
| like_count | int | DEFAULT 0 | 点赞计数（冗余，与 article_likes 同步） |
| created_at | datetime | | 创建时间 |
| updated_at | datetime | | 更新时间 |

### article_likes（文章点赞记录）

```go
type ArticleLike struct {
    ID        uint64    `gorm:"primaryKey;autoIncrement"`
    ArticleID string    `gorm:"column:article_id;type:varchar(36)"`
    DeviceID  string    `gorm:"column:device_id;type:varchar(64)"`
    CreatedAt time.Time `gorm:"column:created_at"`
}
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | bigint unsigned | PK, AUTO_INCREMENT | |
| article_id | varchar(36) | NOT NULL | 关联文章 ID |
| device_id | varchar(64) | NOT NULL | 设备指纹 |
| created_at | datetime | | 点赞时间 |

**索引**：`UNIQUE (article_id, device_id)` — 每设备每文章只能点赞一次

### albums（相册表）

```go
type Album struct {
    ID           uint64  `gorm:"primaryKey;autoIncrement"`
    Name         string  `gorm:"column:name;type:varchar(100)"`
    Description  *string `gorm:"column:description;type:varchar(200)"`
    CoverURL     *string `gorm:"column:cover_url;type:varchar(500)"`
    PasswordHash *string `gorm:"column:password_hash;type:varchar(255)"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | bigint unsigned | PK, AUTO_INCREMENT | |
| name | varchar(100) | NOT NULL | 相册名称 |
| description | varchar(200) | NULLABLE | 相册描述 |
| cover_url | varchar(500) | NULLABLE | 封面 URL（自动取最新照片） |
| password_hash | varchar(255) | NULLABLE | bcrypt 密码哈希 |
| created_at | datetime | | |
| updated_at | datetime | | |

### photos（照片表）

```go
type Photo struct {
    ID           uint64  `gorm:"primaryKey;autoIncrement"`
    AlbumID      uint64  `gorm:"column:album_id"`
    URL          string  `gorm:"column:url;type:varchar(500)"`
    Caption      *string `gorm:"column:caption;type:varchar(200)"`
    PasswordHash *string `gorm:"column:password_hash;type:varchar(255)"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | bigint unsigned | PK, AUTO_INCREMENT | |
| album_id | bigint unsigned | NOT NULL | 所属相册 ID |
| url | varchar(500) | NOT NULL | 图片 URL |
| caption | varchar(200) | NULLABLE | 图片说明 |
| password_hash | varchar(255) | NULLABLE | bcrypt 密码哈希 |
| created_at | datetime | | |
| updated_at | datetime | | |

### photo_likes（照片点赞记录）

```go
type PhotoLike struct {
    ID       uint64    `gorm:"primaryKey;autoIncrement"`
    PhotoID  uint64    `gorm:"column:photo_id"`
    DeviceID string    `gorm:"column:device_id;type:varchar(64)"`
    CreatedAt time.Time
}
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | bigint unsigned | PK, AUTO_INCREMENT | |
| photo_id | bigint unsigned | NOT NULL | 照片 ID |
| device_id | varchar(64) | NOT NULL | 设备指纹 |
| created_at | datetime | | |

**幂等**：使用 `INSERT IGNORE` 保证每设备每照片只能点赞一次

### photo_dislikes（照片踩记录）

结构与 `photo_likes` 完全一致：

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | bigint unsigned | PK, AUTO_INCREMENT | |
| photo_id | bigint unsigned | NOT NULL | 照片 ID |
| device_id | varchar(64) | NOT NULL | 设备指纹 |
| created_at | datetime | | |

### profile（博主资料表）

```go
type Profile struct {
    ID     uint64 `gorm:"primaryKey"`
    Avatar string `gorm:"column:avatar;type:varchar(500)"`
    Bio    string `gorm:"column:bio;type:text"`
}
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | bigint unsigned | PK | 固定为 1（单行表） |
| avatar | varchar(500) | | 头像 URL |
| bio | text | | 个人简介 |

> 无 created_at / updated_at 字段（GORM 模型中未定义时间字段）。通过 `WHERE id = 1` 操作。

### messages（留言板表）

```go
type Message struct {
    ID               uint64  `gorm:"primaryKey;autoIncrement"`
    DeviceID         string  `gorm:"column:device_id;type:varchar(64);uniqueIndex"`
    Nickname         *string `gorm:"column:nickname;type:varchar(100)"`
    Content          string  `gorm:"column:content;type:text"`
    LastModifiedDate *string `gorm:"column:last_modified_date;type:date"`
    IP               *string `gorm:"column:ip;type:varchar(45)"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | bigint unsigned | PK, AUTO_INCREMENT | |
| device_id | varchar(64) | UNIQUE | 设备指纹（每设备一条留言） |
| nickname | varchar(100) | NULLABLE | 昵称（默认显示"匿名"） |
| content | text | NOT NULL | 留言内容 |
| last_modified_date | date | NULLABLE | 最后修改日期（每日限改一次） |
| ip | varchar(45) | NULLABLE | 客户端 IP |
| created_at | datetime | | |
| updated_at | datetime | | |

### changelogs（更新日志表）

```go
type Changelog struct {
    ID        uint64    `gorm:"primaryKey;autoIncrement"`
    Version   string    `gorm:"column:version;type:varchar(20);uniqueIndex"`
    Date      string    `gorm:"column:date;type:varchar(20)"`
    Logs      JSON      `gorm:"column:logs;type:json"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | bigint unsigned | PK, AUTO_INCREMENT | |
| version | varchar(20) | UNIQUE | 版本号（如 `1.9.1`，无 v 前缀） |
| date | varchar(20) | | 发版日期（`YYYY-MM-DD`） |
| logs | json | | 更新内容 JSON 数组 |
| created_at | datetime | | |
| updated_at | datetime | | |

## 查询规范

### 参数化查询（强制）

```go
// ✅ 正确
utils.DB.Where("id = ?", id).First(&article)
utils.DB.Where("title LIKE ?", "%"+keyword+"%").Find(&articles)

// ❌ 禁止拼接
utils.DB.Where(fmt.Sprintf("id = '%s'", id))
```

### 批量查询

IN 查询参数控制在 100 以内：

```go
utils.DB.Where("article_id IN ? AND device_id = ?", ids, deviceID).Find(&likes)
```

### 计数查询

使用 `Count` 方法：

```go
var count int64
utils.DB.Model(&models.Photo{}).Where("album_id = ?", albumID).Count(&count)
```

### 原生表达式

使用 `gorm.Expr` 处理数据库层计算：

```go
utils.DB.Model(&models.Article{}).Where("id = ?", id).
    Update("like_count", gorm.Expr("GREATEST(like_count - 1, 0)"))
```

### 更新操作

使用 `map[string]interface{}` 方式更新，支持零值字段：

```go
updates := map[string]interface{}{
    "nickname": nickname,    // 可能为 nil
    "content":  content,
}
utils.DB.Model(&msg).Updates(updates)
```

## 自定义 JSON 类型

`models/json.go` 定义了 `JSON` 类型，用于 MySQL JSON 字段的序列化/反序列化：

```go
type JSON json.RawMessage
```

实现了以下接口：
- `driver.Valuer` — 写入数据库
- `sql.Scanner` — 从数据库读取
- `json.Marshaler` — JSON 序列化
- `json.Unmarshaler` — JSON 反序列化

使用场景：`Article.Content`（Tiptap JSON）、`Changelog.Logs`（更新日志数组）
