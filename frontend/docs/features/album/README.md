# 相册领域

## 概述

相册功能是博客「生活」Tab 的核心，支持相册集管理、照片浏览（灯箱预览）、密码保护、点赞/踩交互。

## 核心实体

### AlbumItem

```typescript
interface AlbumItem {
  id: number
  name: string
  description: string | null
  coverUrl: string | null      // 自动取最新一张照片作为封面
  photoCount: number
  hasPassword: boolean          // 是否设置了密码保护
  createdAt: string
  updatedAt: string
}
```

### PhotoItem

```typescript
interface PhotoItem {
  id: number
  albumId: number
  url: string
  caption: string | null
  hasPassword: boolean
  createdAt: string
  updatedAt: string
  likes?: number               // 前端运行时附加
  liked?: boolean              // 当前设备是否已点赞
  dislikes?: number            // 前端运行时附加
  disliked?: boolean           // 当前设备是否已踩
}
```

## 当前代码位置

| 功能 | 文件 | 说明 |
|------|------|------|
| 相册集列表 | `pages/home.vue`（生活 Tab） | 网格布局（移动端 2 列 / PC 端 5 列） |
| 照片浏览 | `pages/home.vue`（生活 Tab） | 按年月分组，虚拟滚动 |
| 灯箱预览 | `pages/home.vue` | 全屏遮罩，左右切换/键盘/触摸滑动，双指缩放(0.5x-4x) |
| 密码保护 | `pages/home.vue` | 未解锁显示模糊遮罩，密码验证弹窗，sessionStorage 缓存解锁状态 |
| 点赞/踩 | `pages/home.vue` | 灯箱中爱心按钮（右上角）+ 👎踩按钮 + 飘动动画 |
| 相册管理 | `pages/admin.vue`（相册管理 Tab） | 创建/编辑/删除相册集，批量上传照片 |
| API 封装 | `utils/api.ts` | `apiGetAlbums()`、`apiCreateAlbum()`、`apiUpdateAlbum()`、`apiDeleteAlbum()`、`apiGetPhotos()`、`apiAddPhoto()`、`apiDeletePhoto()`、`apiUpdatePhoto()`、`apiVerifyAlbumPassword()`、`apiVerifyPhotoPassword()` |
| 类型定义 | `types/index.ts` | `AlbumItem`、`AlbumListResponse`、`PhotoItem`、`PhotoListResponse` |
| 服务端 API | `server/api/albums/`、`server/api/photos/` | 相册 + 照片接口 |
| 服务端 DAO | `server/utils/albums.ts` | 相册集/照片 CRUD + 封面自动更新 + bcrypt 密码哈希 |

## API 接口

### 相册集

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/albums` | ❌ | 相册集列表 |
| POST | `/api/albums` | ✅ | 创建相册集（支持 password） |
| PUT | `/api/albums/:id` | ✅ | 更新相册集 |
| DELETE | `/api/albums/:id` | ✅ | 级联删除 |
| GET | `/api/albums/:id/photos` | ❌ | 照片列表 |
| POST | `/api/albums/:id/photos` | ✅ | 添加照片 |
| POST | `/api/albums/:id/verify-password` | ❌ | 验证相册密码 |

### 照片

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| PUT | `/api/photos/:id` | ✅ | 更新照片（caption/password） |
| DELETE | `/api/photos/:id` | ✅ | 删除照片 |
| POST | `/api/photos/:id/likes` | ❌ | 照片点赞 |
| GET | `/api/photos/:id/likes` | ❌ | 查询点赞状态 |
| POST | `/api/photos/:id/dislikes` | ❌ | 照片踩 |
| GET | `/api/photos/:id/dislikes` | ❌ | 查询踩状态 |
| POST | `/api/photos/:id/verify-password` | ❌ | 验证照片密码 |

## 业务流程

### 相册浏览流程

```
1. 切换到「生活」Tab
2. → apiGetAlbums() 获取相册集列表
3. → 网格展示（封面 + 名称 + 照片数 + 锁图标）
4. → 点击相册集
5. → 如果 hasPassword：弹出密码验证弹窗
6.   → apiVerifyAlbumPassword() 验证
7.   → 成功：sessionStorage 缓存 unlocked_album_{id}
8. → apiGetPhotos(albumId) 获取照片列表
9. → 照片按年月分组，虚拟滚动渲染
```

### 灯箱预览流程

```
1. 点击照片缩略图
2. → 检查照片是否有密码保护
3.   → 有且未解锁：弹出密码验证
4. → 打开全屏灯箱遮罩
5. → 支持：左右箭头/键盘/触摸滑动切换
6. → 支持：双指捏合缩放（0.5x-4x）
7. → 支持：PC 端鼠标拖动平移（放大模式）
8. → 右上角爱心按钮 + 👎踩按钮
9. → Escape 关闭
10. → 切换到下一张时自动检查目标照片密码状态
```

### 照片上传流程（管理后台）

```
1. 选择文件
2. → 乐观更新：立即插入占位卡片（loading + 进度条）
3. → chunkedUpload(file, onProgress)
4.   → ≤1.5MB：直传 /api/upload
5.   → >1.5MB：分片上传 → 合并
6. → 成功：apiAddPhoto() 保存到数据库
7. → 替换占位卡片为真实图片
8. → 失败：显示裂图 + 红色提示
```

## 密码保护机制

- **存储**：bcrypt 哈希存储在 `albums.password` / `photos.password` 字段
- **查询**：列表接口只返回 `hasPassword: boolean`，不暴露哈希值
- **验证**：`verify-password` 接口使用 `bcryptjs.compare()` 校验
- **缓存**：客户端 `sessionStorage` 缓存解锁状态（`unlocked_album_{id}` / `unlocked_photo_{id}`），关闭标签页即失效
