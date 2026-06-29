# 微信公众号同步 API

## 概述
博客文章发布后自动同步到微信公众号草稿箱，支持手动触发同步、查看同步日志、查看 access_token 状态和服务器 IP。

---

## [POST] /api/admin/articles/:id/sync-wechat - 手动触发微信同步

### 描述
手动触发指定文章同步到微信公众号草稿箱。如果文章已有草稿，则更新；否则创建新草稿。

### 请求头
| Header | 必填 | 说明 |
|--------|------|------|
| Cookie: auth_token=xxx | ✅ | 管理员 Token |

### 路径参数
| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 文章 UUID |

### 成功响应 (200)
```json
{
  "ok": true,
  "message": "同步成功"
}
```

### 错误响应
| StatusCode | Message | 场景 |
|------------|---------|------|
| 400 | 缺少文章 ID | 路径参数为空 |
| 401 | 未登录 | Token 无效 |
| 404 | 文章不存在 | 文章 ID 不存在 |
| 500 | 同步失败: xxx | 微信 API 调用失败 |

---

## [GET] /api/admin/articles/:id/wechat-sync-logs - 查看同步日志

### 描述
查看指定文章的微信同步操作日志（最近 50 条）。

### 请求头
| Header | 必填 | 说明 |
|--------|------|------|
| Cookie: auth_token=xxx | ✅ | 管理员 Token |

### 路径参数
| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 文章 UUID |

### 成功响应 (200)
```json
{
  "list": [
    {
      "id": 1,
      "articleId": "abc-123",
      "action": "create_draft",
      "status": "success",
      "error": null,
      "createdAt": "2026-06-29T10:00:00Z"
    },
    {
      "id": 2,
      "articleId": "abc-123",
      "action": "upload_image",
      "status": "failed",
      "error": "图片超过 2MB 限制",
      "createdAt": "2026-06-29T10:00:01Z"
    }
  ]
}
```

### action 枚举
| 值 | 说明 |
|------|------|
| create_draft | 创建草稿 |
| update_draft | 更新草稿 |
| upload_image | 上传图片 |

---

## [GET] /api/admin/wechat/access-token-status - 查看 access_token 状态

### 描述
查看微信 access_token 缓存状态，便于排查连接问题。

### 请求头
| Header | 必填 | 说明 |
|--------|------|------|
| Cookie: auth_token=xxx | ✅ | 管理员 Token |

### 成功响应 (200)
```json
{
  "hasToken": true,
  "remainingSec": 5400,
  "lastRefresh": "2026-06-29T09:30:00Z"
}
```

---

## [GET] /api/admin/wechat/server-ip - 查看服务器出口 IP

### 描述
返回 VPS 出口公网 IP，方便用户在微信公众号后台添加 IP 白名单。

### 请求头
| Header | 必填 | 说明 |
|--------|------|------|
| Cookie: auth_token=xxx | ✅ | 管理员 Token |

### 成功响应 (200)
```json
{
  "outboundIP": "203.195.213.129",
  "wxAPIIPs": ["101.226.103.0/25", "..."],
  "hint": "请将 outboundIP 添加到微信公众号后台 → 基本配置 → IP 白名单",
  "mpConsoleURL": "https://mp.weixin.qq.com/"
}
```

---

## 自动同步行为

### 触发条件
1. `POST /api/articles` 创建文章成功后，异步触发同步
2. `PUT /api/articles/:id` 更新文章成功后，异步触发同步

### 同步流程
1. 将 Tiptap JSON 内容转换为公众号兼容 HTML
2. 下载文章中的图片并上传到微信（使用 uploadimg 接口，返回永久 URL）
3. 下载封面图并上传为永久素材（使用 add_material 接口，获取 thumb_media_id）
4. 调用草稿接口创建/更新草稿
5. 更新文章的 `wechat_sync_status` 和 `wechat_draft_media_id`

### 文章新增字段
| 字段 | 类型 | 说明 |
|------|------|------|
| wechat_draft_media_id | string | 微信草稿 media_id |
| wechat_synced_at | datetime | 最后同步时间 |
| wechat_sync_status | string | pending/syncing/success/failed |
| wechat_sync_error | string | 失败原因 |
| wechat_auto_sync | int | 是否启用自动同步（1=是 0=否） |
