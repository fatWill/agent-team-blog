# 文章接口文档

## [GET] /api/articles - 文章列表

### 描述
获取所有文章列表（不含 content 字段），支持按标题模糊搜索。

### 请求参数 (Query)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | ❌ | 标题关键词（模糊搜索） |

### 成功响应 (200)

```json
{
  "list": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "文章标题",
      "summary": "文章摘要",
      "coverImage": "/uploads/cover.jpg",
      "likeCount": 5,
      "createdAt": "2026-01-01T12:00:00+08:00",
      "updatedAt": "2026-01-01T12:00:00+08:00"
    }
  ]
}
```

---

## [GET] /api/articles/:id - 文章详情

### 描述
获取单篇文章完整内容（含 Tiptap JSON content）。

### 请求参数 (Path)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | ✅ | 文章 UUID |

### 成功响应 (200)

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "文章标题",
  "summary": "文章摘要",
  "coverImage": "/uploads/cover.jpg",
  "content": { "type": "doc", "content": [...] },
  "likeCount": 5,
  "createdAt": "2026-01-01T12:00:00+08:00",
  "updatedAt": "2026-01-01T12:00:00+08:00"
}
```

### 错误响应

| 状态码 | statusMessage | 场景 |
|--------|---------------|------|
| 400 | 缺少文章 ID | id 参数为空 |
| 404 | 文章不存在 | 数据库无此记录 |

---

## [POST] /api/articles - 创建文章 🔒

### 描述
创建一篇新文章，ID 由服务端生成（UUID v4）。

### 鉴权
Cookie `auth_token`（httpOnly）

### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | ✅ | 文章标题（不能为空） |
| summary | string | ❌ | 文章摘要 |
| coverImage | string | ❌ | 封面图 URL |
| content | object | ✅ | Tiptap JSON 对象（不能为 null） |

### 请求示例

```json
{
  "title": "Go 并发编程指南",
  "summary": "一篇关于 Go 并发编程的实践指南",
  "coverImage": "/uploads/cover.jpg",
  "content": { "type": "doc", "content": [{ "type": "paragraph", "content": [{ "type": "text", "text": "正文内容" }] }] }
}
```

### 成功响应 (200)

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Go 并发编程指南",
  "summary": "一篇关于 Go 并发编程的实践指南",
  "coverImage": "/uploads/cover.jpg",
  "content": { ... },
  "likeCount": 0,
  "createdAt": "2026-01-01T12:00:00+08:00",
  "updatedAt": "2026-01-01T12:00:00+08:00"
}
```

### 错误响应

| 状态码 | statusMessage | 场景 |
|--------|---------------|------|
| 400 | 参数解析失败 | JSON 格式错误 |
| 400 | 文章标题不能为空 | title 为空字符串 |
| 400 | 文章内容不能为空，且必须为 Tiptap JSON 对象 | content 缺失或为 null |
| 401 | 未登录，请先登录 | 未携带 Cookie |
| 500 | 创建文章失败 | 数据库异常 |

---

## [PUT] /api/articles/:id - 更新文章 🔒

### 描述
部分更新文章，至少提供一个字段。

### 鉴权
Cookie `auth_token`（httpOnly）

### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | ❌ | 文章标题 |
| summary | string | ❌ | 文章摘要 |
| coverImage | string | ❌ | 封面图 URL |
| content | object | ❌ | Tiptap JSON 对象 |

> 至少需要提供以上一个字段。

### 成功响应 (200)

返回更新后的完整文章对象（同创建响应）。

### 错误响应

| 状态码 | statusMessage | 场景 |
|--------|---------------|------|
| 400 | 至少需要提供 title、summary、coverImage 或 content 中的一个字段 | 无有效字段 |
| 404 | 文章不存在 | RowsAffected=0 |

---

## [DELETE] /api/articles/:id - 删除文章 🔒

### 描述
永久删除一篇文章。

### 鉴权
Cookie `auth_token`（httpOnly）

### 成功响应 (200)

```json
{ "ok": true }
```

### 错误响应

| 状态码 | statusMessage | 场景 |
|--------|---------------|------|
| 400 | 缺少文章 ID | id 为空 |
| 404 | 文章不存在 | 数据库无此记录 |

---

## [POST] /api/articles/:id/like - 点赞/取消

### 描述
切换文章点赞状态。已赞则取消，未赞则点赞。同步更新 `articles.like_count` 冗余字段。

### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| deviceId | string | ✅ | 设备指纹（8-64 字符） |

### 成功响应 (200)

```json
{
  "liked": true,
  "likeCount": 6
}
```

### 错误响应

| 状态码 | statusMessage | 场景 |
|--------|---------------|------|
| 400 | 无效的设备 ID | deviceId 长度不在 8-64 范围 |
| 404 | 文章不存在 | 文章 ID 无效 |

---

## [GET] /api/articles/:id/like-status - 单篇点赞状态

### 描述
查询当前设备对某篇文章的点赞状态。

### 请求参数 (Query)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| deviceId | string | ❌ | 设备指纹（不传则 liked=false） |

### 成功响应 (200)

```json
{
  "liked": true,
  "likeCount": 6
}
```

---

## [GET] /api/articles/like-status-batch - 批量点赞状态

### 描述
批量查询当前设备对多篇文章的点赞状态（最多 100 篇）。

### 请求参数 (Query)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| deviceId | string | ✅ | 设备指纹 |
| ids | string | ✅ | 文章 ID 列表（逗号分隔，最多 100 个） |

### 成功响应 (200)

```json
{
  "likedIds": ["uuid-1", "uuid-2"]
}
```
