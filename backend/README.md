# Blog Backend (Go)

博客系统 Go 后端服务，从 Nuxt 3 Server API 迁移而来。

## 技术栈

- **Go 1.22+** + **Gin** HTTP 框架
- **GORM** ORM（MySQL 8.0）
- **go-redis** v9（Redis 7.0+）
- **bcrypt**（密码加密）

## 快速启动

### 环境要求

- Go 1.22+
- MySQL 8.0+（已有 `blog` 数据库和表结构）
- Redis 6.0+

### 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `SERVER_PORT` | `8080` | 服务监听端口 |
| `CORS_ORIGIN` | `https://fatwill.cloud` | 允许的跨域来源 |
| `DB_HOST` | `127.0.0.1` | MySQL 地址 |
| `DB_PORT` | `3306` | MySQL 端口 |
| `DB_USER` | `root` | MySQL 用户名 |
| `DB_PASSWORD` | *(必填)* | MySQL 密码 |
| `DB_NAME` | `blog` | MySQL 数据库名 |
| `REDIS_HOST` | `127.0.0.1` | Redis 地址 |
| `REDIS_PORT` | `6379` | Redis 端口 |
| `REDIS_PASSWORD` | *(必填)* | Redis 密码 |
| `UPLOAD_DIR` | `/root/blog-uploads` | 图片上传目录 |
| `UPLOAD_TMP_DIR` | `/root/blog-uploads/tmp` | 分片上传临时目录 |
| `ADMIN_USERNAME` | `admin` | 管理后台用户名 |
| `ADMIN_PASSWORD` | *(必填)* | 管理后台密码 |
| `COS_ID` | *(必填)* | 腾讯云 COS SecretId |
| `COS_KEY` | *(必填)* | 腾讯云 COS SecretKey |
| `COS_BUCKET` | `fatwill-cloud-1253664788` | COS 存储桶名称 |
| `COS_REGION` | `ap-guangzhou` | COS 地域 |
| `COS_BASE_URL` | `https://assets.fatwill.cloud` | COS CDN 访问域名 |

### 启动

```bash
cd backend

# 设置必要的环境变量
export DB_PASSWORD="<your-db-password>"
export REDIS_PASSWORD="<your-redis-password>"
export ADMIN_USERNAME="<your-admin-username>"
export ADMIN_PASSWORD="<your-admin-password>"

go run main.go
```

服务默认监听 `:8080`，Nginx 需将 `/api/*` 请求转发到此端口。

### 编译

```bash
go build -o blog-backend main.go
./blog-backend
```

## API 路由

### 认证
- `POST /api/auth/login` — 登录
- `GET /api/auth/check` — 检查登录状态（需鉴权）
- `POST /api/auth/logout` — 登出

### 文章
- `GET /api/articles` — 文章列表
- `GET /api/articles/:id` — 文章详情
- `POST /api/articles` — 创建文章（需鉴权）
- `PUT /api/articles/:id` — 更新文章（需鉴权）
- `DELETE /api/articles/:id` — 删除文章（需鉴权）
- `POST /api/articles/:id/like` — 点赞/取消
- `GET /api/articles/:id/like-status` — 点赞状态
- `GET /api/articles/like-status-batch` — 批量点赞状态

### 相册
- `GET /api/albums` — 相册列表
- `POST /api/albums` — 创建相册（需鉴权）
- `PUT /api/albums/:id` — 更新相册（需鉴权）
- `DELETE /api/albums/:id` — 删除相册（需鉴权）
- `GET /api/albums/:id/photos` — 照片列表
- `POST /api/albums/:id/photos` — 添加照片（需鉴权）
- `POST /api/albums/:id/verify-password` — 验证相册密码

### 照片
- `PUT /api/photos/:id` — 更新照片（需鉴权）
- `DELETE /api/photos/:id` — 删除照片（需鉴权）
- `GET /api/photos/:id/likes` — 点赞数
- `POST /api/photos/:id/likes` — 点赞
- `GET /api/photos/:id/dislikes` — 踩数
- `POST /api/photos/:id/dislikes` — 踩/取消踩
- `POST /api/photos/:id/verify-password` — 验证照片密码

### 上传（全部需鉴权）
- `POST /api/upload` — 普通上传
- `POST /api/upload/chunk` — 分片上传
- `POST /api/upload/merge` — 合并分片
- `DELETE /api/upload/chunk` — 取消分片上传

### 其他
- `GET /api/profile` — 个人资料
- `PUT /api/profile` — 更新资料（需鉴权）
- `GET /api/messages` — 留言列表
- `POST /api/messages` — 发布留言（IP 限频）
- `PUT /api/messages/:id` — 修改留言
- `DELETE /api/messages/:id` — 删除留言（需鉴权）
- `GET /api/theme` — 获取主题
- `POST /api/theme` — 保存主题
- `GET /api/changelog` — 更新日志

## 目录结构

```
backend/
├── main.go                # 入口，路由注册
├── config/
│   └── config.go          # 配置加载（环境变量）
├── internal/              # 业务核心（按领域分组）
│   ├── article/handler.go
│   ├── album/handler.go
│   ├── photo/handler.go
│   ├── auth/handler.go
│   ├── upload/handler.go
│   ├── guestbook/handler.go
│   ├── profile/handler.go
│   ├── changelog/handler.go
│   └── theme/handler.go
├── pkg/                   # 基础设施
│   ├── db/db.go           # MySQL 连接
│   ├── rds/rds.go         # Redis 连接
│   └── middleware/middleware.go  # 鉴权 + 限频
├── models/                # 数据模型
│   ├── article.go
│   ├── album.go
│   ├── json.go
│   └── misc.go
└── README.md
```
