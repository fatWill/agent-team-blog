# HTTP API 设计规范

## 基础信息

| 项目 | 值 |
|------|------|
| Base URL | `https://fatwill.cloud/api` |
| 协议 | HTTPS（Nginx SSL 终端） |
| 内容类型 | `application/json`（除文件上传外） |
| 字符编码 | UTF-8 |
| 时区 | UTC+8（CST） |

## RESTful 约定

| HTTP 方法 | 语义 | 示例 |
|-----------|------|------|
| GET | 查询资源 | `GET /api/articles` |
| POST | 创建资源 / 执行操作 | `POST /api/articles` |
| PUT | 更新资源（部分更新） | `PUT /api/articles/:id` |
| DELETE | 删除资源 | `DELETE /api/articles/:id` |

### 路径命名规则

- 资源名使用**复数**：`/articles`, `/albums`, `/photos`, `/messages`
- 子资源使用嵌套路径：`/albums/:id/photos`
- 操作使用动词：`/articles/:id/like`, `/albums/:id/verify-password`

## 鉴权方式

### Cookie Token 认证

- **Cookie 名称**：`auth_token`
- **属性**：`httpOnly=true`, `secure=false`, `path=/`, `sameSite=lax`
- **有效期**：30 天（2592000 秒）
- **滚动续期**：每次通过鉴权中间件后自动刷新

### 鉴权流程

```
1. POST /api/auth/login → 服务端生成 64 位 hex Token
2. Token 存入 Redis（auth_token:{token} → username, TTL=30天）
3. 写入 httpOnly Cookie（auth_token={token}）
4. 后续请求自动携带 Cookie
5. AuthRequired 中间件验证 Token → 滚动续期 Redis TTL + Cookie MaxAge
```

### 未鉴权响应

```json
{
  "error": true,
  "statusCode": 401,
  "statusMessage": "未登录，请先登录"
}
```

## 统一响应格式

### 成功响应

项目**未使用**统一的 `{code, message, data}` 包装，而是根据场景直接返回：

#### 列表类

```json
{
  "list": [...]
}
```

#### 单条记录

```json
{
  "id": "xxx",
  "title": "xxx",
  "createdAt": "2026-01-01T00:00:00+08:00",
  ...
}
```

#### 操作类

```json
{"ok": true}
// 或
{"success": true}
```

#### 点赞/踩类

```json
{
  "liked": true,
  "likeCount": 42
}
// 或
{
  "count": 10,
  "liked": true
}
```

### 错误响应

统一格式：

```json
{
  "error": true,
  "statusCode": 400,
  "statusMessage": "人类可读的错误描述"
}
```

> 注意：`DeleteMessage` 接口使用了简化格式 `{"error": "xxx"}`，与其他接口不完全一致。

## 错误码说明

| HTTP 状态码 | 场景 | statusMessage 示例 |
|-------------|------|-------------------|
| 400 | 参数缺失/格式错误 | "参数解析失败"、"文章标题不能为空" |
| 401 | 未登录/Token 过期 | "未登录，请先登录"、"Token 已过期或无效，请重新登录" |
| 403 | 无权限/操作受限 | "无权修改此留言"、"今天已经修改过了，明天再来吧" |
| 404 | 资源不存在 | "文章不存在"、"相册不存在" |
| 429 | 请求频率过高 | "操作太频繁，请稍后再试" |
| 500 | 服务端异常 | "创建文章失败"、"密码加密失败" |

## 限频策略

| 接口 | 限制 | 实现方式 |
|------|------|----------|
| `POST /api/messages` | 每 IP 每分钟 3 次 | 内存 Map（middleware/ratelimit.go） |

限频触发时返回 HTTP 429：

```json
{
  "error": true,
  "statusCode": 429,
  "statusMessage": "操作太频繁，请稍后再试"
}
```

## CORS 配置

```go
cors.Config{
    AllowOrigins:     []string{"https://fatwill.cloud"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Cookie"},
    ExposeHeaders:    []string{"Set-Cookie"},
    AllowCredentials: true,
}
```

## 文件上传

### 直传（POST /api/upload）

- Content-Type: `multipart/form-data`
- 字段名: `file`
- 支持格式: jpg, jpeg, png, gif, webp
- 双重校验: MIME 类型 + 文件扩展名
- 返回: `{"url": "https://fatwill-cloud-1253664788.cos.ap-guangzhou.myqcloud.com/upload/xxx.jpg"}`
- 存储后端: 腾讯云 COS 对象存储
- 上传策略: ≤2MB PutObject 直传，>2MB COS 原生分片上传

### 分片上传

1. `POST /api/upload/chunk` — 逐片上传（FormData: file + uploadId + chunkIndex + totalChunks + filename）
2. `POST /api/upload/merge` — 合并所有分片（JSON: uploadId + totalChunks + filename）
3. `DELETE /api/upload/chunk` — 取消上传，清理临时文件（JSON: uploadId）

安全措施：`uploadId` 使用正则 `^[\w-]+$` 校验，防止路径穿越攻击。

## 设备标识（deviceId）

多个接口使用 `deviceId` 区分匿名用户：

- 长度限制：8-64 字符
- 用途：文章点赞、照片点赞/踩、留言板归属判断
- 传递方式：Query 参数或 JSON Body

## 时间格式

| 场景 | 格式 | 示例 |
|------|------|------|
| 文章时间 | RFC3339 | `2026-01-01T12:00:00+08:00` |
| 相册/照片时间 | ISO 8601 | `2026-01-01T12:00:00.000Z` |
| 留言时间 | RFC3339 | `2026-01-01T12:00:00+08:00` |
| 留言修改日期 | 日期字符串 | `2026-01-01` |

> 注意：时间格式在不同模块间存在不一致（RFC3339 vs ISO 8601），待后续统一。
