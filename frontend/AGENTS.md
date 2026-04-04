# AGENTS.md - 项目中枢索引

## 项目概述

fatwill 个人博客，基于 Nuxt 3 全栈架构（前端 SSR + 后端 Nitro Server API），支持 Dark/Light 主题切换、文章管理、相册管理、留言板等功能。

## 技术栈速览

| 类别 | 技术 | 版本 |
|------|------|------|
| 框架 | Nuxt 3 | ^3.16.1 |
| 视图层 | Vue 3 | ^3.5.13 |
| 语言 | TypeScript | ^5.8.2 |
| 样式 | Tailwind CSS（@nuxtjs/tailwindcss） | ^6.13.2 |
| 状态管理 | Pinia（@pinia/nuxt） | ^2.3.0 / ^0.9.0 |
| 富文本编辑器 | Tiptap（vue-3 + starter-kit + image + link） | ^3.21.0 |
| UI 组件库 | Ant Design Vue（PC 端） | ^4.2.6 |
| 虚拟滚动 | vue-virtual-scroller | ^2.0.0-beta.8 |
| Canvas 文字排版 | @chenglou/pretext | ^0.0.4 |
| 数据库 | MySQL 8.0（mysql2 驱动） | ^3.20.0 |
| 缓存 | Redis（ioredis） | ^5.10.1 |
| 图片处理 | sharp | ^0.34.5 |
| 加密 | bcryptjs | ^3.0.3 |
| 进度条 | nprogress | ^0.2.0 |
| 包管理 | npm | — |

## 最高优先级读取

1. `docs/eslint-rules.md` — 代码规范与命名约定（最高优先级）
2. `docs/architecture.md` — 整体架构概览与目录结构说明

## 当前目录结构（实际现状）

> ⚠️ **注意**：当前项目**未采用 Feature-First 结构**，所有页面逻辑内聚在各 `.vue` 单文件组件中（尤其是 `home.vue` 和 `admin.vue` 体量较大），待后续重构为 Feature-First 模块化结构。

```
frontend/
├── AGENTS.md                          # 中枢索引（本文件）
├── docs/                              # 工程文档体系
│   ├── architecture.md                # 整体架构概览
│   ├── eslint-rules.md                # 代码规范
│   ├── features/                      # 领域模块文档
│   │   ├── article/README.md
│   │   ├── album/README.md
│   │   ├── auth/README.md
│   │   ├── admin/README.md
│   │   └── guestbook/README.md
│   ├── pages/                         # 页面级文档
│   │   ├── home.md
│   │   ├── blog.md
│   │   └── admin.md
│   └── shared/                        # 共享层文档
│       ├── components.md
│       ├── composables.md
│       └── utils.md
├── app.vue                            # 应用入口（主题 class 注入）
├── nuxt.config.ts                     # Nuxt 配置
├── tailwind.config.ts                 # Tailwind 配置（dark mode: class）
├── package.json                       # 依赖管理
├── tsconfig.json                      # TypeScript 配置
├── assets/
│   └── styles/
│       └── tailwind.css               # Tailwind 入口样式
├── components/
│   └── AppLoading.vue                 # 跨端 Loading 组件
├── composables/
│   ├── useTheme.ts                    # 主题切换（cookie + Redis 持久化）
│   └── useDevice.ts                   # 设备检测（移动端/PC端）
├── layouts/
│   └── default.vue                    # 默认布局（主题初始化）
├── middleware/
│   └── auth.ts                        # 鉴权路由守卫（保护 /admin）
├── pages/
│   ├── index.vue                      # 酷炫首页（/ — Canvas + pretext 代码滚动）
│   ├── home.vue                       # 博客主页（/home — 文章/相册/留言板等 Tab）
│   ├── login.vue                      # 登录页（/login）
│   ├── admin.vue                      # 管理后台（/admin）
│   └── articles/
│       └── [id].vue                   # 文章详情页（/articles/:id）
├── plugins/
│   ├── antd.client.ts                 # Ant Design Vue 客户端插件
│   ├── nprogress.client.ts            # NProgress 路由进度条
│   ├── pinia-hydration-fix.ts         # Pinia SSR hydration 修复
│   └── vue-virtual-scroller.client.ts # 虚拟滚动客户端插件
├── public/
│   └── avatar.png                     # 博主头像
├── stores/
│   └── auth.ts                        # 认证状态管理（Pinia）
├── types/
│   └── index.ts                       # 全局 TypeScript 类型定义
├── utils/
│   ├── api.ts                         # API 请求封装（所有前端接口调用）
│   ├── chunkedUpload.ts               # 分片上传工具
│   ├── imageUrl.ts                    # CDN URL 转换工具（toCdnUrl）
│   ├── ui.ts                          # 统一 UI 工具（跨端 Toast/Dialog）
│   └── mobileUI.ts                    # 移动端 Toast/Dialog 命令式组件
└── server/                            # Nuxt Server（Nitro）
    ├── api/                           # API 路由
    │   ├── auth/                      # 鉴权（login/check/logout）
    │   ├── articles/                  # 文章 CRUD + 点赞
    │   ├── albums/                    # 相册集 CRUD + 照片 + 密码验证
    │   ├── photos/                    # 照片管理 + 点赞/踩 + 密码验证
    │   ├── messages/                  # 留言板 CRUD
    │   ├── profile/                   # 个人资料
    │   ├── upload/                    # 分片上传（chunk/merge/delete）
    │   ├── upload.post.ts             # 小文件直传
    │   ├── changelog.get.ts           # 更新日志
    │   ├── theme.get.ts               # 获取主题偏好
    │   └── theme.post.ts              # 保存主题偏好
    ├── utils/                         # 服务端工具层
    │   ├── db.ts                      # MySQL 连接池（进程级单例）
    │   ├── redis.ts                   # Redis 连接
    │   ├── auth.ts                    # Token 鉴权工具
    │   ├── articles.ts                # 文章 DAO
    │   ├── albums.ts                  # 相册 DAO
    │   ├── messages.ts                # 留言 DAO
    │   ├── changelog.ts               # 更新日志 DAO
    │   ├── profile.ts                 # 个人资料 DAO
    │   └── rateLimit.ts               # IP 限频工具
    └── data/
        ├── articles.json              # 旧数据文件（已废弃，数据在 MySQL）
        └── changelog.json             # 旧数据文件（已废弃，数据在 MySQL）
```

## 领域模块注册表

> 当前项目所有业务逻辑内聚在页面级 `.vue` 文件中，尚未拆分为 Feature-First 模块。下表按业务领域归类记录现状。

| 领域 | 主要代码位置 | 状态 | 说明 | 待重构 |
|------|-------------|------|------|--------|
| 酷炫首页 | `pages/index.vue` | ✅ 已实现 | Canvas + @chenglou/pretext 代码滚动背景、拖拽按钮 | 可提取为 `features/landing/` |
| 文章 | `pages/home.vue`（文章 Tab）、`pages/articles/[id].vue` | ✅ 已实现 | 文章列表、详情（Tiptap 渲染）、点赞、虚拟滚动 | 可提取为 `features/article/` |
| 相册 | `pages/home.vue`（生活 Tab） | ✅ 已实现 | 相册集、照片管理、灯箱预览、密码保护、点赞/踩 | 可提取为 `features/album/` |
| 留言板 | `pages/home.vue`（留言板 Tab） | ✅ 已实现 | 留言 CRUD、编辑限制、相对时间 | 可提取为 `features/guestbook/` |
| 更新日志 | `pages/home.vue`（更新日志 Tab） | ✅ 已实现 | Git-log 风格时间轴 | 可提取为 `features/changelog/` |
| Agent Team | `pages/home.vue`（Agent Team Tab） | ✅ 已实现 | 4 个 AI Agent 卡片展示 | 可提取为 `features/agent-team/` |
| 鉴权 | `stores/auth.ts`、`middleware/auth.ts`、`server/api/auth/` | ✅ 已实现 | Cookie Token 72h 滚动续期、SSR Cookie 转发 | 可提取为 `features/auth/` |
| 管理后台 | `pages/admin.vue` | ✅ 已实现 | 文章管理、相册管理、个人资料、留言管理 | 可提取为 `features/admin/` |

## 共享层注册表

### Composables（`composables/`）

| 模块 | 文件 | 说明 |
|------|------|------|
| 主题切换 | `composables/useTheme.ts` | Dark/Light 切换，Cookie + Redis 持久化，SSR 零闪烁 |
| 设备检测 | `composables/useDevice.ts` | 响应式 `isMobile` ref + 命令式 `isMobileDevice()` 函数 |

### Stores（`stores/`）

| 模块 | 文件 | 说明 |
|------|------|------|
| 认证状态 | `stores/auth.ts` | `isLoggedIn`、`username`、`checkAuth()`（SSR Cookie 转发）、`logout()` |

### Utils（`utils/`）

| 模块 | 文件 | 说明 |
|------|------|------|
| API 封装 | `utils/api.ts` | 所有前端 API 调用函数（文章/相册/照片/留言/认证/资料等） |
| 分片上传 | `utils/chunkedUpload.ts` | ≤1.5MB 直传，>1.5MB 分片上传（进度回调、失败自动清理） |
| CDN URL | `utils/imageUrl.ts` | `toCdnUrl()` — 将 `/uploads/xxx` 转为 `https://cdn.fatwill.cloud/uploads/xxx` |
| 统一 UI | `utils/ui.ts` | `showSuccess()`、`showError()`、`showInfo()`、`showConfirm()` — PC 用 antd，移动端用自定义组件 |
| 移动端 UI | `utils/mobileUI.ts` | `MobileToast`（success/error/show）、`MobileDialog`（confirm）— 仿 antd-mobile 风格 |

### Components（`components/`）

| 组件 | 文件 | 说明 |
|------|------|------|
| AppLoading | `components/AppLoading.vue` | 跨端 Loading 组件，PC 用 ASpin，移动端用自定义 spinner |

### Plugins（`plugins/`）

| 插件 | 文件 | 运行环境 | 说明 |
|------|------|---------|------|
| Ant Design Vue | `plugins/antd.client.ts` | 仅客户端 | 注册 message/modal/spin 样式，配置 message 全局参数 |
| NProgress | `plugins/nprogress.client.ts` | 仅客户端 | 路由导航进度条（beforeEach/afterEach/onError） |
| Pinia Hydration Fix | `plugins/pinia-hydration-fix.ts` | 仅服务端 | JSON 往返清洗 payload，修复 mysql2 RowDataPacket 无原型链问题 |
| Virtual Scroller | `plugins/vue-virtual-scroller.client.ts` | 仅客户端 | 注册 vue-virtual-scroller 全局组件 |

### Types（`types/`）

| 文件 | 说明 |
|------|------|
| `types/index.ts` | 全局类型定义：ArticleListItem、ArticleDetail、AlbumItem、PhotoItem、MessageItem、Profile、ChangelogItem、LoginRequest/Response、TabItem 等 |

## 页面路由映射

| 页面 | 路由 | 文件 | 依赖模块 | 说明 |
|------|------|------|----------|------|
| 酷炫首页 | `/` | `pages/index.vue`（17KB） | @chenglou/pretext | Canvas 代码滚动背景 + 拖拽按钮 |
| 博客主页 | `/home` | `pages/home.vue`（82KB） | useTheme, useDevice, api, types, vue-virtual-scroller, antd | 6 个 Tab：文章/生活/小工具·小游戏/Agent Team/更新日志/留言板 |
| 文章详情 | `/articles/:id` | `pages/articles/[id].vue`（8KB） | useTheme, api, Tiptap | Tiptap 只读渲染 + 点赞 |
| 登录 | `/login` | `pages/login.vue`（4KB） | auth store, api | 账号密码登录，成功跳转 redirect 或 /home |
| 管理后台 | `/admin` | `pages/admin.vue`（57KB） | auth store, api, Tiptap, antd, chunkedUpload | 4 个管理 Tab：文章/相册/个人资料/留言管理 |

## 后端 API（Nuxt Server / Nitro）

### 鉴权

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| POST | `/api/auth/login` | ❌ | 账号密码验证，写入 httpOnly cookie（72h TTL） |
| GET | `/api/auth/check` | Cookie | 验证 token 有效性，滚动续期 lastActiveAt |
| POST | `/api/auth/logout` | Cookie | 清除 token |

### 文章

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/articles` | ❌ | 文章列表（不含 content），支持 `?title=xxx` 筛选 |
| GET | `/api/articles/:id` | ❌ | 文章详情（含 Tiptap JSON content） |
| POST | `/api/articles` | ✅ | 创建文章（title/summary/coverImage/content） |
| PUT | `/api/articles/:id` | ✅ | 更新文章（部分更新） |
| DELETE | `/api/articles/:id` | ✅ | 物理删除文章 |
| POST | `/api/articles/:id/like` | ❌ | 点赞/取消点赞（deviceId 切换） |
| GET | `/api/articles/:id/like-status` | ❌ | 查询单篇点赞状态 |
| GET | `/api/articles/like-status-batch` | ❌ | 批量查询点赞状态 |

### 相册

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/albums` | ❌ | 相册集列表（含 photoCount、coverUrl、hasPassword） |
| POST | `/api/albums` | ✅ | 创建相册集（name/description/password） |
| PUT | `/api/albums/:id` | ✅ | 更新相册集 |
| DELETE | `/api/albums/:id` | ✅ | 级联删除相册及所有照片 |
| GET | `/api/albums/:id/photos` | ❌ | 获取相册下所有照片 |
| POST | `/api/albums/:id/photos` | ✅ | 添加照片并自动更新封面 |
| POST | `/api/albums/:id/verify-password` | ❌ | 验证相册集密码 |

### 照片

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| PUT | `/api/photos/:id` | ✅ | 更新照片信息（caption/password） |
| DELETE | `/api/photos/:id` | ✅ | 删除照片并自动更新相册封面 |
| POST | `/api/photos/:id/likes` | ❌ | 照片点赞 |
| GET | `/api/photos/:id/likes` | ❌ | 查询照片点赞状态 |
| POST | `/api/photos/:id/dislikes` | ❌ | 照片踩 |
| GET | `/api/photos/:id/dislikes` | ❌ | 查询照片踩状态 |
| POST | `/api/photos/:id/verify-password` | ❌ | 验证照片密码 |

### 留言板

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/messages` | ❌ | 留言列表（支持 `?deviceId=xxx` 标记自己的留言） |
| POST | `/api/messages` | ❌ | 新增留言（IP 限频） |
| PUT | `/api/messages/:id` | ❌ | 修改留言（每日一次限制） |
| DELETE | `/api/messages/:id` | ✅ | 管理员删除留言 |

### 其他

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/profile` | ❌ | 获取博主资料（avatar/bio） |
| PUT | `/api/profile` | ✅ | 更新博主资料 |
| GET | `/api/changelog` | ❌ | 获取所有版本更新日志 |
| POST | `/api/upload` | ✅ | 小文件直传（multipart/form-data） |
| POST | `/api/upload/chunk` | ✅ | 上传单个分片 |
| POST | `/api/upload/merge` | ✅ | 合并所有分片 |
| DELETE | `/api/upload/chunk` | ✅ | 清理临时分片文件 |
| GET | `/api/theme` | ❌ | 从 Redis 获取主题偏好 |
| POST | `/api/theme` | ❌ | 保存主题偏好到 Redis |

## 服务端工具层（`server/utils/`）

| 文件 | 说明 |
|------|------|
| `db.ts` | MySQL 连接池（mysql2，进程级单例，从 runtimeConfig 读取配置） |
| `redis.ts` | Redis 连接（ioredis，从 runtimeConfig 读取配置） |
| `auth.ts` | Token 鉴权工具（生成/验证/续期，httpOnly cookie，72h TTL） |
| `articles.ts` | 文章 DAO（CRUD + 点赞计数） |
| `albums.ts` | 相册 DAO（相册集 CRUD + 照片 CRUD + 封面自动更新 + 密码保护） |
| `messages.ts` | 留言 DAO（CRUD + deviceId 关联） |
| `changelog.ts` | 更新日志 DAO（查询 changelogs 表） |
| `profile.ts` | 个人资料 DAO（读写 profile 表） |
| `rateLimit.ts` | IP 限频工具（内存 Map 实现） |

## 数据库表注册表

| 表名 | 说明 | 主键 | 索引 |
|------|------|------|------|
| `articles` | 文章表 | `id` (varchar(36) UUID) | `idx_created_at` |
| `changelogs` | 更新日志表 | `id` (bigint auto_increment) | `uk_version`, `idx_date` |
| `auth_tokens` | 登录 Token 表 | `token` (varchar(64)) | `idx_username` |
| `profile` | 博主个人资料表 | `id` (int, 固定为 1) | — |
| `albums` | 相册集表 | `id` (int auto_increment) | — |
| `photos` | 照片表 | `id` (int auto_increment) | `idx_album_id` |
| `article_likes` | 文章点赞表 | `id` (int auto_increment) | `uq_article_device`, `idx_article_id` |
| `messages` | 留言板表 | `id` (int auto_increment) | `uq_device_id` |

## 环境配置

数据库连接通过 `.env` 文件注入（Nuxt runtimeConfig 读取）：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `DB_HOST` | MySQL 主机地址 | `localhost` |
| `DB_PORT` | MySQL 端口 | `3306` |
| `DB_USER` | 数据库用户名 | `root` |
| `DB_PASSWORD` | 数据库密码 | — |
| `DB_NAME` | 数据库名 | `blog` |
| `REDIS_HOST` | Redis 主机地址 | `127.0.0.1` |
| `REDIS_PORT` | Redis 端口 | `6379` |
| `REDIS_PASSWORD` | Redis 密码 | — |

## 版本号规范

格式：`X.Y.Z`（Semantic Versioning）

| 位置 | 含义 | 何时递增 |
|------|------|---------|
| X（主版本） | 重大重构或破坏性变更 | 架构级改动 |
| Y（次版本） | 常规功能迭代 | **默认每次发版递增此位** |
| Z（补丁版本） | 紧急 bug 修复 | 仅修复线上问题 |

## 更新日志规范

- 存储在 MySQL `changelogs` 表中
- 每条最多 **20 字**，最多 **5 条**
- 内容风格：Emoji 前缀 + 简短描述，如 `"🎉 个人博客正式上线"`

## 变更日志（最近重要变更）

- 2026-04-04: 默认 Dark 模式 + 修正 GitHub 链接（fatwillzeng → fatwill）
- 2026-04-04: SSR 阶段 Cookie 转发修复重启后需重新登录
- 2026-04-04: 新增「Agent Team」Tab，展示 4 个 AI Agent 卡片
- 2026-04-04: Tab 导航添加 Emoji 图标前缀 + 美化样式
- 2026-04-03: 首页重写为 Canvas + @chenglou/pretext 代码滚动背景 + 拖拽按钮
- 2026-04-02: 文章点赞 + 留言板 Tab 功能（v1.9.0）
- 2026-04-02: 相册/照片密码保护功能（v1.8.0）
- 2026-04-01: CDN 图片 URL 转换、ant-design-vue 集成、虚拟滚动优化
- 2026-04-01: 首屏 CSR→SSR 改造（useAsyncData 预取 + Pinia hydration 修复）
- 2026-03-31: 数据从 JSON 文件迁移到 MySQL 8.0、相册功能、分片上传、Token 持久化
- 2026-03-30: 全面重写博客前端（Dark/Light 主题、Tab 导航、Tiptap 编辑器）
