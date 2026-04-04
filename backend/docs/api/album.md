# 相册 & 照片接口文档

## 相册接口

### [GET] /api/albums - 相册列表

#### 描述
获取所有相册列表，包含照片数量和是否加密信息。使用 LEFT JOIN 子查询统计 photo_count。

#### 成功响应 (200)

```json
{
  "list": [
    {
      "id": 1,
      "name": "旅行照片",
      "description": "2026年春节旅行",
      "coverUrl": "https://assets.fatwill.cloud/upload/20260101/cover.jpg",
      "photoCount": 12,
      "hasPassword": false,
      "createdAt": "2026-01-01T12:00:00.000Z",
      "updatedAt": "2026-01-01T12:00:00.000Z"
    }
  ]
}
```

---

### [POST] /api/albums - 创建相册 🔒

#### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | ✅ | 相册名称 |
| description | string | ❌ | 相册描述 |
| password | string | ❌ | 访问密码（bcrypt 加密存储） |

#### 成功响应 (200)

返回 `AlbumListItem` 对象（photoCount=0）。

#### 错误响应

| 状态码 | statusMessage | 场景 |
|--------|---------------|------|
| 400 | 相册名称不能为空 | name 为空 |
| 500 | 密码加密失败 | bcrypt 异常 |
| 500 | 创建相册失败 | 数据库异常 |

---

### [PUT] /api/albums/:id - 更新相册 🔒

#### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | ❌ | 相册名称 |
| description | string | ❌ | 相册描述 |
| password | string | ❌ | 密码（空字符串=清除密码） |

> 至少提供一个字段。password 传空字符串 `""` 表示清除密码。

#### 成功响应 (200)

返回更新后的 `AlbumListItem` 对象。

---

### [DELETE] /api/albums/:id - 删除相册 🔒

#### 描述
删除相册并**级联删除**所有关联照片。

#### 成功响应 (200)

```json
{ "ok": true }
```

---

### [POST] /api/albums/:id/verify-password - 验证相册密码

#### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| password | string | ✅ | 待验证的密码 |

#### 成功响应 (200)

```json
{ "success": true }
```

> `success` 为 false 表示密码错误。

---

## 照片接口

### [GET] /api/albums/:id/photos - 照片列表

#### 描述
获取指定相册的所有照片，按创建时间倒序。

#### 成功响应 (200)

```json
{
  "list": [
    {
      "id": 1,
      "albumId": 1,
      "url": "https://assets.fatwill.cloud/upload/20260101/photo.jpg",
      "caption": "夕阳下",
      "hasPassword": false,
      "createdAt": "2026-01-01T12:00:00.000Z",
      "updatedAt": "2026-01-01T12:00:00.000Z"
    }
  ]
}
```

---

### [POST] /api/albums/:id/photos - 添加照片 🔒

#### 描述
向指定相册添加一张照片。添加后自动更新相册封面为最新照片。

#### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| url | string | ✅ | 图片 URL |
| caption | string | ❌ | 图片说明 |
| password | string | ❌ | 访问密码（bcrypt 加密） |

#### 成功响应 (200)

返回 `PhotoListItem` 对象。

---

### [PUT] /api/photos/:id - 更新照片 🔒

#### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| caption | string | ❌ | 图片说明 |
| password | string | ❌ | 密码（空字符串=清除密码） |

#### 成功响应 (200)

返回更新后的 `PhotoListItem` 对象。

---

### [DELETE] /api/photos/:id - 删除照片 🔒

#### 描述
删除照片后自动更新所属相册的封面（取最新照片，无照片则封面设为 null）。

#### 成功响应 (200)

```json
{ "ok": true }
```

---

### [POST] /api/photos/:id/verify-password - 验证照片密码

#### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| password | string | ✅ | 待验证的密码 |

#### 成功响应 (200)

```json
{ "success": true }
```

---

## 照片点赞/踩接口

### [GET] /api/photos/:id/likes - 获取点赞信息

#### 请求参数 (Query)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| deviceId | string | ❌ | 设备指纹（传入则返回 liked 状态） |

#### 成功响应 (200)

```json
{
  "count": 10,
  "liked": true
}
```

---

### [POST] /api/photos/:id/likes - 点赞

#### 描述
幂等操作，使用 `INSERT IGNORE` 保证每设备每照片只点赞一次。

#### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| deviceId | string | ✅ | 设备指纹（8-64 字符） |

#### 成功响应 (200)

```json
{
  "count": 11,
  "liked": true
}
```

---

### [GET] /api/photos/:id/dislikes - 获取踩信息

#### 请求参数 (Query)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| deviceId | string | ❌ | 设备指纹 |

#### 成功响应 (200)

```json
{
  "count": 3,
  "disliked": false
}
```

---

### [POST] /api/photos/:id/dislikes - 踩/取消踩

#### 描述
支持 `action` 参数控制踩或取消踩。

#### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| deviceId | string | ✅ | - | 设备指纹（8-64 字符） |
| action | string | ❌ | `"dislike"` | `"dislike"` 踩 / `"undislike"` 取消踩 |

#### 成功响应 (200)

```json
{
  "count": 4,
  "disliked": true
}
```
