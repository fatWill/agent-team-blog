# 留言板接口文档

## [GET] /api/messages - 留言列表

### 描述
获取所有留言，按创建时间倒序。根据 `deviceId` 参数判断每条留言是否属于当前设备。昵称根据 deviceId 自动生成（`用户` + deviceId 后4位）。

### 请求参数 (Query)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| deviceId | string | ❌ | 设备指纹（用于判断 isOwn） |

### 成功响应 (200)

```json
{
  "list": [
    {
      "id": 1,
      "nickname": "用户g7h8",
      "content": "这个博客不错！",
      "isOwn": false,
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
| `nickname` | 自动生成：`用户` + deviceId 后4位 |

---

## [POST] /api/messages - 创建留言

### 描述
每个设备只能发一条留言。昵称自动生成（`用户` + deviceId 后4位），不接收 nickname 参数。

### 限频
IP 限频：每分钟 3 次（`middleware.RateLimit(3)`）

### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| deviceId | string | ✅ | 设备指纹（8-64 字符） |
| content | string | ✅ | 留言内容 |

### 请求示例

```json
{
  "deviceId": "a1b2c3d4e5f6g7h8",
  "content": "这个博客不错！"
}
```

### 成功响应 (200)

```json
{
  "id": 1,
  "nickname": "用户g7h8",
  "content": "这个博客不错！",
  "isOwn": true,
  "createdAt": "2026-01-01T12:00:00+08:00",
  "updatedAt": "2026-01-01T12:00:00+08:00"
}
```

### 错误响应

| 状态码 | statusMessage | 场景 |
|--------|---------------|------|
| 400 | 无效的设备 ID | deviceId 长度不在 8-64 |
| 400 | 留言内容不能为空 | content 为空 |
| 400 | 每个设备只能留言一次 | 该 deviceId 已有留言 |
| 429 | 操作太频繁，请稍后再试 | IP 限频触发 |

### 业务逻辑

```
1. 校验 deviceId（8-64字符）和 content（非空）
2. 查询 device_id 是否已有留言 → 有则返回 400
3. 自动生成昵称：用户 + deviceId 后4位
4. 获取客户端 IP
5. 创建新留言（唯一索引兜底并发场景）
```

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
