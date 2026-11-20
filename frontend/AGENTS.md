# AGENTS.md - 前端项目中枢索引

> 仓库级运维、EdgeOne/CDN 缓存架构与总规范见根目录 [AGENTS.md](../AGENTS.md)。

## ⚠️ 文档同步与提交铁律（必守）

- **同步更新**：任何前端代码改动（新增/删改页面路由、Feature 模块、组件、Composable、Plugin、Middleware、Nitro API 代理、依赖或构建配置），**必须在同一次改动中同步更新本文件，禁止长期滞后**。与原内容冲突时直接修改原条目，并在底部「变更日志」追加记录。
- **改完直接提交**：改动并验证完成后**直接 `git commit && git push`，不用再问**（只暂存本次改动文件，不误提交无关变更）。

---

## 项目概述

fatwill 个人博客前端，基于 Nuxt 3 SSR + Go 后端 API 架构。生产环境由 Nginx 直接将 `/api/*` 转发给 Go 后端（`127.0.0.1:8080`），少数需字段映射或流式透传的接口（微信管理、图一乐 SSE 聊天）以及 SSR 阶段请求经由 Nitro Server 层透传代理。支持 Dark/Light 主题切换、文章、生活相册、成长日记、装修管理、玩具（图一乐 / SBTI / 太空游戏）、留言板与管理后台等功能。

## 技术栈速览

| 类别 | 技术 | 版本 |
|------|------|------|
| 框架 | Nuxt 3 | ^3.16.1 |
| 视图层 | Vue 3 | ^3.5.13 |
| 语言 | TypeScript（strict） | ^5.8.2 |
| 样式 | Tailwind CSS（`@nuxtjs/tailwindcss`，darkMode: class） | ^6.13.2 |
| 状态管理 | Pinia（`@pinia/nuxt`） | ^2.3.1 / ^0.9.0 |
| 富文本编辑器 | Tiptap（vue-3 + starter-kit + image + link + table + html） | ^3.21.0 / ^3.22.3 |
| Markdown / 高亮 | `markdown-it` + `highlight.js` | ^14.1.1 / ^11.11.1 |
| UI 组件库 | Ant Design Vue（PC 端按需/客户端加载） | ^4.2.6 |
| 图表 / PDF / 二维码 | `echarts`、`pdfjs-dist`、`qrcode` | ^6.0.0 / ^5.6.205 / ^1.5.4 |
| 微信分享 | `weixin-js-sdk` | ^1.6.5 |
| Canvas 文字排版 | `@chenglou/pretext` | ^0.0.4 |
| 日期 / 进度条 | `dayjs`、`nprogress` | ^1.11.23 / ^0.2.0 |
| 后端 API | Go 后端（`http://127.0.0.1:8080`） | — |
| 包管理 | npm | — |

## 最高优先级读取

1. **本文件** `frontend/AGENTS.md` 与根 [AGENTS.md](../AGENTS.md)
2. `docs/eslint-rules.md` — 代码规范与命名约定
3. `docs/architecture.md` — 整体架构概览与目录结构说明

## 目录结构（Feature-First + Nuxt 约定混合架构）

> 类型定义和 API 服务层按领域拆分在 `features/`，跨领域通用组件/composables/工具位于 `shared/`；Nuxt 约定目录（`components/`、`composables/`、`utils/`、`types/`）保留向后兼容 re-export 或直接存放通用模块，确保 Nuxt 自动导入正常工作。
> 页面级 `.vue` 文件（`home.vue`、`admin.vue`）体量较大，后续按需继续拆分子组件。

```
frontend/
├── AGENTS.md                          # 前端中枢索引（本文件）
├── app.vue                            # 应用入口（主题 class、运行时 canonical/preconnect 注入）
├── nuxt.config.ts                     # Nuxt 配置（runtimeConfig、Tab 独立路由、nuxtLink 交互预取）
├── tailwind.config.ts                 # Tailwind 配置（dark mode: class，扫描 features/**）
├── package.json / tsconfig.json
├── .env / .env.example / .env.production
├── assets/styles/tailwind.css         # Tailwind 入口样式
├── docs/                              # 工程文档体系（architecture / eslint-rules / features / pages / shared）
├── features/                          # 🏗️ Feature-First 领域模块
│   ├── article/                       #   文章领域（types.ts / api.ts / index.ts）
│   ├── album/                         #   相册领域（types.ts / api.ts / index.ts）
│   ├── auth/                          #   鉴权领域（types.ts / api.ts / index.ts）
│   ├── guestbook/                     #   留言板领域（types.ts / api.ts / index.ts）
│   ├── changelog/                     #   更新日志领域（types.ts / index.ts）
│   ├── renovation/                    #   装修知识领域（types.ts / api.ts / index.ts，已对接后端 API）
│   ├── material/                      #   装修材料清单领域（types.ts / api.ts / index.ts，已对接后端 API）
│   ├── tuyile/                        #   图一乐领域（types/api/utils/styles + 10 组件 + 3 composables）
│   └── admin/                         #   管理后台领域（预留）
├── shared/                            # 🔗 跨领域共享层
│   ├── types.ts                       #   通用类型（NavItem、TabItem、Profile）
│   ├── components/AppLoading.vue      #   跨端 Loading 组件
│   ├── composables/                   #   useTheme.ts / useDevice.ts
│   └── utils/                         #   api.ts / ui.ts / mobileUI.ts / imageUrl.ts / chunkedUpload.ts
├── components/                        # Nuxt 自动导入组件
│   ├── AppLoading.vue                 #   跨端 Loading 组件
│   ├── MediaViewer.vue                #   图片/视频媒体灯箱预览组件
│   └── SearchDialog.vue               #   Cmd+K 全站搜索弹窗
├── composables/                       # Nuxt 自动导入 Composables
│   ├── useAuth.ts                     #   鉴权状态与登录/退出封装
│   ├── useDevice.ts                   #   → shared/composables/useDevice.ts
│   ├── useSiteUrl.ts                  #   域名配置统一入口（siteUrl/assetsUrl/cdnUrl 等）
│   ├── useTheme.ts                    #   → shared/composables/useTheme.ts
│   └── useWechatShare.ts              #   微信 JS-SDK 分享初始化（支持 signUrl 重签兜底）
├── stores/auth.ts                     # Pinia 认证状态管理
├── middleware/
│   ├── auth.ts                        # 鉴权路由守卫（保护 /admin）
│   └── tab-redirect.global.ts         # 全局重定向（/home 301 → /articles，兼容 ?tab=xxx）
├── layouts/default.vue                # 默认布局
├── pages/                             # 页面路由（详见下方「页面路由映射」）
├── plugins/                           # 客户端/服务端插件（详见下方「Plugins」）
├── public/                            # 静态资源（图标、china.json、test/* 静态测试页等）
├── scripts/migrate-image-urls.ts      # 一次性历史数据迁移脚本
└── server/                            # Nuxt Server（Nitro 代理与 SEO 路由）
    ├── api/                           #   透传代理路由（auth/articles/albums/photos/messages/pv/profile/
    │                                  #   upload/theme/changelog/growth-diary/renovation/toys/tuyile/wechat/admin）
    ├── routes/                        #   动态 SEO 路由（robots.txt.ts / sitemap.xml.ts）
    ├── utils/                         #   proxy.ts（通用代理）/ tuyileProxy.ts（SSE 流式透传）
    └── data/                          #   articles.json / changelog.json（底稿数据）
```

## 领域模块注册表

| 领域 | Feature 路径 | 包含内容 | 状态 | 页面级 UI |
|------|-------------|---------|------|----------|
| 文章 | `features/article/` | `types.ts` + `api.ts` + `index.ts` | ✅ 已对接 | `pages/home.vue`（文章 Tab）、`pages/articles/[id].vue` |
| 相册 | `features/album/` | `types.ts` + `api.ts` + `index.ts` | ✅ 已对接 | `pages/home.vue`（生活 Tab） |
| 鉴权 | `features/auth/` | `types.ts` + `api.ts` + `index.ts` | ✅ 已对接 | `pages/login.vue`、`stores/auth.ts`、`composables/useAuth.ts` |
| 留言板 | `features/guestbook/` | `types.ts` + `api.ts` + `index.ts` | ✅ 已对接 | `pages/home.vue`（留言板 Tab） |
| 更新日志 | `features/changelog/` | `types.ts` + `index.ts` | ✅ 已对接 | `pages/home.vue`（更新日志 Tab） |
| 装修 | `features/renovation/` | `types.ts` + `api.ts` + `index.ts` | ✅ 已对接后端 `/api/renovation/articles` | `pages/home.vue`（装修 Tab）、`pages/renovation/*` |
| 材料清单 | `features/material/` | `types.ts` + `api.ts` + `index.ts` | ✅ 已对接后端 `/api/materials` | `pages/admin.vue` / 材料附件联动 `pages/pdf/[id].vue` |
| 图一乐 | `features/tuyile/` | `types.ts` + `api.ts` + `utils.ts` + `styles/` + `composables/` + `components/` + `index.ts` | ✅ 已上线（二期） | `pages/toys/tuyile/index.vue`、`pages/toys/tuyile/share/[hash].vue` |
| 管理后台 | `features/admin/` | （预留） | 🚧 待拆分 | `pages/admin.vue` |
| 酷炫首页 / Agent Team | — | — | ⏸️ 独立实现 | `pages/index.vue`、`pages/home.vue` |

## 共享层注册表

### Composables

| 模块 | 路径 | 说明 |
|------|------|------|
| 鉴权封装 | `composables/useAuth.ts` | 登录态访问、SSR `useState` 联动 |
| 主题切换 | `shared/composables/useTheme.ts`（`composables/useTheme.ts` re-export） | Dark/Light 切换，Cookie 持久化，SSR 零闪烁 |
| 设备检测 | `shared/composables/useDevice.ts`（`composables/useDevice.ts` re-export） | 响应式 `isMobile` ref + 命令式 `isMobileDevice()` |
| 站点 URL | `composables/useSiteUrl.ts` | 域名统一入口：`siteUrl`/`assetsUrl`/`cdnUrl` + `articleUrl()`/`pageUrl()`/`defaultOgImage()`/`assetUrl()` |
| 微信分享 | `composables/useWechatShare.ts` | 微信公众号 JS-SDK 签名与自定义卡片，支持可选 `signUrl` 重签兜底 |

### Components & Utils

| 类型 | 模块 | 路径 | 说明 |
|------|------|------|------|
| Component | `AppLoading` | `components/AppLoading.vue` / `shared/components/AppLoading.vue` | 跨端 Loading 组件 |
| Component | `MediaViewer` | `components/MediaViewer.vue` | 图片/视频全屏灯箱预览（含缩放、切换、视频播放） |
| Component | `SearchDialog` | `components/SearchDialog.vue` | Cmd+K 搜索对话框 |
| Store | `auth` | `stores/auth.ts` | Pinia 认证状态管理 |
| Util | 通用 API | `shared/utils/api.ts`（`utils/api.ts` re-export） | profile、upload 等跨领域 API |
| Util | 分片上传 | `shared/utils/chunkedUpload.ts` | ≤1.5MB 直传，>1.5MB 分片上传 |
| Util | CDN URL | `shared/utils/imageUrl.ts` | `toCdnUrl()` / `toThumbUrl()` / `toWebpUrl()`（启动由插件注入 assets host） |
| Util | 统一 UI | `shared/utils/ui.ts` + `shared/utils/mobileUI.ts` | `showSuccess/Error/Info/Confirm` + 移动端 Toast/Dialog |

### Plugins（`plugins/`）

| 插件 | 文件 | 运行环境 | 说明 |
|------|------|----------|------|
| ImageUrl Host | `0.image-url-host.ts` | 服务端 + 客户端 | 最早期把 `runtimeConfig.public.assetsUrl` 注入 `shared/utils/imageUrl.ts` |
| Sandbox Safe | `0.sandbox-safe.client.ts` | 仅客户端（`enforce: 'pre'`） | 页面被嵌入无 `allow-same-origin` 的 sandbox iframe 时，兜底 `cookie`/`localStorage`/`sessionStorage` 防崩 |
| Auth Check | `auth-check.server.ts` | 仅服务端 | SSR 阶段直连 Go 后端 `/api/auth/check` 确定登录态写入 `useState`，避免 Nitro 内部循环调用 |
| Ant Design Vue | `antd.client.ts` | 仅客户端 | 按需/客户端加载 antd message/modal 等能力 |
| NProgress | `nprogress.client.ts` | 仅客户端 | 路由导航进度条 |
| Perf Tracker | `perf-tracker.client.ts` | 仅客户端 | 采集 FCP/LCP/TTFB/DOM Ready/Load/JS Load 上报 `/api/perf/report` |
| Pinia Hydration Fix | `pinia-hydration-fix.ts` | 仅服务端 | JSON 往返清洗 payload，修复 SSR hydration 原型链问题 |
| PV Tracker | `pv-tracker.client.ts` | 仅客户端 | 路由变化自动静默上报 `/api/pv/record` |
| 图一乐入口 URL | `tuyile-entry-url.client.ts` | 仅客户端 | 记录首次进入 URL，供 iOS 微信 JS-SDK 重签兜底 |

> 注：`nuxt.config.ts` 已开启 `experimental.defaults.nuxtLink.prefetchOn: { visibility: false, interaction: true }`，避免列表页首屏预取详情/后台大 chunk 抢占带宽；虚拟滚动插件已移除。

## 页面路由映射

| 页面 | 路由 | 文件 | 说明 |
|------|------|------|------|
| 酷炫首页 | `/` | `pages/index.vue` | Canvas + `@chenglou/pretext` 代码滚动背景 + 拖拽按钮 |
| 📝 文章 Tab | `/articles` | `pages/home.vue`（`pages:extend` 注册） | 默认 Tab，文章列表 |
| 📷 生活 Tab | `/life` | `pages/home.vue`（`pages:extend` 注册） | 生活相册与成长日记入口 |
| 🏠 装修 Tab | `/renovation` | `pages/home.vue`（`pages:extend` 注册） | 装修功能入口与知识文章列表 |
| 🎮 玩具 Tab | `/toys` | `pages/home.vue`（`pages:extend` 注册） | 玩具集合 Tab（注意 `nuxt.config.ts` 中为 `/toys`） |
| 🤖 Agent Team | `/agent-team` | `pages/home.vue`（`pages:extend` 注册） | Agent Team 介绍 Tab |
| 💬 留言板 | `/guestbook` | `pages/home.vue`（`pages:extend` 注册） | 留言板 Tab |
| 📋 更新日志 | `/changelog` | `pages/home.vue`（`pages:extend` 注册） | 更新日志 Tab |
| 旧主页重定向 | `/home` | `middleware/tab-redirect.global.ts` | 301 重定向到 `/articles`，兼容 `?tab=xxx` |
| 文章详情 | `/articles/:id` | `pages/articles/[id].vue` | Tiptap 只读渲染、左侧 TOC、点赞、阅读量、图片事件委托预览与 WebP 输出 |
| 成长日记 | `/life/growth-diary` | `pages/life/growth-diary.vue` | 阳阳的成长日记（按月浏览，登录支持批量编辑、图片/视频） |
| 装修预算 | `/renovation/budget` | `pages/renovation/budget.vue` | 装修成本预算表（预算/实际支出、5 大分类、支持登录编辑） |
| 装修进度 | `/renovation/progress` | `pages/renovation/progress.vue` | 装修进度与踩坑日记页面 |
| PDF 预览 | `/pdf/:id` | `pages/pdf/[id].vue` | 纯客户端（`ssr: false`）`pdfjs-dist` Canvas 渲染附件 PDF |
| 太空游戏 | `/toys/starfling` | `pages/toys/starfling.vue` | 全屏 iframe 嵌入 |
| SBTI 测试 | `/toys/sbti` | `pages/toys/sbti.vue` | SBTI 人格测试 |
| 图一乐 | `/toys/tuyile` | `pages/toys/tuyile/index.vue` | 老黄历 + AG-UI 运势聊天（`view=almanac\|chat`、`date=`），文档 `docs/pages/tuyile.md` |
| 图一乐分享 | `/toys/tuyile/share/:hash` | `pages/toys/tuyile/share/[hash].vue` | 对话分享只读页（SSR、OG、微信卡片、noindex），文档 `docs/pages/tuyile-share.md` |
| 登录 | `/login` | `pages/login.vue` | 博主账号密码登录 |
| 管理后台 | `/admin` | `pages/admin.vue` | 鉴权保护：内容管理、相册、材料清单、PV/UV 与性能统计、微信同步 |
| 问卷测试页 | `/test/wj-*` | `pages/test/wj-{html,html2,jspsych,jspsych2,jspsych3}-test.vue` | 问卷 / jsPsych 联调测试页（对应 `public/test/*`） |

## Nitro Server 代理层（`server/`）

> 生产 Nginx 中 `/api/*` 默认直连 Go 后端 `127.0.0.1:8080`；其中 `/api/admin/wechat/*`、`/api/admin/articles/:id/wechat-sync-logs`、`/api/toys/tuyile/chat`（SSE 禁止缓冲）以及 SSR/本地开发请求经 Nitro `server/api/*` 转发。完整后端 API 清单以 [backend/AGENTS.md](../backend/AGENTS.md) 为准。

- `server/utils/proxy.ts`：`proxyToBackend()`（JSON 透传）+ `proxyRawToBackend()`（multipart 上传原始透传），自动透传 Cookie / Set-Cookie / X-Forwarded-For。
- `server/utils/tuyileProxy.ts`：`proxyTuyileJson()` + `proxyTuyileStream()`（SSE 流式逐块透传，客户端断开即中止上游）。
- Nitro 已覆盖代理组：`auth/*`、`articles/*`（含 `random`、`view`、`like-status-batch`）、`albums/*`、`photos/*`、`messages/*`、`profile/*`、`upload/*`、`theme`、`changelog`、`pv/*`（含 `geo`）、`growth-diary/*`、`renovation/{budget,pitfall}/*`、`toys/tuyile/*`、`wechat/jssdk-config`、`admin/wechat/*`、`admin/articles/:id/*`。

## 环境配置与域名体系

通过 `.env` / `.env.production` 注入（Nuxt `runtimeConfig`）：

| 变量 | 作用 | 默认值 | 生效方式 |
|------|------|--------|---------|
| `BACKEND_URL` | Go 后端内网地址 | `http://127.0.0.1:8080` | 运行时 |
| `NUXT_PUBLIC_SITE_URL` | 站点主域名（canonical / OG / sitemap / robots / 分享） | `https://fatwill.cn` | 运行时 |
| `NUXT_PUBLIC_ASSETS_URL` | 图片/素材 CDN（EO 直连 COS） | `https://assets.fatwill.cn` | 运行时 |
| `NUXT_PUBLIC_CDN_URL` | Nuxt 构建产物 CDN（`app.cdnURL`） | `https://cdn.fatwill.cn` | ⚠️ **构建期（换域须重 build）** |

- 统一读取入口：`composables/useSiteUrl.ts`（`siteUrl`、`assetsUrl`、`cdnUrl`、`articleUrl()`、`pageUrl()`、`defaultOgImage()`、`assetUrl()`）。
- 线上加速与两层缓存（EdgeOne `blog-cache` + VPS Nginx 24h SSR 缓存）详见根 [AGENTS.md](../AGENTS.md)。

## SEO 与动态路由（`server/routes/`）

- **`/robots.txt`**（`server/routes/robots.txt.ts`）：默认拒绝，显式 Allow `/$`、`/articles`、`/life$`、`/tools$`、`/agent-team$`、`/guestbook$`、`/changelog$`；Disallow `/admin`、`/login`、`/api/` 及兜底 `/`；Sitemap 动态读取 `siteUrl`。*(注：玩具 Tab 路由在 `nuxt.config.ts` 中为 `/toys`，若需放开玩具页收录记得同步更新此处。)*
- **`/sitemap.xml`**（`server/routes/sitemap.xml.ts`）：动态拉取后端全部文章，输出 `/articles`（priority 1.0）及全部 `/articles/:id`（priority 0.8），缓存 `max-age=3600`。

## 版本号与更新日志规范

- 版本号 `X.Y.Z`（Semantic Versioning）：默认常规功能迭代递增次版本 `Y`，紧急修复递增补丁位 `Z`。
- 更新日志通过 Go 后端 `pkg/db/changelog_seed.go` 幂等播种，每条 ≤ 20 字，最多 5 条，风格为 Emoji 前缀 + 简短描述。

## 变更日志（最近重要变更）

- 2026-11-20：全面梳理并同步 `frontend/AGENTS.md`（补齐 `features/material`、`MediaViewer`、`useAuth`、`useWechatShare`、`0.sandbox-safe.client`、`auth-check.server`、`perf-tracker.client`、`/life/growth-diary`、`/pdf/:id`、`/test/wj-*`、`/toys` 路由与 `nuxtLink` 交互预取说明，移除已不存在页面记录）；新增 `AGENTS.md` 强制随改动同步更新铁律；关联根 `AGENTS.md` 的 EdgeOne 加速与缓存架构。
- 2026-10-13：图一乐二期：顶部运势 banner（`FortuneBanner`）、年月日选择（antd DatePicker 异步 chunk）、思考状态条（`CUSTOM tuyile.status` / `REASONING_*`）、轻松风文案与对话分享页 `/toys/tuyile/share/:hash`；新增 `tuyile-entry-url.client.ts`；`dayjs` 升级为直接依赖。
- 2026-10-12：新增玩具「图一乐」（`/toys/tuyile`，`features/tuyile/`）及 Nitro `server/utils/tuyileProxy.ts` SSE 流式透传。
- 2026-09-30：修复文章图片点击预览失效，改为容器事件委托并在渲染阶段输出 WebP；底部备案信息追加粤公网安备链接与警徽图标。
- 2026-09-08：主域配置全面切换为 `fatwill.cn` / `assets.fatwill.cn` / `cdn.fatwill.cn`（新增 `.env.production`）。
- 2026-08-31：域名配置化改造（`useSiteUrl.ts`、`0.image-url-host.ts`、动态 `robots.txt.ts`）。
