# AGENTS.md - 项目中枢索引

## 项目概述

fatwill 个人博客，基于 Nuxt 3 前端 SSR + Go 后端 API 架构。前端 Nitro Server 层为纯透传代理，所有业务逻辑和数据存储由 Go 后端（`http://127.0.0.1:8080`）处理。支持 Dark/Light 主题切换、文章管理、相册管理、留言板等功能。

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
| 进度条 | nprogress | ^0.2.0 |
| 后端 API | Go 后端（透传代理） | http://127.0.0.1:8080 |
| 包管理 | npm | — |

## 最高优先级读取

1. `docs/eslint-rules.md` — 代码规范与命名约定（最高优先级）
2. `docs/architecture.md` — 整体架构概览与目录结构说明

## 目录结构（Feature-First + Nuxt 约定混合架构）

> 已完成 Feature-First 重构。类型定义和 API 服务层已按领域拆分到 `features/` 目录，通用工具和 composables 已移至 `shared/` 目录。
> Nuxt 约定目录（`components/`、`composables/`、`utils/`）保留为**向后兼容 re-export 层**，确保 Nuxt 自动导入机制正常工作。
> 页面级 `.vue` 文件（`home.vue` 82KB、`admin.vue` 57KB）的组件拆分为后续优化方向。

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
│
├── features/                          # 🏗️ Feature-First 领域模块
│   ├── article/                       # 文章领域
│   │   ├── types.ts                   #   类型定义（ArticleListItem、ArticleDetail 等）
│   │   ├── api.ts                     #   API 服务层（CRUD + 点赞）
│   │   └── index.ts                   #   模块公共导出入口
│   ├── album/                         # 相册领域
│   │   ├── types.ts                   #   类型定义（AlbumItem、PhotoItem 等）
│   │   ├── api.ts                     #   API 服务层（相册/照片 CRUD + 密码验证）
│   │   └── index.ts                   #   模块公共导出入口
│   ├── auth/                          # 鉴权领域
│   │   ├── types.ts                   #   类型定义（LoginRequest、LoginResponse）
│   │   ├── api.ts                     #   API 服务层（login）
│   │   └── index.ts                   #   模块公共导出入口
│   ├── guestbook/                     # 留言板领域
│   │   ├── types.ts                   #   类型定义（MessageItem、MessageListResponse）
│   │   ├── api.ts                     #   API 服务层（留言 CRUD）
│   │   └── index.ts                   #   模块公共导出入口
│   ├── changelog/                     # 更新日志领域
│   │   ├── types.ts                   #   类型定义（ChangelogItem、ChangelogResponse）
│   │   └── index.ts                   #   模块公共导出入口
│   └── admin/                         # 管理后台领域（预留）
│
├── shared/                            # 🔗 跨领域共享层
│   ├── types.ts                       # 通用类型（NavItem、TabItem、Profile）
│   ├── components/
│   │   └── AppLoading.vue             # 跨端 Loading 组件（实际实现）
│   ├── composables/
│   │   ├── useTheme.ts                # 主题切换（实际实现）
│   │   └── useDevice.ts               # 设备检测（实际实现）
│   └── utils/
│       ├── api.ts                     # 通用 API（profile、upload）
│       ├── ui.ts                      # 统一 UI 工具（实际实现）
│       ├── mobileUI.ts                # 移动端 Toast/Dialog（实际实现）
│       ├── imageUrl.ts                # CDN URL 转换（实际实现）
│       └── chunkedUpload.ts           # 分片上传（实际实现）
│
├── app.vue                            # 应用入口（主题 class 注入）
├── nuxt.config.ts                     # Nuxt 配置
├── tailwind.config.ts                 # Tailwind 配置（dark mode: class）
├── package.json                       # 依赖管理
├── tsconfig.json                      # TypeScript 配置
├── assets/
│   └── styles/
│       └── tailwind.css               # Tailwind 入口样式
├── components/                        # Nuxt 约定（re-export → shared/components/）
│   └── AppLoading.vue                 #   → shared/components/AppLoading.vue
├── composables/                       # Nuxt 约定（re-export → shared/composables/）
│   ├── useTheme.ts                    #   → shared/composables/useTheme.ts
│   └── useDevice.ts                   #   → shared/composables/useDevice.ts
├── utils/                             # Nuxt 约定（re-export → features/*/api.ts + shared/utils/）
│   ├── api.ts                         #   → features/*/api.ts + shared/utils/api.ts
│   ├── ui.ts                          #   → shared/utils/ui.ts
│   ├── mobileUI.ts                    #   → shared/utils/mobileUI.ts
│   ├── imageUrl.ts                    #   → shared/utils/imageUrl.ts
│   └── chunkedUpload.ts              #   → shared/utils/chunkedUpload.ts
├── types/                             # Nuxt 约定（re-export → features/*/types.ts + shared/types.ts）
│   └── index.ts                       #   → 各 feature 类型聚合导出
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
└── server/                            # Nuxt Server（Nitro）
    ├── api/                           # API 路由
    │   ├── auth/                      # 鉴权（login/check/logout）
    │   ├── articles/                  # 文章 CRUD + 点赞
    │   ├── albums/                    # 相册集 CRUD + 照片 + 密码验证
    │   ├── photos/                    # 照片管理 + 点赞/踩 + 密码验证
    │   ├── messages/                  # 留言板 CRUD
    │   ├── pv/                        # PV/UV 数据统计（record/overview/trend/top-pages/logs）
    │   ├── profile/                   # 个人资料
    │   ├── upload/                    # 分片上传（chunk/merge/delete）
    │   ├── upload.post.ts             # 小文件直传
    │   ├── changelog.get.ts           # 更新日志
    │   ├── theme.get.ts               # 获取主题偏好
    │   └── theme.post.ts              # 保存主题偏好
    ├── utils/                         # 服务端工具层
    │   └── proxy.ts                   # 通用代理工具（proxyToBackend + proxyRawToBackend）
    └── .env                           # 环境变量（BACKEND_URL）
```

## 领域模块注册表

> 类型定义和 API 服务层已按领域拆分到 `features/` 目录。页面级 UI 组件仍内聚在 `pages/*.vue` 中，后续可进一步拆分。

| 领域 | Feature 路径 | 包含内容 | 状态 | 页面级 UI |
|------|-------------|---------|------|----------|
| 文章 | `features/article/` | types.ts + api.ts + index.ts | ✅ 已拆分 | `pages/home.vue`（文章 Tab）、`pages/articles/[id].vue` |
| 相册 | `features/album/` | types.ts + api.ts + index.ts | ✅ 已拆分 | `pages/home.vue`（生活 Tab） |
| 鉴权 | `features/auth/` | types.ts + api.ts + index.ts | ✅ 已拆分 | `pages/login.vue`、`stores/auth.ts`、`middleware/auth.ts` |
| 留言板 | `features/guestbook/` | types.ts + api.ts + index.ts | ✅ 已拆分 | `pages/home.vue`（留言板 Tab） |
| 更新日志 | `features/changelog/` | types.ts + index.ts | ✅ 已拆分 | `pages/home.vue`（更新日志 Tab） |
| 管理后台 | `features/admin/` | （预留） | 🚧 待拆分 | `pages/admin.vue`（57KB，后续拆分组件） |
| 酷炫首页 | — | — | ⏸️ 不拆分 | `pages/index.vue`（独立页面，无需 feature 模块） |
| Agent Team | — | — | ⏸️ 不拆分 | `pages/home.vue`（Agent Team Tab，数据写死在前端） |

## 共享层注册表

### Composables（`shared/composables/` → `composables/` re-export）

| 模块 | 实际路径 | Nuxt 自动导入路径 | 说明 |
|------|---------|------------------|------|
| 主题切换 | `shared/composables/useTheme.ts` | `composables/useTheme.ts` | Dark/Light 切换，Cookie 持久化，SSR 零闪烁 |
| 设备检测 | `shared/composables/useDevice.ts` | `composables/useDevice.ts` | 响应式 `isMobile` ref + 命令式 `isMobileDevice()` 函数 |

### Stores（`stores/`）

| 模块 | 文件 | 说明 |
|------|------|------|
| 认证状态 | `stores/auth.ts` | `isLoggedIn`、`username`、`checkAuth()`（SSR Cookie 转发）、`logout()` |

### Utils（`shared/utils/` → `utils/` re-export）

| 模块 | 实际路径 | Nuxt 自动导入路径 | 说明 |
|------|---------|------------------|------|
| 通用 API | `shared/utils/api.ts` | `utils/api.ts` | profile、upload 等跨领域 API |
| 分片上传 | `shared/utils/chunkedUpload.ts` | `utils/chunkedUpload.ts` | ≤1.5MB 直传，>1.5MB 分片上传 |
| CDN URL | `shared/utils/imageUrl.ts` | `utils/imageUrl.ts` | `toCdnUrl()` |
| 统一 UI | `shared/utils/ui.ts` | `utils/ui.ts` | `showSuccess/Error/Info/Confirm` |
| 移动端 UI | `shared/utils/mobileUI.ts` | `utils/mobileUI.ts` | `MobileToast/MobileDialog` |

### Components（`shared/components/` → `components/` re-export）

| 组件 | 实际路径 | Nuxt 自动导入路径 | 说明 |
|------|---------|------------------|------|
| AppLoading | `shared/components/AppLoading.vue` | `components/AppLoading.vue` | 跨端 Loading 组件 |

### Types（`shared/types.ts` → `types/index.ts` re-export）

| 文件 | 说明 |
|------|------|
| `shared/types.ts` | 通用类型：NavItem、TabItem、Profile |
| `types/index.ts` | 向后兼容聚合导出（re-export 所有 feature 类型 + shared 类型） |

### Plugins（`plugins/`）

| 插件 | 文件 | 运行环境 | 说明 |
|------|------|---------|------|
| Ant Design Vue | `plugins/antd.client.ts` | 仅客户端 | 注册 message/modal/spin 样式，配置 message 全局参数 |
| NProgress | `plugins/nprogress.client.ts` | 仅客户端 | 路由导航进度条（beforeEach/afterEach/onError） |
| Pinia Hydration Fix | `plugins/pinia-hydration-fix.ts` | 仅服务端 | JSON 往返清洗 payload，修复 SSR hydration 原型链问题 |
| Virtual Scroller | `plugins/vue-virtual-scroller.client.ts` | 仅客户端 | 注册 vue-virtual-scroller 全局组件 |
| PV Tracker | `plugins/pv-tracker.client.ts` | 仅客户端 | 路由变化时自动上报 PV/UV（fire and forget） |

## 页面路由映射

| 页面 | 路由 | 文件 | 依赖模块 | 说明 |
|------|------|------|----------|------|
| 酷炫首页 | `/` | `pages/index.vue`（17KB） | @chenglou/pretext | Canvas 代码滚动背景 + 拖拽按钮 |
| 博客主页 | `/home` | `pages/home.vue`（82KB） | useTheme, useDevice, api, types, vue-virtual-scroller, antd | 6 个 Tab：文章/生活/小工具·小游戏/Agent Team/更新日志/留言板 |
| 文章详情 | `/articles/:id` | `pages/articles/[id].vue`（8KB） | useTheme, api, Tiptap | Tiptap 只读渲染 + 点赞 |
| 登录 | `/login` | `pages/login.vue`（4KB） | auth store, api | 账号密码登录，成功跳转 redirect 或 /home |
| 管理后台 | `/admin` | `pages/admin.vue`（57KB） | auth store, api, Tiptap, antd, chunkedUpload | 6 个管理 Tab：文章/相册/个人资料/留言管理/数据统计 |

## 后端 API（透传代理 → Go 后端 `http://127.0.0.1:8080`）

> 所有 `/api/*` 路由均为 Nitro 薄代理层，通过 `proxyToBackend()` / `proxyRawToBackend()` 转发到 Go 后端，自动透传 Cookie 和客户端 IP。

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
| GET | `/api/theme` | ❌ | 获取主题偏好 |
| POST | `/api/theme` | ❌ | 保存主题偏好 |

### PV/UV 数据统计

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| POST | `/api/pv/record` | ❌ | 上报访问记录（path/device_id/referer） |
| GET | `/api/pv/overview` | ✅ | 今日/总 PV/UV 概览数据 |
| GET | `/api/pv/trend` | ✅ | PV/UV 趋势数据（`?days=7\|30`） |
| GET | `/api/pv/top-pages` | ✅ | Top5 页面报表（`?days=7\|30`） |
| GET | `/api/pv/logs` | ✅ | 访问日志分页查询（`?page=&page_size=&path=&date=`） |

## 服务端代理层（`server/utils/`）

> 前端 Nitro Server 层为**纯透传代理**，不直连任何数据库或缓存，所有请求转发到 Go 后端 `http://127.0.0.1:8080`。

| 文件 | 说明 |
|------|------|
| `proxy.ts` | 通用代理工具：`proxyToBackend()`（JSON 请求透传）+ `proxyRawToBackend()`（multipart 原始透传，用于文件上传）。自动透传 Cookie / Set-Cookie / X-Forwarded-For 等 header |

## 环境配置

前端通过 `.env` 文件注入环境变量（Nuxt runtimeConfig 读取）：

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `BACKEND_URL` | Go 后端 API 地址 | `http://127.0.0.1:8080` |

## 版本号规范

格式：`X.Y.Z`（Semantic Versioning）

| 位置 | 含义 | 何时递增 |
|------|------|---------|
| X（主版本） | 重大重构或破坏性变更 | 架构级改动 |
| Y（次版本） | 常规功能迭代 | **默认每次发版递增此位** |
| Z（补丁版本） | 紧急 bug 修复 | 仅修复线上问题 |

## 更新日志规范

- 通过 Go 后端 API 管理
- 每条最多 **20 字**，最多 **5 条**
- 内容风格：Emoji 前缀 + 简短描述，如 `"🎉 个人博客正式上线"`

## 变更日志（最近重要变更）

- 2026-04-06: 新增数据统计功能（PV/UV 折线图、Top5 页面报表、访问日志查询）+ 前台 PV 自动埋点插件
- 2026-04-05: 前端 Server Routes 全部改为透传代理，调用 Go 后端 API，移除 better-sqlite3/ioredis/cos-nodejs-sdk/bcryptjs/sharp/uuid 依赖
- 2026-04-05: 上传接口迁移至腾讯云 COS，toCdnUrl 适配新 URL 格式
- 2026-04-04: Feature-First 目录结构重构（types/api 按领域拆分到 features/，通用层移至 shared/）
- 2026-04-04: 默认 Dark 模式 + 修正 GitHub 链接（fatwillzeng → fatwill）
- 2026-04-04: SSR 阶段 Cookie 转发修复重启后需重新登录
- 2026-04-04: 新增「Agent Team」Tab，展示 4 个 AI Agent 卡片
- 2026-04-04: Tab 导航添加 Emoji 图标前缀 + 美化样式
- 2026-04-03: 首页重写为 Canvas + @chenglou/pretext 代码滚动背景 + 拖拽按钮
- 2026-04-02: 文章点赞 + 留言板 Tab 功能（v1.9.0）
- 2026-04-02: 相册/照片密码保护功能（v1.8.0）
- 2026-04-01: CDN 图片 URL 转换、ant-design-vue 集成、虚拟滚动优化
- 2026-04-01: 首屏 CSR→SSR 改造（useAsyncData 预取 + Pinia hydration 修复）
- 2026-03-31: 相册功能、分片上传、Token 持久化
- 2026-03-30: 全面重写博客前端（Dark/Light 主题、Tab 导航、Tiptap 编辑器）
