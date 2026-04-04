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
│     /api/* → localhost:8080    /uploads/* → 静态       │
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
│  │  MySQL   │  │  Redis   │  │ 文件系统  │            │
│  │ (GORM)   │  │(go-redis)│  │ (uploads) │            │
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

## 连接池配置（针对 2 核 2G 服务器优化）

### MySQL
- `MaxOpenConns = 10`：最大打开连接数
- `MaxIdleConns = 5`：最大空闲连接数
- GORM 日志级别：`Warn`（仅记录慢查询和错误）

### Redis
- 默认连接池（go-redis 默认 10 个连接）
- 单 DB（DB 0）

## 部署架构

```
Nginx (:443)
    │
    ├── /api/*       → proxy_pass http://127.0.0.1:8080
    ├── /uploads/*   → 静态文件 (/root/blog-uploads/)
    └── 其他          → Nuxt 3 SSR (:3000)

Go Backend (:8080)
    ├── MySQL (127.0.0.1:3306)
    └── Redis (127.0.0.1:6379)
```

- systemd 管理 Go 服务进程
- 交叉编译：`GOOS=linux GOARCH=amd64 go build -o blog-backend .`
