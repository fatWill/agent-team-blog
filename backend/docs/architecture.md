# 整体架构概览

## 系统定位

fatwill 个人博客后端 API 服务，为前端 Nuxt 3 应用提供数据接口。采用单体架构，部署在 2 核 2G 云服务器上。

## 技术架构图

```
┌──────────────────────────────────────────────────────┐
│                     Client (Browser)                  │
│                  Nuxt 3 SSR / CSR                     │
└──────────────────────┬───────────────────────────────┘
                       │ HTTPS
┌──────────────────────▼───────────────────────────────┐
│                      Nginx                            │
│          反向代理 + SSL 终端 + 静态资源                 │
│     /api/* → localhost:8080    图片 → COS (fatwill-cloud-1253664788.cos.ap-guangzhou.myqcloud.com)│
└──────────────────────┬───────────────────────────────┘
                       │ HTTP
┌──────────────────────▼───────────────────────────────┐
│              Go Backend (Gin, :8080)                  │
│                                                       │
│  ┌─────────┐  ┌────────────┐  ┌──────────────────┐   │
│  │ Middleware│  │  Handlers  │  │     Models       │   │
│  │ ┌───────┐│  │ ┌────────┐ │  │ ┌──────────────┐ │   │
│  │ │ Auth  ││  │ │Articles│ │  │ │ Article      │ │   │
│  │ │ Rate  ││  │ │Albums  │ │  │ │ Album/Photo  │ │   │
│  │ │ Limit ││  │ │Photos  │ │  │ │ Message      │ │   │
│  │ │ CORS  ││  │ │Upload  │ │  │ │ Profile      │ │   │
│  │ └───────┘│  │ │Messages│ │  │ │ Changelog    │ │   │
│  └─────────┘  │ │Profile │ │  │ │ JSON (自定义) │ │   │
│               │ │Theme   │ │  │ └──────────────┘ │   │
│               │ │Changelog│ │  └──────────────────┘   │
│               │ └────────┘ │                          │
│               └────────────┘                          │
│                      │                                │
│         ┌────────────┼────────────┐                   │
│         ▼            ▼            ▼                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐            │
│  │  SQLite  │  │  Redis   │  │ 腾讯云COS │            │
│  │ (GORM)   │  │(go-redis)│  │ (图片存储) │            │
│  └──────────┘  └──────────┘  └──────────┘            │
└──────────────────────────────────────────────────────┘
```

## 当前目录结构说明

```
backend/
├── main.go              # 应用入口 + 路由注册
├── config/
│   └── config.go        # 配置结构体（环境变量加载）
├── handlers/            # HTTP 处理器（9 个文件）
├── middleware/           # 中间件（鉴权 + 限频）
├── models/              # 数据模型（GORM struct）
├── utils/               # 基础设施（DB + Redis 初始化）
└── docs/                # 工程文档
```

## 分层说明

### 当前架构：扁平 Handler 模式

当前项目采用**扁平化**的 Handler 模式，所有业务逻辑（参数绑定、校验、数据库操作、响应封装）都在 `handlers/` 层完成：

```
请求 → Middleware(CORS/Auth/RateLimit) → Handler(参数+逻辑+DB+响应) → 数据库
```

| 层级 | 目录 | 职责 |
|------|------|------|
| **入口** | `main.go` | 配置加载、依赖初始化、路由注册 |
| **中间件** | `middleware/` | 鉴权（Token 验证 + 续期）、IP 限频 |
| **处理器** | `handlers/` | 参数绑定 + 业务逻辑 + 数据库操作 + 响应封装 |
| **模型** | `models/` | GORM 结构体定义（纯数据结构，无方法） |
| **基础设施** | `utils/` | MySQL 连接池、Redis 客户端（全局变量） |
| **配置** | `config/` | 环境变量读取 + 默认值 |

### 数据流

```
Client Request
    │
    ▼
Gin Router (main.go registerRoutes)
    │
    ├── CORS Middleware (gin-contrib/cors)
    ├── Auth Middleware (middleware/auth.go)  ← 读/写 Redis + Cookie
    ├── RateLimit Middleware (middleware/ratelimit.go) ← 内存 Map
    │
    ▼
Handler Function (handlers/*.go)
    │
    ├── c.ShouldBindJSON / c.Query / c.Param  ← 参数绑定
    ├── utils.DB.Where(...).Find(...)          ← GORM 查询
    ├── utils.RDB.Get / Set                    ← Redis 操作
    │
    ▼
c.JSON(statusCode, gin.H{...})               ← 响应
```

## 技术选型说明

| 选型 | 理由 |
|------|------|
| **Go + Gin** | 高性能、低内存占用，适合 2 核 2G 服务器 |
| **GORM** | Go 生态最成熟的 ORM，开发效率高 |
| **go-redis** | Go 官方推荐的 Redis 客户端 |
| **Cookie Token** | 个人博客单用户场景，httpOnly Cookie 比 JWT 更安全（防 XSS） |
| **bcrypt** | 密码哈希行业标准，抗彩虹表攻击 |
| **UUID** | 文章 ID 使用 UUID v4，避免自增 ID 被遍历 |
| **内存 Map 限频** | 单实例部署，无需 Redis 限频，内存 Map 性能更好 |
| **环境变量配置** | 通过 systemd `Environment=` 注入，无需配置文件管理 |
| **域名零硬编码** | 所有对外域名（站点、COS、下载白名单）均从环境变量读取，默认值保持线上现状；换域名只需改 systemd `Environment=` 并重启，不动代码 |

## 域名配置（运行时注入）

后端不在代码中硬编码任何对外域名，统一由 `config.Load()` 从环境变量读取：

| 变量 | 作用域 | 消费方 |
|------|--------|--------|
| `CORS_ORIGIN` | 浏览器跨域来源白名单 | `main.go` CORS 中间件 |
| `SITE_URL` | 站点根 URL，拼接绝对页面链接 | `internal/wechat`（`wechat.SetSiteURL`） |
| `COS_BUCKET` / `COS_REGION` / `COS_BASE_URL` | COS SDK 与原始访问域名 | `internal/upload`（`upload.SetCOSConfig`） |
| `COS_CUSTOM_DOMAIN` | 返回给前端的图片 URL 域名 | `internal/upload` |
| `DOWNLOAD_ALLOWED_HOSTS` | 下载代理域名白名单（逗号分隔） | `internal/download`（`download.SetDownloadConfig`） |

注入方式沿用项目既有的**包级 setter** 风格（与 `upload.SetCOSConfig` 一致）：`main()` 加载 `cfg` 后显式调用各领域包的 `Set*` 函数，包内保留与线上一致的默认值兜底，避免忘记注入导致行为突变。

## 连接配置（针对 2 核 2G 服务器优化）

### SQLite
- `MaxOpenConns = 1`：单文件数据库，串行写入
- `MaxIdleConns = 1`
- WAL 模式：支持并发读
- GORM 日志级别：`Warn`（仅记录慢查询和错误）

### Redis
- 默认连接池（go-redis 默认 10 个连接）
- 单 DB（DB 0）

## 部署架构

```
Nginx (:443)
    │
    ├── /api/*       → proxy_pass http://127.0.0.1:8080
    └── 其他          → Nuxt 3 SSR (:3000)

图片存储: 腾讯云 COS (fatwill-cloud-1253664788.cos.ap-guangzhou.myqcloud.com)

Go Backend (:8080)
    ├── SQLite (/root/blog-data/blog.db)
    └── Redis (127.0.0.1:6379)
```

- systemd 管理 Go 服务进程
- 交叉编译：`GOOS=linux GOARCH=amd64 go build -o blog-backend .`
