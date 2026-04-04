# 认证接口文档

## [POST] /api/auth/login - 登录

### 描述
使用账号密码登录，成功后服务端生成 64 位 hex Token，存入 Redis 并写入 httpOnly Cookie。

### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | ✅ | 用户名 |
| password | string | ✅ | 密码 |

### 请求示例

```json
{
  "username": "fafa",
  "password": "xxxxxx"
}
```

### 成功响应 (200)

```json
{
  "token": "a1b2c3d4e5f6...（64位hex字符串）"
}
```

**Set-Cookie**：`auth_token=a1b2c3d4...; Path=/; HttpOnly; Max-Age=2592000`

### 错误响应

| 状态码 | statusMessage | 场景 |
|--------|---------------|------|
| 400 | 请提供用户名和密码 | 参数缺失 |
| 401 | 用户名或密码错误 | 认证失败 |
| 500 | 生成 Token 失败 | 随机数生成异常 |
| 500 | 保存 Token 失败 | Redis 写入异常 |

### Token 存储

- **Redis Key**：`auth_token:{token}`
- **Redis Value**：username
- **TTL**：30 天（`30 * 24 * time.Hour`）

---

## [GET] /api/auth/check - 验证登录态 🔒

### 描述
验证当前 Cookie 中的 Token 是否有效。经过 `AuthRequired` 中间件，Token 有效则自动滚动续期。

### 鉴权
Cookie `auth_token`（httpOnly）

### 成功响应 (200)

```json
{
  "ok": true,
  "username": "fafa"
}
```

### 错误响应

| 状态码 | statusMessage | 场景 |
|--------|---------------|------|
| 401 | 未登录，请先登录 | 无 Cookie |
| 401 | Token 已过期或无效，请重新登录 | Redis 中 Token 不存在 |

---

## [POST] /api/auth/logout - 退出登录

### 描述
退出登录。删除 Redis 中的 Token，清除客户端 Cookie。

### 成功响应 (200)

```json
{
  "ok": true
}
```

**Set-Cookie**：`auth_token=; Path=/; HttpOnly; Max-Age=-1`（清除 Cookie）

### 说明

- 即使未携带有效 Cookie，接口也会返回 200
- Redis Token 删除使用 `DEL` 命令

---

## 鉴权中间件详解

### AuthRequired

```
请求 → 读取 Cookie(auth_token) → Redis GET(auth_token:{token})
  │
  ├── Cookie 不存在 → 401 "未登录，请先登录"
  ├── Redis 查无此 Token → 401 "Token 已过期或无效"
  └── Token 有效 →
        ├── Redis EXPIRE（续期 30 天）
        ├── Set-Cookie（刷新 MaxAge）
        ├── c.Set("username", username)
        └── c.Next()
```

### 常量定义

```go
const (
    TokenTTL       = 30 * 24 * time.Hour  // 30 天
    TokenKeyPrefix = "auth_token:"
    CookieName     = "auth_token"
)
```
