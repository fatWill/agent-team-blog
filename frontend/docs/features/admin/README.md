# 管理后台领域

## 概述

管理后台（`/admin`）是博主的内容管理中心，需登录后才能访问。包含文章管理、相册管理、个人资料管理和留言管理四个功能模块。

## 当前代码位置

| 功能 | 文件 | 说明 |
|------|------|------|
| 管理后台 UI | `pages/admin.vue`（57KB，大文件） | 所有管理功能内聚在单文件中 |
| 鉴权保护 | `middleware/auth.ts` | 路由守卫，未登录跳转 /login |
| API 封装 | `utils/api.ts` | 所有管理操作的 API 调用 |
| 分片上传 | `utils/chunkedUpload.ts` | 图片上传工具 |

## 功能模块

### Tab 结构

管理后台通过 Tab 切换四个功能区：

| Tab | 功能 | 说明 |
|-----|------|------|
| 写文章 | 文章编辑器 | Tiptap 富文本编辑器 + 封面图上传 + 预览 |
| 文章管理 | 文章列表管理 | 搜索（300ms 防抖）、编辑、删除 |
| 相册管理 | 相册集和照片管理 | 创建/编辑/删除相册集、批量上传照片、照片密码设置 |
| 留言管理 | 留言列表管理 | 查看所有留言、删除留言 |

### 文章编辑器

- **编辑器**：Tiptap（@tiptap/vue-3 + starter-kit + image + link）
- **封面图**：支持上传封面图片，预览显示
- **保存**：创建 `POST /api/articles` 或更新 `PUT /api/articles/:id`
- **操作反馈**：成功用 `showSuccess()`，失败用 `showError()`

### 文章管理

- **搜索**：标题关键词搜索，300ms 防抖
- **编辑**：点击文章进入编辑模式（复用写文章 Tab 的 Tiptap 编辑器）
- **删除**：`showConfirm()` 确认弹窗 → `apiDeleteArticle()`

### 相册管理

- **相册集操作**：
  - 创建：Modal 弹窗输入名称/描述/密码
  - 编辑：Modal 弹窗修改名称/描述/密码
  - 删除：确认弹窗 → 级联删除相册及所有照片
- **照片操作**：
  - 上传：选择文件 → 乐观更新（占位卡片 + 进度条）→ `chunkedUpload()` → `apiAddPhoto()`
  - 删除：照片卡片 hover 显示删除按钮
  - 密码设置：照片 hover 显示密码设置按钮

### 留言管理

- **列表展示**：显示所有留言（昵称、内容、时间）
- **删除**：确认弹窗 → `apiDeleteMessage(id)`

## 顶部导航

| 元素 | 说明 |
|------|------|
| ← 首页 | NuxtLink 跳转 /home |
| 管理后台 | 当前页面标题 |
| 🌙/☀️ 按钮 | 主题切换（useTheme） |
| 退出登录 | 调用 authStore.logout() + 跳转登录页 |

## UI 组件使用

| 场景 | 组件 |
|------|------|
| Loading | `AppLoading` 组件 |
| Toast 提示 | `showSuccess()` / `showError()` — 跨端适配 |
| 确认弹窗 | `showConfirm()` — 跨端适配 |
| 表单弹窗 | ant-design-vue Modal（PC 端） |

## 待优化方向

1. `admin.vue` 文件体量过大（57KB），建议拆分为子组件：
   - `AdminArticleEditor.vue` — 文章编辑器
   - `AdminArticleList.vue` — 文章管理列表
   - `AdminAlbumManager.vue` — 相册管理
   - `AdminMessageManager.vue` — 留言管理
2. 个人资料管理功能可独立为 Tab
