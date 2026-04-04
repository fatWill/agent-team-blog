# 留言板接口文档

## [GET] /api/messages - 留言列表

### 描述
获取所有留言，按创建时间倒序。根据 `deviceId` 参数判断每条留言是否属于当前设备，以及是否可编辑。

### 请求参数 (Query)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| deviceId | string | ❌ | 设备指纹（用于判断 isOwn 和 canEdit） |

### 成功响应 (200)

```json
{
  "list": [
    {
      "id": 1,
      "nickname": "路人甲",
      "content": "这个博客不错！",
      "isOwn": false,
      "canEdit": false,
      "createdAt": "2026-01-01T12:00:00+08:00",
      "updatedAt": "2026-01-01T12:00:00+08:00"
    }
  ]
}
```

### 字段说明

| 字段 | 说明 |
|------|------|
| `isOwn` | 当前 deviceId 与留言 device_id 匹配时为 true |
| `canEdit` | isOwn=true 且今天未修改过时为 true |
| `nickname` | 为 null 或空时显示 "匿名" |

---

## [POST] /api/messages - 创建/修改留言

### 描述
每个设备只能有一条留言。如果该 deviceId 已有留言，则视为修改操作（受每日一次限制）。

### 限频
IP 限频：每分钟 3 次（`middleware.RateLimit(3)`）

### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| deviceId | string | ✅ | 设备指纹（8-64 字符） |
| nickname | string | ❌ | 昵称 |
| content | string | ✅ | 留言内容 |

### 请求示例

```json
{
  "deviceId": "a1b2c3d4e5f6g7h8",
  "nickname": "路人甲",
  "content": "这个博客不错！"
}
```

### 成功响应 (200)

```json
{
  "id": 1,
  "nickname": "路人甲",
  "content": "这个博客不错！",
  "isOwn": true,
  "canEdit": true,
  "createdAt": "2026-01-01T12:00:00+08:00",
  "updatedAt": "2026-01-01T12:00:00+08:00"
}
```

> 修改操作成功后 `canEdit` 为 false（当天已修改）。

### 错误响应

| 状态码 | statusMessage | 场景 |
|--------|---------------|------|
| 400 | 无效的设备 ID | deviceId 长度不在 8-64 |
| 400 | 留言内容不能为空 | content 为空 |
| 403 | 今天已经修改过了，明天再来吧 | 当天已修改（UTC+8） |
| 429 | 操作太频繁，请稍后再试 | IP 限频触发 |

### 业务逻辑

```
1. 校验 deviceId（8-64字符）和 content（非空）
2. 获取客户端 IP（X-Forwarded-For → X-Real-IP → ClientIP）
3. 查询 device_id 是否已有留言
   ├── 有 → 检查 last_modified_date 是否为今天
   │     ├── 是今天 → 403 拒绝
   │     └── 不是 → 更新 nickname/content/last_modified_date/ip
   └── 无 → 创建新留言
```

---

## [PUT] /api/messages/:id - 修改留言

### 描述
修改指定 ID 的留言，需验证设备归属和每日修改限制。

### 请求参数

**Path**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | uint64 | ✅ | 留言 ID |

**Body - JSON**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| deviceId | string | ✅ | 设备指纹（用于验证归属） |
| nickname | string | ❌ | 昵称 |
| content | string | ✅ | 留言内容 |

### 成功响应 (200)

```json
{
  "id": 1,
  "nickname": "路人甲",
  "content": "修改后的内容",
  "isOwn": true,
  "canEdit": false,
  "createdAt": "2026-01-01T12:00:00+08:00",
  "updatedAt": "2026-01-02T10:00:00+08:00"
}
```

### 错误响应

| 状态码 | statusMessage | 场景 |
|--------|---------------|------|
| 400 | 无效的留言 ID | id 非正整数 |
| 400 | 无效的设备 ID | deviceId 长度不合法 |
| 400 | 留言内容不能为空 | content 为空 |
| 403 | 无权修改此留言 | deviceId 与留言 device_id 不匹配 |
| 403 | 今天已经修改过了，明天再来吧 | 当天已修改 |
| 404 | 留言不存在 | 数据库无此记录 |

---

## [DELETE] /api/messages/:id - 删除留言 🔒

### 描述
管理员删除指定留言。需要鉴权。

### 鉴权
Cookie `auth_token`（httpOnly）

### 请求参数 (Path)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | uint64 | ✅ | 留言 ID |

### 成功响应 (200)

```json
{ "success": true }
```

### 错误响应

| 状态码 | 响应 | 场景 |
|--------|------|------|
| 400 | `{"error": "无效的留言 ID"}` | id 非正整数 |
| 401 | 标准 401 响应 | 未登录 |
| 404 | `{"error": "留言不存在"}` | 数据库无此记录 |
| 500 | `{"error": "删除留言失败"}` | 数据库异常 |

> ⚠️ 此接口错误响应格式与其他接口不一致（使用 `{"error": "msg"}` 而非标准格式），待后续统一。
