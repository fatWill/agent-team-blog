# Go 编码规范

本文档定义项目的 Go 编码规范，所有后端代码必须遵循。

## 命名约定

| 类型 | 规则 | 示例 |
|------|------|------|
| 包名 | 小写单词，不用下划线和复数 | `handlers`, `models`, `middleware` |
| 文件名 | 小写 + 下划线（snake_case） | `articles.go`, `album.go`, `ratelimit.go` |
| 结构体 | PascalCase | `Article`, `AlbumListItem`, `PhotoDislike` |
| 方法/函数 | PascalCase（导出）/ camelCase（私有） | `GetArticles`, `getTodayCST`, `updateAlbumCover` |
| 常量 | PascalCase | `TokenTTL`, `TokenKeyPrefix`, `CookieName` |
| 变量 | camelCase | `articleID`, `deviceID`, `uploadCfg` |
| JSON Tag | camelCase（前端约定） | `json:"coverImage"`, `json:"likeCount"` |
| GORM Column Tag | snake_case | `gorm:"column:cover_image"`, `gorm:"column:like_count"` |

## 包导入顺序

标准库 → 第三方库 → 项目内部包，组间空行分隔：

```go
import (
    "fmt"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"

    "github.com/fatWill/agent-team-blog/backend/models"
    "github.com/fatWill/agent-team-blog/backend/utils"
)
```

## Handler 函数规范

### 函数签名

所有 handler 函数接收 `*gin.Context` 参数，无返回值：

```go
// GetArticles GET /api/articles
func GetArticles(c *gin.Context) {
    // ...
}
```

### 函数注释

导出函数必须有注释，格式为：`// 函数名 HTTP方法 路径`

```go
// CreateArticle POST /api/articles
func CreateArticle(c *gin.Context) { ... }

// DeleteMessage DELETE /api/messages/:id（需鉴权）
func DeleteMessage(c *gin.Context) { ... }
```

### 参数绑定

- JSON Body：使用匿名结构体 + `c.ShouldBindJSON`
- URL 参数：使用 `c.Param("id")` + `strconv.ParseUint`
- Query 参数：使用 `c.Query("key")`

```go
// Body 参数
var body struct {
    Title   string `json:"title"`
    Content string `json:"content"`
}
if err := c.ShouldBindJSON(&body); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "参数解析失败"})
    return
}

// Path 参数
id, err := strconv.ParseUint(c.Param("id"), 10, 64)
if err != nil || id == 0 {
    c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的 ID"})
    return
}
```

### 输入校验

参数绑定后立即校验，校验失败立即返回：

```go
if strings.TrimSpace(body.Title) == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "标题不能为空"})
    return
}
```

## 错误处理规范

### 统一错误响应格式

```go
c.JSON(http.StatusXxx, gin.H{
    "error":         true,
    "statusCode":    4xx/5xx,
    "statusMessage": "人类可读的错误描述",
})
```

### 常见错误码

| HTTP 状态码 | 场景 |
|-------------|------|
| 400 | 参数缺失、格式错误、校验失败 |
| 401 | 未登录、Token 过期 |
| 403 | 无权限、操作被限制（如每日修改限制） |
| 404 | 资源不存在 |
| 429 | IP 限频触发 |
| 500 | 数据库操作失败、文件操作失败 |

### 成功响应格式

根据业务场景使用不同格式：

```go
// 列表
c.JSON(http.StatusOK, gin.H{"list": list})

// 单条记录
c.JSON(http.StatusOK, gin.H{"id": id, "title": title, ...})

// 操作成功
c.JSON(http.StatusOK, gin.H{"ok": true})
c.JSON(http.StatusOK, gin.H{"success": true})
```

## 数据库操作规范

### 全局 DB 实例

通过 `utils.DB`（`*gorm.DB`）进行所有数据库操作：

```go
// 查询
utils.DB.Where("id = ?", id).First(&article)

// 创建
utils.DB.Create(&article)

// 更新（Map 方式，支持零值更新）
utils.DB.Model(&models.Article{}).Where("id = ?", id).Updates(map[string]interface{}{
    "title": title,
    "content": content,
})

// 删除
utils.DB.Where("id = ?", id).Delete(&models.Article{})
```

### 参数化查询

**禁止**拼接 SQL 字符串，必须使用参数化查询：

```go
// ✅ 正确
utils.DB.Where("title LIKE ?", "%"+title+"%")

// ❌ 错误（SQL 注入风险）
utils.DB.Where("title LIKE '%" + title + "%'")
```

### 原生 SQL

仅在 GORM 链式调用无法满足时使用 `Raw`：

```go
utils.DB.Raw(`
    SELECT a.id, a.name, IFNULL(pc.cnt, 0) AS photo_count
    FROM albums a
    LEFT JOIN (...) pc ON pc.album_id = a.id
    ORDER BY a.created_at DESC
`).Scan(&rows)
```

## Redis 操作规范

### 全局 Redis 实例

通过 `utils.RDB`（`*redis.Client`）进行操作：

```go
ctx := context.Background()

// 读取
val, err := utils.RDB.Get(ctx, key).Result()

// 写入（带 TTL）
utils.RDB.Set(ctx, key, value, 30*24*time.Hour)

// 续期
utils.RDB.Expire(ctx, key, ttl)

// 删除
utils.RDB.Del(ctx, key)
```

### Key 命名

| 前缀 | 格式 | 说明 |
|------|------|------|
| `auth_token:` | `auth_token:{64位hex}` | 登录 Token |
| `theme:` | `theme:{uid}` | 主题偏好 |

## 安全规范

| 风险 | 防护措施 | 实现位置 |
|------|----------|----------|
| SQL 注入 | GORM 参数化查询 | handlers/ |
| 密码泄露 | bcrypt 哈希存储，JSON tag `-` 隐藏 | models/、handlers/ |
| 路径穿越 | uploadId 正则校验 `^[\w-]+$` | handlers/upload.go |
| CSRF | Cookie httpOnly + SameSite | middleware/auth.go |
| 暴力破解 | IP 限频中间件 | middleware/ratelimit.go |
| 文件类型 | MIME + 扩展名双重校验 | handlers/upload.go |

## 代码质量红线

- **禁止**在 handler 中硬编码数据库连接信息
- **禁止**在日志中输出密码、Token 等敏感信息
- **禁止**使用 `interface{}` 做参数透传（当前项目无此问题）
- **必须**所有导出函数有注释
- **必须**error 被处理或显式忽略（`_ = fn()`）
- **必须**`strings.TrimSpace` 处理用户输入
