# AGENTS.md - 项目中枢索引

## 项目概述
fatwillzeng 个人博客，基于 Nuxt 3 + Vue 3 + TypeScript + Tailwind CSS 构建，支持 Dark/Light 主题切换。

## 技术栈速览
- **框架**: Nuxt 3.21.2 + Vue 3.5
- **语言**: TypeScript 5.8
- **样式**: Tailwind CSS（dark mode via `class` 策略）
- **状态管理**: Pinia
- **富文本**: Tiptap（@tiptap/vue-3 + starter-kit + image + link）
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

## 共享层注册表

| 模块 | 路径 | 说明 |
|------|------|------|
| 主题切换 | `composables/useTheme.ts` | Dark/Light 主题切换，持久化 localStorage |
| API 服务 | `utils/api.ts` | 封装登录、文章 CRUD 接口 |
| 认证 Store | `stores/auth.ts` | Pinia 认证状态管理 |
| 类型定义 | `types/index.ts` | 全局 TS 类型：文章、登录、Tab 等 |

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
| 文章列表 | GET | `/api/articles` | 返回列表（不含 content） |
| 文章详情 | GET | `/api/articles/:id` | 返回详情（含 Tiptap JSON content） |
| 创建文章 | POST | `/api/articles` | 需鉴权（cookie auth_token） |
| 更新文章 | PUT | `/api/articles/:id` | 需鉴权，部分更新，自动刷新 updatedAt |
| 更新日志 | GET | `/api/changelog` | 返回所有版本更新日志，数据源 server/data/changelog.json |

## 版本号规范（用户明确要求）

版本号格式：`X.Y.Z`（遵循 Semantic Versioning）

| 位置 | 含义 | 何时递增 |
|------|------|---------|
| X（主版本） | 重大重构或破坏性变更 | 特殊需求，架构级改动 |
| Y（次版本） | 常规功能迭代 | **默认每次发版递增此位** |
| Z（补丁版本） | 紧急 bug 修复 | 仅修复线上问题时使用 |

**默认递增规则**：非特殊需求，每次发版只递增中间位（Y），如 `1.0.0 → 1.1.0 → 1.2.0`。

## 更新日志规范（用户明确要求）

- 每次前端发版必须在 `/server/data/changelog.json` 中新增一条版本记录
- 字段说明：
  - `version`：版本号字符串，如 `"1.1.0"`
  - `date`：发版日期，格式 `"YYYY-MM-DD"`
  - `logs`：更新内容数组，**每条最多 20 字，最多 5 条**
- 不需要登录即可查看的更新内容都应记录；涉及敏感功能（如管理后台权限逻辑）的条目可省略
- 内容风格：用 Emoji 前缀 + 简短描述，如 `"🎉 个人博客正式上线"`

## 变更日志
- 2026-03-31: 为所有数据结构统一添加 createdAt/updatedAt 时间戳字段（ChangelogItem 接口 + changelog.json 历史数据补充）
- 2026-03-31: 新增「更新日志」Tab，时间轴风格展示版本历史；新增版本号规范和更新日志规则到 AGENTS.md
- 2026-03-31: 新增更新文章接口（PUT /api/articles/:id），支持 updatedAt 自动刷新；新增 UpdateArticleRequest 类型和 apiUpdateArticle 前端调用
- 2026-03-30: 全面重写首页，实现 Dark/Light 主题切换、个人信息区域、Tab 导航、文章列表
- 2026-03-30: 新增文章详情页（Tiptap 渲染）、登录页、管理后台（Tiptap 编辑器）
- 2026-03-30: 新增 useTheme composable、auth store、API 服务层
