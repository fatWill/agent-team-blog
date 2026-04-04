# 管理后台（/admin）

## 基本信息

| 属性 | 值 |
|------|------|
| 路由 | `/admin` |
| 文件 | `pages/admin.vue`（57KB，大文件） |
| 渲染模式 | CSR（鉴权保护，无 SEO 需求） |
| 鉴权 | `middleware/auth.ts` 路由守卫 |
| 依赖 | auth store, api.ts, Tiptap, ant-design-vue, chunkedUpload, useTheme |

## 功能描述

博主的内容管理中心，需登录后才能访问。通过 Tab 切换四个管理功能区。

## 页面结构

### 顶部导航栏

| 元素 | 位置 | 功能 |
|------|------|------|
| ← 首页 | 左侧 | NuxtLink 跳转 /home |
| 管理后台 | 左侧 | 当前页面标题 |
| 🌙/☀️ | 右侧 | 主题切换按钮 |
| 退出登录 | 右侧 | 清除认证状态 + 跳转登录页 |

### Tab 导航

| Tab | 功能 | 主要操作 |
|-----|------|---------|
| 写文章 | Tiptap 富文本编辑器 | 标题/摘要/封面图/正文编辑 → 发布 |
| 文章管理 | 文章列表 | 搜索（300ms 防抖）/ 编辑 / 删除 |
| 相册管理 | 相册集 + 照片管理 | 创建/编辑/删除相册 / 上传/删除照片 / 密码设置 |
| 留言管理 | 留言列表 | 查看 / 删除 |

## 鉴权保护

访问 `/admin` 时，`middleware/auth.ts` 会：

1. 检查内存中的登录状态
2. 检查 Cookie 是否存在
3. 调用 `authStore.checkAuth()` 验证 Token
4. 未登录跳转 `/login?redirect=/admin`

## 文章编辑器

### Tiptap 配置

- 扩展：StarterKit + Image + Link
- 输出格式：JSON（存储在 MySQL `articles.content` 字段）
- 图片插入：通过 `chunkedUpload()` 上传后插入 URL

### 编辑流程

```
新建文章：填写表单 → apiCreateArticle() → 成功提示 → 重置表单
编辑文章：从文章管理 Tab 点击编辑 → 加载数据到编辑器 → apiUpdateArticle() → 成功提示
```

## 相册管理

### 相册集操作

- **创建**：Modal 弹窗 → 输入名称/描述/密码 → `apiCreateAlbum()`
- **编辑**：Modal 弹窗 → 修改名称/描述/密码 → `apiUpdateAlbum()`
- **删除**：`showConfirm()` → `apiDeleteAlbum()`（级联删除所有照片）

### 照片操作

- **上传**：选择文件 → 乐观更新（占位卡片 + 进度条）→ `chunkedUpload()` → `apiAddPhoto()`
- **删除**：照片 hover 显示删除按钮 → `apiDeletePhoto()`
- **密码**：照片 hover 显示密码设置按钮 → Modal 输入密码 → `apiUpdatePhoto()`
- **布局**：3 列网格

## 操作反馈

| 操作 | 成功 | 失败 |
|------|------|------|
| 创建/更新 | `showSuccess('操作成功')` | `showError(errorMessage)` |
| 删除 | `showSuccess('删除成功')` | `showError('删除失败')` |
| 上传 | 进度条 → 替换为真实图片 | 裂图 + 红色提示 |

## Meta

```typescript
useHead({ title: '管理后台 - fatwill' })
```
