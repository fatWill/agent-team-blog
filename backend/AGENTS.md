# AGENTS.md - 后端项目中枢索引

## 项目概述

fatwill 个人博客后端服务，基于 Go + Gin + GORM + SQLite + Redis 构建的 RESTful API，提供文章管理、相册管理、留言板、鉴权、文件上传等核心能力。服务部署在 2 核 2G 云服务器上，需注重性能优化。

## 技术栈速览

| 组件 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.26.1 |
| Web 框架 | Gin | v1.12.0 |
| ORM | GORM + glebarez/sqlite | v1.31.1 + v1.11.0 |
| 数据库 | SQLite (modernc.org/sqlite) | 纯 Go 实现 |
| 缓存 | Redis (go-redis) | v9.18.0 |
| 跨域 | gin-contrib/cors | v1.7.7 |
| 密码哈希 | golang.org/x/crypto/bcrypt | v0.49.0 |
| UUID | google/uuid | v1.6.0 |
| 包管理 | Go Modules | go.mod |

## 最高优先级读取

Agent 收到任务后，**必须按以下优先级依次读取**：

1. **本文件** `AGENTS.md` — 项目全貌、模块注册、路由总表
2. `docs/coding-conventions.md` — Go 编码规范（最高优先级）
3. `docs/sqlite.md` — SQLite 使用规范 + 全量表结构
4. `docs/http-api.md` — HTTP API 设计规范 + 错误码体系
5. `docs/architecture.md` — 整体架构概览

## 当前目录结构

```
backend/
├── AGENTS.md              # 🔑 中枢索引（最高优先级）
├── README.md              # 项目说明
├── main.go                # 入口：配置加载、依赖初始化、路由注册
├── go.mod / go.sum        # Go Modules 依赖管理
├── config/
│   └── config.go          # 配置结构体 + 环境变量加载
├── models/                # 📦 数据模型层（GORM struct + DTO）
│   ├── article.go         # Article / ArticleListItem / ArticleLike
│   ├── album.go           # Album / Photo / PhotoLike / PhotoDislike + ListItem
│   ├── misc.go            # Profile / Message / Changelog + ListItem
│   └── json.go            # 自定义 JSON 类型（MySQL JSON 字段）
├── internal/              # 🏗️ 业务核心（按领域分组）
│   ├── article/
│   │   └── handler.go     # 文章 CRUD + 点赞
│   ├── album/
│   │   └── handler.go     # 相册 CRUD + 密码验证 + UpdateAlbumCover
│   ├── photo/
│   │   └── handler.go     # 照片 CRUD + 点赞/踩 + 密码验证
│   ├── auth/
│   │   └── handler.go     # 登录、退出、鉴权检查
│   ├── upload/
│   │   └── handler.go     # 图片上传（直传 + 分片）+ RandomString
│   ├── guestbook/
│   │   └── handler.go     # 留言板 CRUD
│   ├── profile/
│   │   └── handler.go     # 博主个人资料
│   ├── changelog/
│   │   └── handler.go     # 更新日志查询
│   └── theme/
│       └── handler.go     # 主题偏好（Redis 存储）
├── pkg/                   # 🔗 基础设施层（可跨领域引用）
│   ├── db/
│   │   └── db.go          # SQLite 连接初始化 + 自动建表（MaxOpen=1, WAL 模式）
│   ├── rds/
│   │   └── rds.go         # Redis 客户端初始化
│   └── middleware/
│       └── middleware.go   # Token 鉴权（30天滚动续期）+ IP 限频
└── docs/                  # 📚 工程文档体系
    ├── architecture.md    # 整体架构概览
    ├── coding-conventions.md  # Go 编码规范
    ├── sqlite.md          # SQLite 使用规范 + 表结构
    ├── http-api.md        # HTTP API 设计规范
    └── api/               # 按领域分组的接口文档
        ├── article.md     # 文章接口
        ├── auth.md        # 认证接口
        ├── album.md       # 相册 + 照片接口
        ├── guestbook.md   # 留言板接口
        ├── profile.md     # 个人资料接口
        └── changelog.md   # 更新日志接口
```

## 领域模块注册表

| 模块 | 包路径 | 说明 | 跨模块依赖 |
|------|--------|------|-----------|
| article | `internal/article/` | 文章 CRUD、点赞 | — |
| album | `internal/album/` | 相册 CRUD、密码验证 | — |
| photo | `internal/photo/` | 照片 CRUD、点赞/踩、密码验证 | → `album.UpdateAlbumCover` |
| auth | `internal/auth/` | 登录、退出、鉴权检查 | → `pkg/middleware` |
| upload | `internal/upload/` | 图片上传（直传+分片） | — |
| guestbook | `internal/guestbook/` | 留言板 CRUD | → `pkg/middleware.GetClientIP` |
| profile | `internal/profile/` | 博主个人资料 | — |
| changelog | `internal/changelog/` | 更新日志查询 | — |
| theme | `internal/theme/` | 主题偏好 | → `upload.RandomString` |

## 路由注册表（完整 API 清单）

### 认证 (`/api/auth`)

| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| POST | `/api/auth/login` | ❌ | `auth.Login` | 账号密码登录，写入 httpOnly Cookie |
| POST | `/api/auth/logout` | ❌ | `auth.Logout` | 退出登录，删除 Redis Token |
| GET | `/api/auth/check` | ✅ | `auth.AuthCheck` | 验证 Token 有效性 |

### 文章 (`/api/articles`)

| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| GET | `/api/articles` | ❌ | `article.GetArticles` | 文章列表（支持 `?title=` 模糊搜索） |
| GET | `/api/articles/like-status-batch` | ❌ | `article.GetArticleLikeStatusBatch` | 批量查询点赞状态（最多 100 篇） |
| GET | `/api/articles/:id` | ❌ | `article.GetArticle` | 文章详情（含 content JSON） |
| GET | `/api/articles/:id/like-status` | ❌ | `article.GetArticleLikeStatus` | 单篇文章点赞状态 |
| POST | `/api/articles/:id/like` | ❌ | `article.LikeArticle` | 点赞/取消切换 |
| POST | `/api/articles` | ✅ | `article.CreateArticle` | 创建文章（UUID 主键） |
| PUT | `/api/articles/:id` | ✅ | `article.UpdateArticle` | 更新文章 |
| DELETE | `/api/articles/:id` | ✅ | `article.DeleteArticle` | 删除文章 |

### 相册 (`/api/albums`)

| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| GET | `/api/albums` | ❌ | `album.GetAlbums` | 相册列表（含 photoCount、hasPassword） |
| GET | `/api/albums/:id/photos` | ❌ | `album.GetAlbumPhotos` | 相册照片列表 |
| POST | `/api/albums/:id/verify-password` | ❌ | `album.VerifyAlbumPassword` | 验证相册密码 |
| POST | `/api/albums` | ✅ | `album.CreateAlbum` | 创建相册 |
| PUT | `/api/albums/:id` | ✅ | `album.UpdateAlbum` | 更新相册 |
| DELETE | `/api/albums/:id` | ✅ | `album.DeleteAlbum` | 删除相册（级联删除照片） |
| POST | `/api/albums/:id/photos` | ✅ | `album.AddAlbumPhoto` | 添加照片（自动更新封面） |

### 照片 (`/api/photos`)

| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| GET | `/api/photos/:id/likes` | ❌ | `photo.GetPhotoLikes` | 获取点赞数 + 当前设备是否已赞 |
| POST | `/api/photos/:id/likes` | ❌ | `photo.PostPhotoLike` | 点赞（幂等，INSERT IGNORE） |
| GET | `/api/photos/:id/dislikes` | ❌ | `photo.GetPhotoDislikes` | 获取踩数 + 当前设备是否已踩 |
| POST | `/api/photos/:id/dislikes` | ❌ | `photo.PostPhotoDislike` | 踩/取消踩（action 参数） |
| POST | `/api/photos/:id/verify-password` | ❌ | `photo.VerifyPhotoPassword` | 验证照片密码 |
| PUT | `/api/photos/:id` | ✅ | `photo.UpdatePhoto` | 更新照片信息 |
| DELETE | `/api/photos/:id` | ✅ | `photo.DeletePhoto` | 删除照片（自动更新封面） |

### 上传 (`/api/upload`，全部需鉴权)

| 方法 | 路径 | Handler | 说明 |
|------|------|---------|------|
| POST | `/api/upload` | `upload.Upload` | 直传图片（jpg/jpeg/png/gif/webp） |
| POST | `/api/upload/chunk` | `upload.UploadChunk` | 上传单个分片 |
| POST | `/api/upload/merge` | `upload.MergeChunks` | 合并分片 → 返回最终 URL |
| DELETE | `/api/upload/chunk` | `upload.CancelChunkUpload` | 取消上传，清理临时分片 |

### 个人资料

| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| GET | `/api/profile` | ❌ | `profile.GetProfile` | 获取博主资料 |
| PUT | `/api/profile` | ✅ | `profile.UpdateProfile` | 更新博主资料 |

### 留言板

| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| GET | `/api/messages` | ❌ | `guestbook.GetMessages` | 留言列表（含 isOwn、canEdit） |
| POST | `/api/messages` | ❌（IP 限频 3/min） | `guestbook.CreateMessage` | 新增/修改留言（每设备一条） |
| PUT | `/api/messages/:id` | ❌ | `guestbook.UpdateMessage` | 修改留言（每日一次，UTC+8） |
| DELETE | `/api/messages/:id` | ✅ | `guestbook.DeleteMessage` | 删除留言（管理员） |

### 主题

| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| GET | `/api/theme` | ❌ | `theme.GetTheme` | 获取主题偏好 |
| POST | `/api/theme` | ❌ | `theme.SaveTheme` | 保存主题偏好 |

### 更新日志

| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| GET | `/api/changelog` | ❌ | `changelog.GetChangelog` | 更新日志列表（按 id DESC） |

## 数据库表注册表

| 表名 | Model 路径 | 主键 | 说明 |
|------|-----------|------|------|
| `articles` | `models/article.go → Article` | `id` varchar(36) UUID | 文章主表 |
| `article_likes` | `models/article.go → ArticleLike` | `id` uint64 自增 | 文章点赞记录 |
| `albums` | `models/album.go → Album` | `id` uint64 自增 | 相册表 |
| `photos` | `models/album.go → Photo` | `id` uint64 自增 | 照片表 |
| `photo_likes` | `models/album.go → PhotoLike` | `id` uint64 自增 | 照片点赞记录 |
| `photo_dislikes` | `models/album.go → PhotoDislike` | `id` uint64 自增 | 照片踩记录 |
| `profile` | `models/misc.go → Profile` | `id` uint64（固定=1） | 博主资料（单行） |
| `messages` | `models/misc.go → Message` | `id` uint64 自增 | 留言板 |
| `changelogs` | `models/misc.go → Changelog` | `id` uint64 自增 | 更新日志 |

## 鉴权机制

- **方式**：Cookie-based Token（httpOnly）
- **Token 存储**：Redis，Key = `auth_token:{64位hex}`，Value = username
- **有效期**：30 天（`TokenTTL = 30 * 24 * time.Hour`）
- **滚动续期**：每次通过 `AuthRequired()` 中间件后，自动刷新 Redis TTL 和 Cookie MaxAge
- **Cookie 配置**：`name=auth_token, httpOnly=true, secure=false, path=/`

## 缓存策略注册表

| Key 模式 | TTL | 存储位置 | 说明 |
|----------|-----|----------|------|
| `auth_token:{token}` | 30 天（滚动续期） | Redis | 用户登录态 |
| `theme:{uid}` | 30 天 | Redis | 用户主题偏好（light/dark） |

> IP 限频使用**内存 Map**（`pkg/middleware/middleware.go`），非 Redis。每 5 分钟清理过期条目。

## 环境变量配置

通过环境变量注入，未设置时使用默认值（见 `config/config.go`）：

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `SERVER_PORT` | `8080` | 服务监听端口 |
| `CORS_ORIGIN` | `https://fatwill.cloud` | 允许跨域来源 |
| `DB_PATH` | `/root/blog-data/blog.db` | SQLite 数据库文件路径 |
| `REDIS_HOST` | `127.0.0.1` | Redis 主机 |
| `REDIS_PORT` | `6379` | Redis 端口 |
| `REDIS_PASSWORD` | *(见配置)* | Redis 密码 |
| `UPLOAD_DIR` | `/root/blog-uploads` | 分片临时文件存储目录 |
| `UPLOAD_TMP_DIR` | `/root/blog-uploads/tmp` | 分片临时目录 |
| `COS_SECRET_ID` | *(必填)* | 腾讯云 COS SecretID |
| `COS_SECRET_KEY` | *(必填)* | 腾讯云 COS SecretKey |
| `COS_BUCKET` | `fatwill-cloud-1253664788` | COS Bucket 名称 |
| `COS_REGION` | `ap-guangzhou` | COS 地域 |
| `COS_BASE_URL` | `https://fatwill-cloud-1253664788.cos.ap-guangzhou.myqcloud.com` | COS 访问域名（后续切换自定义域名只需改此项） |

## 部署信息

| 项目 | 值 |
|------|------|
| 服务器 | 203.195.213.129（2 核 2G） |
| 部署目录 | `/root/blog-backend/` |
| 二进制文件 | `/root/blog-backend/blog-backend` |
| systemd 服务 | `blog-backend` |
| 重启命令 | `systemctl restart blog-backend` |
| 日志查看 | `journalctl -u blog-backend -n 100 --no-pager` |
| SSH Key | `~/.ssh/id_ed25519_vps` |
| 交叉编译 | `GOOS=linux GOARCH=amd64 go build -o blog-backend .` |

## 开发规范速查

### 新增接口流程
1. 在 `models/` 中定义结构体（如需新表/新 DTO）
2. 在 `internal/<domain>/handler.go` 中实现 handler
3. 在 `main.go` → `registerRoutes()` 注册路由
4. 需鉴权的接口加 `authMW` 中间件
5. 编译验证：`go build ./...`
6. 更新接口文档：`docs/api/<domain>.md`
7. 更新本文件路由注册表和领域模块注册表

### Git 提交规范
```
feat(backend agent): 简要描述
fix(backend agent): 简要描述
docs(backend agent): 简要描述
refactor(backend agent): 简要描述
```

### 更新日志规范
- `version`：无 v 前缀（如 `1.9.1`）
- `date`：`YYYY-MM-DD`
- `logs`：JSON 数组，每条 ≤20 字，最多 5 条

## 变更日志

- 2026-04-05: **数据库从 MySQL 迁移至 SQLite** — 使用 modernc.org/sqlite 纯 Go 驱动，无 CGO 依赖；自动建表；WAL 模式优化性能
- 2026-04-05: **图片存储迁移至腾讯云 COS** — 上传直接写入 COS，删除照片/相册时异步清理 COS 对象；新增 `docs/api/upload.md` 接口文档
- 2026-04-04: 创建 AGENTS.md，补全后端项目中枢索引文档
- 2026-04-04: 登录态有效期从 72h 延长至 30 天，支持滚动续期
- 2026-04-04: 补全 docs/ 文档体系（architecture、coding-conventions、mysql、http-api、api/*）
- 2026-04-04: **DDD 分层重构** — 从扁平 handlers/ 重组为 internal/<domain>/ 按领域分组；基础设施移入 pkg/（db、rds、middleware）
