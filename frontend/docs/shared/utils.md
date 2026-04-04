# 工具函数清单

## api.ts

- **路径**：`utils/api.ts`
- **职责**：封装所有前端 API 调用，统一使用 Nuxt 的 `$fetch`

### 导出函数

| 函数 | 签名 | 说明 |
|------|------|------|
| `apiLogin` | `(data: LoginRequest) => Promise<LoginResponse>` | 登录 |
| `apiFetchArticles` | `() => Promise<ArticleListResponse>` | 获取文章列表 |
| `apiFetchArticle` | `(id: string) => Promise<ArticleDetail>` | 获取文章详情 |
| `apiCreateArticle` | `(data: CreateArticleRequest) => Promise<ArticleDetail>` | 创建文章 |
| `apiUpdateArticle` | `(id: string, data: UpdateArticleRequest) => Promise<ArticleDetail>` | 更新文章 |
| `apiDeleteArticle` | `(id: string) => Promise<{ ok: boolean }>` | 删除文章 |
| `apiGetArticles` | `(params?: { title?: string }) => Promise<ArticleListResponse>` | 获取文章列表（支持筛选） |
| `apiUploadImage` | `(file: File) => Promise<{ url: string }>` | 上传图片 |
| `apiGetProfile` | `() => Promise<Profile>` | 获取博主资料 |
| `apiUpdateProfile` | `(data: Partial<Profile>) => Promise<Profile>` | 更新博主资料 |
| `apiGetAlbums` | `() => Promise<AlbumListResponse>` | 获取相册集列表 |
| `apiCreateAlbum` | `(data) => Promise<AlbumItem>` | 创建相册集 |
| `apiUpdateAlbum` | `(id: number, data) => Promise<AlbumItem>` | 更新相册集 |
| `apiDeleteAlbum` | `(id: number) => Promise<{ ok: boolean }>` | 删除相册集 |
| `apiGetPhotos` | `(albumId: number) => Promise<PhotoListResponse>` | 获取照片列表 |
| `apiAddPhoto` | `(albumId: number, data) => Promise<PhotoItem>` | 添加照片 |
| `apiDeletePhoto` | `(id: number) => Promise<{ ok: boolean }>` | 删除照片 |
| `apiUpdatePhoto` | `(id: number, data) => Promise<PhotoItem>` | 更新照片 |
| `apiVerifyAlbumPassword` | `(albumId: number, password: string) => Promise<{ success: boolean }>` | 验证相册密码 |
| `apiVerifyPhotoPassword` | `(photoId: number, password: string) => Promise<{ success: boolean }>` | 验证照片密码 |
| `apiToggleArticleLike` | `(articleId: string, deviceId: string, signal?) => Promise<{ liked, likeCount }>` | 文章点赞切换 |
| `apiGetArticleLikeStatus` | `(articleId: string, deviceId: string) => Promise<{ liked, likeCount }>` | 查询点赞状态 |
| `apiGetArticleLikeStatusBatch` | `(ids: string[], deviceId: string) => Promise<{ likedIds: string[] }>` | 批量查询点赞状态 |
| `apiGetMessages` | `(deviceId?: string) => Promise<MessageListResponse>` | 获取留言列表 |
| `apiCreateMessage` | `(data) => Promise<MessageItem>` | 新增留言 |
| `apiUpdateMessage` | `(id: number, data) => Promise<MessageItem>` | 修改留言 |
| `apiDeleteMessage` | `(id: number) => Promise<{ success: boolean }>` | 删除留言（管理员） |

---

## chunkedUpload.ts

- **路径**：`utils/chunkedUpload.ts`
- **职责**：文件分片上传工具，支持进度回调和失败自动清理

### 导出

| 名称 | 类型 | 说明 |
|------|------|------|
| `chunkedUpload` | `(file: File, onProgress?) => Promise<string>` | 主函数，返回上传后的 URL |
| `UploadProgress` | interface | `{ loaded: number, total: number, percent: number }` |

### 实现策略

| 文件大小 | 策略 | 接口 |
|----------|------|------|
| ≤ 1.5MB | 直传（fetch + credentials） | `POST /api/upload` |
| > 1.5MB | 分片上传（1.5MB/片） | `POST /api/upload/chunk` → `POST /api/upload/merge` |

- 分片上传失败时自动调用 `DELETE /api/upload/chunk` 清理临时文件
- 使用 `fetch`（非 `$fetch`）确保 `credentials: 'include'` 携带 Cookie

---

## imageUrl.ts

- **路径**：`utils/imageUrl.ts`
- **职责**：图片 CDN URL 转换

### 导出

| 函数 | 签名 | 说明 |
|------|------|------|
| `toCdnUrl` | `(url: string \| null \| undefined) => string` | 将 `/uploads/xxx` 转为 CDN 完整 URL |

### 转换规则

| 输入 | 输出 |
|------|------|
| `null` / `undefined` / `''` | `''` |
| `http://...` / `https://...` | 原样返回 |
| `/uploads/xxx.jpg` | `https://cdn.fatwill.cloud/uploads/xxx.jpg` |
| 其他 | 原样返回 |

---

## ui.ts

- **路径**：`utils/ui.ts`
- **职责**：统一 UI 工具，自动根据设备类型选择 PC 端（antd）或移动端（自定义）组件

### 导出函数

| 函数 | 签名 | 说明 |
|------|------|------|
| `showSuccess` | `(content: string) => void` | 成功提示 |
| `showError` | `(content: string) => void` | 错误提示 |
| `showInfo` | `(content: string) => void` | 普通信息提示 |
| `showConfirm` | `(options: ConfirmOptions) => void` | 确认弹窗 |

### ConfirmOptions

```typescript
{
  title: string
  content: string
  okText?: string        // 默认 '确定'
  cancelText?: string    // 默认 '取消'
  danger?: boolean       // 是否为危险操作
  onOk: () => void | Promise<void>
}
```

### 跨端适配

| 场景 | PC 端（≥768px） | 移动端（<768px） |
|------|----------------|-----------------|
| Toast | `message.success/error/info` | `MobileToast.success/error/show` |
| 确认弹窗 | `Modal.confirm` | `MobileDialog.confirm` |

---

## mobileUI.ts

- **路径**：`utils/mobileUI.ts`
- **职责**：移动端命令式 UI 组件（仿 antd-mobile 风格）

### 导出

| 名称 | 类型 | 说明 |
|------|------|------|
| `MobileToast` | object | 命令式 Toast 提示 |
| `MobileDialog` | object | 命令式确认弹窗 |

### MobileToast API

| 方法 | 说明 |
|------|------|
| `MobileToast.success(content)` | 成功提示（✓ 图标） |
| `MobileToast.error(content)` | 错误提示（✕ 图标） |
| `MobileToast.show({ content, icon?, duration? })` | 通用提示 |

### MobileDialog API

| 方法 | 说明 |
|------|------|
| `MobileDialog.confirm({ title, content, okText?, cancelText?, danger?, onOk })` | 确认弹窗 |

### 实现方式

- 使用 Vue 3 的 `createApp` + `defineComponent` + `h` 函数动态创建 DOM
- Transition 动画（淡入淡出）
- 自动在 `document.body` 上挂载/卸载
- Dark 模式通过检查 `document.documentElement.classList.contains('dark')` 适配
