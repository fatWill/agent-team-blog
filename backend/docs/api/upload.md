# 上传接口文档

> 所有上传接口均需鉴权（Cookie Token）

## 图片存储说明

- **存储后端**：腾讯云 COS 对象存储
- **访问域名**：`https://assets.fatwill.cloud`
- **存储路径规范**：`upload/YYYYMMDD/时间戳-随机串.ext`
- **完整 URL 示例**：`https://assets.fatwill.cloud/upload/20260405/1743811200000-a1b2c3d4.jpg`
- **支持格式**：jpg / jpeg / png / gif / webp

---

## [POST] /api/upload - 直传图片

### 描述

上传单张图片到腾讯云 COS，返回公开访问 URL。

### 请求头

| Header | 必填 | 说明 |
|--------|------|------|
| Cookie | ✅ | `auth_token=xxx`（httpOnly） |
| Content-Type | ✅ | `multipart/form-data` |

### 请求参数 (FormData)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | File | ✅ | 图片文件（jpg/jpeg/png/gif/webp） |

### 成功响应 (200)

```json
{
  "url": "https://assets.fatwill.cloud/upload/20260405/1743811200000-a1b2c3d4.jpg"
}
```

### 错误响应

| statusCode | statusMessage | 场景 |
|------------|---------------|------|
| 400 | 请选择要上传的文件 | 未选择文件 |
| 400 | 不支持的文件格式，仅支持 jpg/jpeg/png/gif/webp | 格式不支持 |
| 401 | 未授权 | 未登录 |
| 500 | COS 客户端未初始化 | COS 配置缺失 |
| 500 | 上传 COS 失败 | COS 服务异常 |

---

## [POST] /api/upload/chunk - 上传分片

### 描述

大文件分片上传。分片临时保存到服务器本地，待全部分片上传完成后调用合并接口。

### 请求头

| Header | 必填 | 说明 |
|--------|------|------|
| Cookie | ✅ | `auth_token=xxx` |
| Content-Type | ✅ | `multipart/form-data` |

### 请求参数 (FormData)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | File | ✅ | 分片数据 |
| uploadId | string | ✅ | 上传会话 ID（仅允许 `\w-` 字符） |
| chunkIndex | int | ✅ | 当前分片索引（从 0 开始） |
| totalChunks | int | ✅ | 总分片数 |
| filename | string | ✅ | 原始文件名 |

### 成功响应 (200)

```json
{
  "received": 0,
  "total": 5
}
```

### 错误响应

| statusCode | statusMessage | 场景 |
|------------|---------------|------|
| 400 | 缺少必要参数 | 参数不完整 |
| 400 | uploadId 包含非法字符 | 路径穿越防护 |
| 400 | chunkIndex 或 totalChunks 参数无效 | 参数范围错误 |
| 401 | 未授权 | 未登录 |

---

## [POST] /api/upload/merge - 合并分片

### 描述

合并已上传的所有分片，上传完整文件到腾讯云 COS，返回公开访问 URL。合并成功后自动清理服务器本地临时文件。

### 请求头

| Header | 必填 | 说明 |
|--------|------|------|
| Cookie | ✅ | `auth_token=xxx` |
| Content-Type | ✅ | `application/json` |

### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| uploadId | string | ✅ | 上传会话 ID |
| totalChunks | int | ✅ | 总分片数 |
| filename | string | ✅ | 原始文件名（用于提取扩展名） |

### 请求示例

```json
{
  "uploadId": "abc123",
  "totalChunks": 5,
  "filename": "large-photo.jpg"
}
```

### 成功响应 (200)

```json
{
  "url": "https://assets.fatwill.cloud/upload/20260405/1743811200000-a1b2c3d4.jpg"
}
```

### 错误响应

| statusCode | statusMessage | 场景 |
|------------|---------------|------|
| 400 | 缺少必要参数 | 参数不完整 |
| 400 | 不支持的文件格式 | 文件扩展名不在白名单 |
| 400 | 上传会话不存在或已过期 | 分片已清理或未上传 |
| 400 | 分片 N 缺失，请重新上传 | 某个分片文件丢失 |
| 401 | 未授权 | 未登录 |
| 500 | 上传 COS 失败 | COS 服务异常 |

---

## [DELETE] /api/upload/chunk - 取消上传

### 描述

取消分片上传，清理服务器上的临时分片文件。

### 请求头

| Header | 必填 | 说明 |
|--------|------|------|
| Cookie | ✅ | `auth_token=xxx` |
| Content-Type | ✅ | `application/json` |

### 请求参数 (Body - JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| uploadId | string | ✅ | 上传会话 ID |

### 请求示例

```json
{
  "uploadId": "abc123"
}
```

### 成功响应 (200)

```json
{
  "ok": true
}
```

---

## COS 删除说明

当以下操作发生时，后端会自动从 COS 删除对应的图片文件：

| 触发场景 | 接口 | 说明 |
|----------|------|------|
| 删除单张照片 | `DELETE /api/photos/:id` | 异步删除该照片的 COS 对象 |
| 删除相册（级联） | `DELETE /api/albums/:id` | 异步批量删除相册下所有照片的 COS 对象 |

> 删除操作为异步执行，不阻塞 API 响应。若 COS 删除失败会记录日志，不影响数据库删除结果。
> 仅删除 URL 前缀为 `https://assets.fatwill.cloud/` 的对象，旧的本地路径 URL 会被跳过。
