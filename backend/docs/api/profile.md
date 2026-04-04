# 个人资料接口文档

## [GET] /api/profile - 获取博主资料

### 描述
获取博主个人资料（头像和简介）。profile 表为单行表（id=1），如果记录不存在则返回空值。

### 成功响应 (200)

```json
{
  "avatar": "/uploads/avatar.jpg",
  "bio": "一个热爱编程的开发者"
}
```

> 记录不存在时返回 `{"avatar": "", "bio": ""}`，不会返回错误。

---

## [PUT] /api/profile - 更新博主资料 🔒

### 描述
更新博主个人资料，支持部分更新。操作 `profile` 表中 id=1 的记录。

### 鉴权
Cookie `auth_token`（httpOnly）

### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| avatar | string | ❌ | 头像 URL |
| bio | string | ❌ | 个人简介 |

> 至少需要提供一个字段。

### 请求示例

```json
{
  "avatar": "/uploads/new-avatar.jpg",
  "bio": "更新后的个人简介"
}
```

### 成功响应 (200)

```json
{
  "avatar": "/uploads/new-avatar.jpg",
  "bio": "更新后的个人简介"
}
```

### 错误响应

| 状态码 | statusMessage | 场景 |
|--------|---------------|------|
| 400 | 参数解析失败 | JSON 格式错误 |
| 400 | 至少需要提供 avatar 或 bio 中的一个字段 | 无有效字段 |
| 401 | 未登录，请先登录 | 未携带 Cookie |
