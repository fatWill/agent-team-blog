# PV/UV 统计接口文档

## [POST] /api/pv/record - 上报访问记录

### 描述
前端页面加载时调用，记录一次页面访问。服务端自动从 User-Agent 解析设备类型/浏览器/操作系统，从请求头获取 IP。

### 鉴权
❌ 无需鉴权

### 防刷策略
同一 `device_id` + `path` 组合，**60 秒内只记录一次**。使用 Redis key `pv:{device_id}:{path}`（TTL 60s）。Redis 不可用时降级为直接写入。

### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| path | string | ✅ | 访问路径，如 `/home`、`/articles/123` |
| device_id | string | ✅ | 设备唯一 ID（前端生成的 UUID） |
| referer | string | ❌ | 来源页面 |

### 请求示例

```json
{
  "path": "/home",
  "device_id": "uuid-xxx-xxx",
  "referer": "https://google.com"
}
```

### 成功响应 (200)

```json
{ "ok": true }
```

### 错误响应

| 状态码 | 说明 |
|--------|------|
| 400 | `path` 或 `device_id` 为空 |
| 500 | 数据库写入失败 |

### 服务端自动解析字段

| 字段 | 来源 | 说明 |
|------|------|------|
| ip | 请求头 `X-Forwarded-For` > `X-Real-IP` > `RemoteAddr` | 客户端 IP |
| user_agent | 请求头 `User-Agent` | 原始 UA 字符串 |
| device_type | UA 解析 | `mobile` / `tablet` / `desktop` |
| browser | UA 解析 | `Chrome` / `Safari` / `Firefox` / `Edge` / `Opera` / `Other` |
| os | UA 解析 | `Windows` / `macOS` / `iOS` / `iPadOS` / `Android` / `Linux` / `Other` |

---

## [GET] /api/pv/trend - PV/UV 趋势 🔒

### 描述
返回最近 N 天每天的 PV（总访问数）和 UV（独立设备数），用于趋势图展示。

### 鉴权
✅ Cookie `auth_token`

### 请求参数 (Query)

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| days | int | ❌ | 7 | 天数，仅支持 `7` 或 `30` |

### 成功响应 (200)

```json
{
  "ok": true,
  "data": [
    { "date": "2026-04-01", "pv": 123, "uv": 45 },
    { "date": "2026-04-02", "pv": 89, "uv": 32 },
    { "date": "2026-04-03", "pv": 0, "uv": 0 }
  ]
}
```

### 字段说明

| 字段 | 说明 |
|------|------|
| `date` | 日期（UTC+8），格式 `YYYY-MM-DD` |
| `pv` | 当天总访问次数 |
| `uv` | 当天独立设备数（按 `device_id` 去重） |

> 没有数据的日期会自动补全为 `pv=0, uv=0`。

---

## [GET] /api/pv/top-pages - Top5 页面 🔒

### 描述
返回最近 N 天访问量最高的 5 个页面。

### 鉴权
✅ Cookie `auth_token`

### 请求参数 (Query)

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| days | int | ❌ | 7 | 天数，仅支持 `7` 或 `30` |

### 成功响应 (200)

```json
{
  "ok": true,
  "data": [
    { "path": "/home", "pv": 500, "uv": 200 },
    { "path": "/articles/abc-123", "pv": 120, "uv": 80 },
    { "path": "/albums", "pv": 90, "uv": 60 },
    { "path": "/changelog", "pv": 50, "uv": 30 },
    { "path": "/articles/def-456", "pv": 40, "uv": 25 }
  ]
}
```

---

## [GET] /api/pv/logs - 访问日志列表 🔒

### 描述
分页查询访问日志，支持按路径和日期筛选。

### 鉴权
✅ Cookie `auth_token`

### 请求参数 (Query)

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| page | int | ❌ | 1 | 页码 |
| page_size | int | ❌ | 20 | 每页条数（最大 100） |
| path | string | ❌ | - | 按路径精确筛选 |
| date | string | ❌ | - | 按日期筛选，格式 `2026-04-01` |

### 成功响应 (200)

```json
{
  "ok": true,
  "data": {
    "list": [
      {
        "id": 1,
        "path": "/home",
        "device_id": "uuid-xxx",
        "ip": "1.2.3.4",
        "device_type": "mobile",
        "browser": "Chrome",
        "os": "iOS",
        "referer": "",
        "created_at": "2026-04-06 10:00:00"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 字段说明

| 字段 | 说明 |
|------|------|
| `created_at` | 格式 `YYYY-MM-DD HH:mm:ss`（UTC+8） |

---

## [GET] /api/pv/overview - 统计概览 🔒

### 描述
返回今日和总计的 PV/UV 数据。

### 鉴权
✅ Cookie `auth_token`

### 成功响应 (200)

```json
{
  "ok": true,
  "data": {
    "today_pv": 50,
    "today_uv": 30,
    "total_pv": 10000,
    "total_uv": 5000
  }
}
```

### 字段说明

| 字段 | 说明 |
|------|------|
| `today_pv` | 今日（UTC+8）总访问次数 |
| `today_uv` | 今日独立设备数 |
| `total_pv` | 历史总访问次数 |
| `total_uv` | 历史总独立设备数 |
