# AGENTS.md - 后端项目中枢索引

## 项目概述
fatwill 博客后端，基于 Go + Gin 框架构建的 RESTful API 服务，提供文章、相册、留言板、鉴权等核心接口。

## 技术栈速览
- **语言**: Go 1.21+
- **框架**: Gin（gin-gonic/gin）
- **数据库**: MySQL 8.0（go-sql-driver/mysql，database/sql 标准库）
- **缓存**: Redis（go-redis/redis/v9）
- **跨域**: gin-contrib/cors
- **密码哈希**: golang.org/x/crypto/bcrypt
- **包管理**: Go Modules（go.mod）

## 项目结构

```
backend/
├── main.go              # 入口：配置加载、依赖初始化、路由注册
├── config/
│   └── config.go        # 配置结构体 + 环境变量加载（getEnv）
├── handlers/            # HTTP 处理器（业务逻辑层）
│   ├── auth.go          # 登录、退出、鉴权检查
│   ├── articles.go      # 文章 CRUD + 点赞
│   ├── albums.go        # 相册集 CRUD + 密码验证
│   ├── photos.go        # 照片 CRUD + 点赞/踩 + 密码验证
│   ├── upload.go        # 图片上传（直传 + 分片上传）
│   ├── profile.go       # 博主个人资料
│   ├── messages.go      # 留言板 CRUD
│   ├── changelog.go     # 更新日志查询
│   └── theme.go         # 主题偏好（Redis 存储）
├── middleware/
│   ├── auth.go          # Token 鉴权中间件（TokenTTL=30天，滚动续期）
│   └── ratelimit.go     # IP 限频中间件
├── models/
│   ├── article.go       # Article / ArticleLike 结构体
│   ├── album.go         # Album / Photo / Like / Dislike 结构体
│   ├── misc.go          # Profile / Message / Changelog 结构体
│   └── json.go          # JSON 序列化辅助类型
└── utils/
    ├── db.go            # MySQL 连接池初始化（utils.DB）
    └── redis.go         # Redis 客户端初始化（utils.RDB）
```

## 路由注册表（main.go registerRoutes）

### 认证
| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| POST | `/api/auth/login` | ❌ | 账号密码登录，写入 httpOnly Cookie |
| POST | `/api/auth/logout` | ❌ | 退出登录，删除 Token |
| GET | `/api/auth/check` | ✅ | 验证当前 Cookie Token 是否有效 |

### 文章
| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/articles` | ❌ | 文章列表（支持 `?title=` 筛选） |
| GET | `/api/articles/like-status-batch` | ❌ | 批量查询点赞状态 |
| GET | `/api/articles/:id` | ❌ | 文章详情（含 content） |
| GET | `/api/articles/:id/like-status` | ❌ | 单篇文章点赞状态 |
| POST | `/api/articles/:id/like` | ❌ | 文章点赞/取消（切换） |
| POST | `/api/articles` | ✅ | 创建文章 |
| PUT | `/api/articles/:id` | ✅ | 更新文章 |
| DELETE | `/api/articles/:id` | ✅ | 删除文章 |

### 相册
| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/albums` | ❌ | 相册列表（含 photoCount、coverUrl） |
| GET | `/api/albums/:id/photos` | ❌ | 相册照片列表 |
| POST | `/api/albums/:id/verify-password` | ❌ | 验证相册密码 |
| POST | `/api/albums` | ✅ | 创建相册 |
| PUT | `/api/albums/:id` | ✅ | 更新相册 |
| DELETE | `/api/albums/:id` | ✅ | 删除相册（级联删除照片） |
| POST | `/api/albums/:id/photos` | ✅ | 向相册添加照片 |

### 照片
| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/photos/:id/likes` | ❌ | 获取照片点赞数 |
| POST | `/api/photos/:id/likes` | ❌ | 照片点赞（每设备一次） |
| GET | `/api/photos/:id/dislikes` | ❌ | 获取照片踩数 |
| POST | `/api/photos/:id/dislikes` | ❌ | 照片踩（每设备一次） |
| POST | `/api/photos/:id/verify-password` | ❌ | 验证照片密码 |
| PUT | `/api/photos/:id` | ✅ | 更新照片信息 |
| DELETE | `/api/photos/:id` | ✅ | 删除照片（自动更新封面） |

### 上传（全部需鉴权）
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/upload` | 直传（≤2MB） |
| POST | `/api/upload/chunk` | 上传单个分片 |
| POST | `/api/upload/merge` | 合并分片，返回最终 URL |
| DELETE | `/api/upload/chunk` | 取消上传，清理临时分片 |

### 其他
| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/profile` | ❌ | 获取博主资料 |
| PUT | `/api/profile` | ✅ | 更新博主资料 |
| GET | `/api/messages` | ❌ | 留言列表 |
| POST | `/api/messages` | ❌（限频） | 新增留言（IP 每分钟 3 次） |
| PUT | `/api/messages/:id` | ❌ | 修改留言（每日一次） |
| DELETE | `/api/messages/:id` | ✅ | 删除留言 |
| GET | `/api/theme` | ❌ | 获取主题偏好 |
| POST | `/api/theme` | ❌ | 保存主题偏好 |
| GET | `/api/changelog` | ❌ | 更新日志列表 |

## 鉴权机制（middleware/auth.go）

- **Token 存储**：Redis，Key = `auth_token:{token}`，Value = username
- **Token 有效期**：`TokenTTL = 30 * 24 * time.Hour`（30 天）
- **滚动续期**：每次通过 `AuthRequired()` 中间件验证后，自动刷新 Redis TTL 和 Cookie MaxAge
- **Cookie 名称**：`auth_token`，`HttpOnly=true`，`Secure=false`
- **常量**：
  ```go
  const (
      TokenTTL       = 30 * 24 * time.Hour
      TokenKeyPrefix = "auth_token:"
      CookieName     = "auth_token"
  )
  ```
- **核心函数**：
  - `SaveToken(token, username string) error`
  - `VerifyToken(token string) (string, error)`
  - `DeleteToken(token string) error`
  - `AuthRequired() gin.HandlerFunc`

## 数据库表注册表

| 表名 | 说明 | 主键 |
|------|------|------|
| `articles` | 文章表 | `id` (varchar(36) UUID) |
| `article_likes` | 文章点赞表 | `id` (int auto_increment) |
| `albums` | 相册集表 | `id` (int auto_increment) |
| `photos` | 照片表 | `id` (int auto_increment) |
| `photo_likes` | 照片点赞表 | `id` (int auto_increment) |
| `photo_dislikes` | 照片踩表 | `id` (int auto_increment) |
| `profile` | 博主资料表 | `id` (int, 固定为1) |
| `messages` | 留言板表 | `id` (int auto_increment) |
| `changelogs` | 更新日志表 | `id` (bigint auto_increment) |

## 环境变量配置

通过环境变量注入，未设置时使用默认值（见 `config/config.go`）：

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `SERVER_PORT` | `8080` | 服务监听端口 |
| `CORS_ORIGIN` | `https://fatwill.cloud` | 允许跨域来源 |
| `DB_HOST` | `127.0.0.1` | MySQL 主机 |
| `DB_PORT` | `3306` | MySQL 端口 |
| `DB_USER` | `root` | MySQL 用户名 |
| `DB_PASSWORD` | — | MySQL 密码 |
| `DB_NAME` | `blog` | 数据库名 |
| `REDIS_HOST` | `127.0.0.1` | Redis 主机 |
| `REDIS_PORT` | `6379` | Redis 端口 |
| `REDIS_PASSWORD` | — | Redis 密码 |
| `UPLOAD_DIR` | `/root/blog-uploads` | 图片存储目录 |
| `UPLOAD_TMP_DIR` | `/root/blog-uploads/tmp` | 分片临时目录 |

生产环境通过 systemd service 的 `Environment=` 字段注入（`/etc/systemd/system/blog-backend.service`）。

## 开发规范

### 新增接口流程
1. 在 `models/` 中定义请求/响应结构体
2. 在 `handlers/` 对应文件中实现 handler 函数
3. 在 `main.go` 的 `registerRoutes` 中注册路由
4. 需鉴权的接口加上 `auth` 中间件参数

### 编译与验证
```bash
cd /Users/fatwill/github/agent-team-blog/backend
go build .
```

### 交叉编译（用于部署到 Linux 服务器）
```bash
GOOS=linux GOARCH=amd64 go build -o blog-backend .
```

### Git 提交规范（用户明确要求）
- **后端独立提交**：后端开发完成后，由「博客后端」子 Agent 自行执行 git add + git commit + git push
- commit message 格式：`feat(backend agent): 简要描述` 或 `fix(backend agent): 简要描述`
- push 前需执行：`ssh-add ~/.ssh/id_ed25519_github`
- 远程仓库：`git@github.com:fatWill/agent-team-blog.git`（SSH 格式，main 分支）

### 更新日志规范
每次发版后需向 MySQL `changelogs` 表插入一条记录：
- `version`：版本号（格式 `X.Y.Z`，默认递增中间位）
- `date`：发版日期（`YYYY-MM-DD`）
- `logs`：JSON 数组，每条 ≤20 字，最多 5 条，可用 Emoji 前缀

## 部署信息

- **服务器**：203.195.213.129
- **部署目录**：`/root/blog-backend/`
- **二进制文件**：`/root/blog-backend/blog-backend`
- **systemd 服务**：`blog-backend`
- **重启命令**：`systemctl restart blog-backend`
- **日志查看**：`journalctl -u blog-backend -n 100 --no-pager`
- **SSH Key**：`~/.ssh/id_ed25519_vps`

## 变更日志
- 2026-04-04: 创建 AGENTS.md，补全后端项目中枢索引文档
- 2026-04-04: 登录态有效期从 72h 延长至 30 天，支持滚动续期（TokenTTL 常量统一管理）
