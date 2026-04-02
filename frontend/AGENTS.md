# AGENTS.md - 项目中枢索引

## 项目概述
fatwillzeng 个人博客，基于 Nuxt 3 + Vue 3 + TypeScript + Tailwind CSS 构建，支持 Dark/Light 主题切换。

## 技术栈速览
- **框架**: Nuxt 3.21.2 + Vue 3.5
- **语言**: TypeScript 5.8
- **样式**: Tailwind CSS（dark mode via `class` 策略）
- **状态管理**: Pinia
- **富文本**: Tiptap（@tiptap/vue-3 + starter-kit + image + link）
- **数据库**: MySQL 8.0（mysql2 驱动，连接池模式）
- **包管理**: npm

## 最高优先级读取
1. `docs/eslint-rules.md` - ESLint/代码规范（最高优先级）
2. `docs/architecture.md` - 整体架构说明

## 领域模块注册表

| 模块 | 路径 | 状态 | 说明 |
|------|------|------|------|
| 首页 | `pages/index.vue` | ✅ 已实现 | 个人信息、Tab 导航（文章/生活/小工具·小游戏/agent team/更新日志）、文章列表 |
| 文章详情 | `pages/articles/[id].vue` | ✅ 已实现 | Tiptap 只读渲染文章内容 |
| 登录 | `pages/login.vue` | ✅ 已实现 | 账号密码登录，成功后跳转 /admin |
| 管理后台 | `pages/admin.vue` | ✅ 已实现 | Tiptap 编辑器，录入并发布文章 |
| 相册 | `server/api/albums/`, `server/api/photos/` | ✅ 已实现 | 相册集 CRUD、照片管理、封面自动更新 |
| 文章点赞 | `server/api/articles/[id]/like*.ts` | ✅ 已实现 | 文章点赞/取消、状态查询 |
| 留言板 | `server/api/messages/` | ✅ 已实现 | 留言 CRUD、IP 限频、每日修改限制 |

## 共享层注册表

| 模块 | 路径 | 说明 |
|------|------|------|
| 主题切换 | `composables/useTheme.ts` | Dark/Light 主题切换，持久化 localStorage |
| API 服务 | `utils/api.ts` | 封装登录、文章 CRUD、图片上传、个人资料接口 |
| 认证 Store | `stores/auth.ts` | Pinia 认证状态管理 |
| 类型定义 | `types/index.ts` | 全局 TS 类型：文章、登录、Tab、Profile 等 |
| 设备检测 | `composables/useDevice.ts` | 移动端/PC端检测（响应式 + 命令式） |
| 统一 UI | `utils/ui.ts` | 跨端 Loading/Toast/Dialog 统一调用（PC=antd，移动端=自定义） |
| 移动端 UI | `utils/mobileUI.ts` | 移动端 Toast/Dialog 命令式组件（仿 antd-mobile 风格） |
| Loading 组件 | `components/AppLoading.vue` | 跨端 Loading 组件（PC=ASpin，移动端=自定义 spinner） |
| 分片上传 | `utils/chunkedUpload.ts` | 分片上传工具（进度回调、断点续传、自动清理） |
| CDN URL 转换 | `utils/imageUrl.ts` | 将 /uploads/ 路径转为 CDN 完整 URL（toCdnUrl） |
| 数据库连接 | `server/utils/db.ts` | MySQL 连接池（mysql2，进程级单例） |
| 文章 DAO | `server/utils/articles.ts` | 文章 CRUD 操作（MySQL） |
| 更新日志 DAO | `server/utils/changelog.ts` | 更新日志查询（MySQL） |
| 个人资料 DAO | `server/utils/profile.ts` | 博主资料读写（MySQL profile 表） |
| 相册 DAO | `server/utils/albums.ts` | 相册集和照片 CRUD（MySQL albums/photos 表） |
| 留言 DAO | `server/utils/messages.ts` | 留言 CRUD（MySQL messages 表） |
| IP 限频 | `server/utils/rateLimit.ts` | 内存 Map 实现的 IP 限频工具 |

## 页面路由映射

| 页面 | 路由 | 文件 | 依赖模块 |
|------|------|------|----------|
| 首页 | `/` | `pages/index.vue` | useTheme, api, types |
| 文章详情 | `/articles/:id` | `pages/articles/[id].vue` | useTheme, api, Tiptap |
| 登录 | `/login` | `pages/login.vue` | auth store, api |
| 管理后台 | `/admin` | `pages/admin.vue` | auth store, api, Tiptap |

## 后端 API（Nuxt Server）

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 登录 | POST | `/api/auth/login` | 账号密码验证，写入 httpOnly cookie |
| 文章列表 | GET | `/api/articles` | 返回列表（不含 content），支持 `?title=xxx` 关键词筛选 |
| 文章详情 | GET | `/api/articles/:id` | 返回详情（含 Tiptap JSON content、coverImage） |
| 创建文章 | POST | `/api/articles` | 需鉴权，支持 coverImage 封面图 |
| 更新文章 | PUT | `/api/articles/:id` | 需鉴权，部分更新，支持 coverImage |
| 删除文章 | DELETE | `/api/articles/:id` | 需鉴权，物理删除 |
| 更新日志 | GET | `/api/changelog` | 返回所有版本更新日志，数据源 MySQL changelogs 表 |
| 图片上传 | POST | `/api/upload` | 需鉴权，multipart/form-data，支持 jpg/png/gif/webp |
| 上传分片 | POST | `/api/upload/chunk` | 需鉴权，上传单个分片（uploadId/chunkIndex/totalChunks/filename） |
| 合并分片 | POST | `/api/upload/merge` | 需鉴权，合并所有分片，返回最终 URL |
| 取消上传 | DELETE | `/api/upload/chunk` | 需鉴权，清理临时分片文件 |
| 获取资料 | GET | `/api/profile` | 公开接口，返回博主头像和简介 |
| 更新资料 | PUT | `/api/profile` | 需鉴权，更新博主头像和简介 |
| 相册列表 | GET | `/api/albums` | 公开接口，返回所有相册集（含 photoCount、coverUrl） |
| 创建相册 | POST | `/api/albums` | 需鉴权，创建相册集 |
| 更新相册 | PUT | `/api/albums/:id` | 需鉴权，更新相册集名称和描述 |
| 删除相册 | DELETE | `/api/albums/:id` | 需鉴权，级联删除相册及所有照片 |
| 照片列表 | GET | `/api/albums/:id/photos` | 公开接口，获取相册下所有照片 |
| 添加照片 | POST | `/api/albums/:id/photos` | 需鉴权，添加照片并自动更新封面 |
| 删除照片 | DELETE | `/api/photos/:id` | 需鉴权，删除照片并自动更新封面 |
| 更新照片 | PUT | `/api/photos/:id` | 需鉴权，更新照片信息（caption、password） |
| 相册验证密码 | POST | `/api/albums/:id/verify-password` | 公开接口，验证相册集密码 |
| 照片验证密码 | POST | `/api/photos/:id/verify-password` | 公开接口，验证照片密码 |
| 文章点赞 | POST | `/api/articles/:id/like` | 公开接口，点赞/取消点赞（切换） |
| 文章点赞状态 | GET | `/api/articles/:id/like-status` | 公开接口，查询点赞状态 |
| 留言列表 | GET | `/api/messages` | 公开接口，返回所有留言 |
| 新增留言 | POST | `/api/messages` | 公开接口，新增/修改留言（IP 限频） |
| 修改留言 | PUT | `/api/messages/:id` | 公开接口，修改留言（每日一次） |

## 版本号规范（用户明确要求）

版本号格式：`X.Y.Z`（遵循 Semantic Versioning）

| 位置 | 含义 | 何时递增 |
|------|------|---------|
| X（主版本） | 重大重构或破坏性变更 | 特殊需求，架构级改动 |
| Y（次版本） | 常规功能迭代 | **默认每次发版递增此位** |
| Z（补丁版本） | 紧急 bug 修复 | 仅修复线上问题时使用 |

**默认递增规则**：非特殊需求，每次发版只递增中间位（Y），如 `1.0.0 → 1.1.0 → 1.2.0`。

## 更新日志规范（用户明确要求）

- 更新日志存储在 MySQL `changelogs` 表中（旧的 `server/data/changelog.json` 已废弃）
- 字段说明：
  - `version`：版本号字符串，如 `"1.1.0"`（UNIQUE 约束）
  - `date`：发版日期，格式 `"YYYY-MM-DD"`
  - `logs`：更新内容 JSON 数组，**每条最多 20 字，最多 5 条**
  - `created_at` / `updated_at`：自动维护的时间戳
- 不需要登录即可查看的更新内容都应记录；涉及敏感功能（如管理后台权限逻辑）的条目可省略
- 内容风格：用 Emoji 前缀 + 简短描述，如 `"🎉 个人博客正式上线"`

## 数据库表注册表

| 表名 | 说明 | 主键 | 索引 |
|------|------|------|------|
| `articles` | 文章表 | `id` (varchar(36) UUID) | `idx_created_at` |
| `changelogs` | 更新日志表 | `id` (bigint auto_increment) | `uk_version`, `idx_date` |
| `auth_tokens` | 登录 token 表 | `token` (varchar(64)) | `idx_username` |
| `profile` | 博主个人资料表 | `id` (int, 固定为1) | — |
| `albums` | 相册集表 | `id` (int auto_increment) | — |
| `photos` | 照片表 | `id` (int auto_increment) | `idx_album_id` |
| `article_likes` | 文章点赞表 | `id` (int auto_increment) | `uq_article_device`, `idx_article_id` |
| `messages` | 留言板表 | `id` (int auto_increment) | `uq_device_id` |

## 环境配置

数据库连接通过 `.env` 文件注入（Nuxt runtimeConfig 读取）：
- `DB_HOST` — MySQL 主机地址（服务器用 localhost，本地开发用远程 IP）
- `DB_PORT` — 端口号（默认 3306）
- `DB_USER` — 数据库用户名
- `DB_PASSWORD` — 数据库密码
- `DB_NAME` — 数据库名（blog）

## 变更日志
- 2026-04-02: 新增文章点赞功能（article_likes 表、like_count 冗余字段、点赞/取消切换接口、状态查询接口）；新增留言板功能（messages 表、CRUD 接口、IP 限频、每日修改限制、deviceId 脱敏）
- 2026-04-02: 前端密码保护功能完整实现（v1.8.0）：相册集/照片锁图标+模糊遮罩、密码验证弹窗（支持回车提交）、sessionStorage 缓存解锁状态、灯箱左右切换自动检查密码、后台创建/编辑相册密码、照片密码设置弹窗
- 2026-04-02: 相册和照片支持密码保护（bcrypt 哈希存储）；新增 verify-password 验证接口、照片 PUT 更新接口；查询接口返回 hasPassword 字段
- 2026-04-01: 修复 SSR 500 错误：用 ClientOnly 包裹 DynamicScroller（vue-virtual-scroller 不支持 SSR），文章列表 fallback 渲染静态卡片保证首屏 HTML 包含内容；三个首屏接口（articles/profile/changelog）已通过 useAsyncData 预取
- 2026-04-01: 新增 toCdnUrl 工具函数，所有图片 URL 统一走 CDN（cdn.fatwill.cloud）；覆盖首页头像/文章封面/相册/照片/灯箱及后台全部图片展示
- 2026-04-01: 统一 UI 组件跨端适配：PC 端 antd（Spin/Modal/message），移动端自定义轻量组件（MobileToast/MobileDialog）；封装 useDevice composable 和 AppLoading 组件；修复相册虚拟滚动日期标题被覆盖
- 2026-04-01: 引入 ant-design-vue@4.x；骨架屏替换为 Spin loading；alert/confirm 替换为 Modal.confirm + message；操作成功添加 message.success toast；修复相册虚拟滚动 header/grid 重叠
- 2026-04-01: 文章列表和相册照片列表添加虚拟滚动（vue-virtual-scroller DynamicScroller），大幅精简 DOM 节点提升滚动性能
- 2026-03-31: 图片上传支持进度条和分片上传（>2MB 自动分片）；移除前端 5MB 大小限制；新增 chunkedUpload 工具函数
- 2026-03-31: 新增相册功能（albums/photos 表）；完整 CRUD API（7个接口）；相册 DAO 层；前端类型定义和 API 封装
- 2026-03-31: 新增图片上传接口（POST /api/upload）和个人资料接口（GET/PUT /api/profile）；新建 profile 表和 profile DAO；uploads 目录静态伺服
- 2026-03-31: Token 持久化到 MySQL auth_tokens 表，重启/部署不再丢失登录态
- 2026-03-31: articles 表新增 cover_image 字段；文章接口全面支持 coverImage；新增 title 关键词筛选；新增 DELETE /api/articles/:id 删除接口
- 2026-03-31: 将文章和更新日志数据从 JSON 文件迁移到 MySQL 8.0；新增 mysql2 驱动、连接池工具（server/utils/db.ts）、articles/changelog DAO 层；改写全部 API 路由；新增 .env 环境配置
- 2026-03-31: Admin 入口图标 + Nuxt middleware 鉴权前置 + cookie 72h 滚动续期；新增 /api/auth/check、auth middleware、login redirect
- 2026-03-31: 为所有数据结构统一添加 createdAt/updatedAt 时间戳字段（ChangelogItem 接口 + changelog.json 历史数据补充）
- 2026-03-31: 新增「更新日志」Tab，时间轴风格展示版本历史；新增版本号规范和更新日志规则到 AGENTS.md
- 2026-03-31: 新增更新文章接口（PUT /api/articles/:id），支持 updatedAt 自动刷新；新增 UpdateArticleRequest 类型和 apiUpdateArticle 前端调用
- 2026-03-30: 全面重写首页，实现 Dark/Light 主题切换、个人信息区域、Tab 导航、文章列表
- 2026-03-30: 新增文章详情页（Tiptap 渲染）、登录页、管理后台（Tiptap 编辑器）
- 2026-03-30: 新增 useTheme composable、auth store、API 服务层
