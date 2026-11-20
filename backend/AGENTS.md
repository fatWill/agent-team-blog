# AGENTS.md - 后端项目中枢索引

> 仓库级运维、EdgeOne/CDN 缓存架构与总规范见根目录 [AGENTS.md](../AGENTS.md)。

## ⚠️ 文档同步与提交铁律（必守）

- **同步更新**：任何后端代码改动（新增/修改 API 路由、领域模块 `internal/`、基础包 `pkg/`、SQLite 表结构或迁移、环境变量、部署方式），**必须在同一次改动中同步更新本文件，禁止长期滞后**。与原内容冲突时直接修改原条目，并在底部「变更日志」追加记录。
- **改完直接提交**：改动并验证完成后**直接 `git commit && git push`，不用再问**（只暂存本次改动文件，不误提交无关变更）。

---

## 项目概述

fatwill 个人博客后端服务，基于 Go + Gin + GORM + SQLite + Redis 构建的 RESTful API，提供文章、相册、成长日记、装修（文章/预算/踩坑/材料清单）、留言板、鉴权、COS 文件上传、下载代理、PV/UV 与性能监控、微信公众号同步、图一乐（老黄历 + AG-UI SSE 运势聊天）等核心能力。服务部署在 2 核 2GiB 1Mbps 云服务器上，需严格控制内存与带宽开销。

## 技术栈速览

| 组件 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.26.1 |
| Web 框架 | Gin | v1.12.0 |
| ORM / 数据库 | GORM + `glebarez/sqlite`（纯 Go `modernc.org/sqlite`，无 CGO） | v1.31.1 + v1.11.0 |
| 缓存 | Redis（`go-redis/v9`） | v9.18.0 |
| 对象存储 | 腾讯云 COS（`tencentyun/cos-go-sdk-v5`） | v0.7.73 |
| IP 地理库 | `ip2region`（离线库 `data/ip2region.xdb`） | golang binding |
| 跨域 / 鉴权 | `gin-contrib/cors` + `golang.org/x/crypto/bcrypt` + `google/uuid` | v1.7.7 / v0.49.0 / v1.6.0 |
| 农历 / 爬虫 | `6tail/lunar-go` + `PuerkitoBio/goquery` | v1.4.6 / v1.11.0 |
| 包管理 | Go Modules | `go.mod` |

## 最高优先级读取

Agent 收到后端任务后，**按以下优先级读取**：
1. **本文件** `backend/AGENTS.md` 与根 [AGENTS.md](../AGENTS.md)
2. `docs/coding-conventions.md` — Go 编码规范
3. `docs/sqlite.md` — SQLite 使用规范 + 表结构说明
4. `docs/http-api.md` — HTTP API 设计规范 + 错误码体系
5. `docs/architecture.md` — 整体架构概览

## 当前目录结构

```
backend/
├── AGENTS.md                        # 🔑 后端中枢索引（本文件）
├── README.md / go.mod / go.sum
├── main.go                          # 入口：配置加载、DB/备份/Redis/ipgeo/微信worker 初始化、路由注册
├── cmd/almanac-crawler/             # 老黄历抓取与入库工具（main.go / import.go）
├── config/config.go                 # 配置结构体 + 环境变量加载
├── data/ip2region.xdb               # IP 离线地理位置数据库
├── deploy/env.production.template   # 生产环境变量清单模板
├── scripts/seed_changelog_2.15.0.sql
├── models/                          # 📦 数据模型层（GORM struct + DTO）
│   ├── article.go                   #   Article / ArticleListItem / ArticleLike / WechatSyncLog
│   ├── album.go                     #   Album / Photo / PhotoLike / PhotoDislike
│   ├── misc.go                      #   Profile / Message / Changelog
│   ├── pv.go                        #   PageView / PVTrendItem / TopPageItem / PVOverview
│   ├── perf.go                      #   前端 Web Vitals 性能指标模型
│   ├── renovation.go                #   装修知识文章模型
│   ├── budget.go                    #   装修成本预算模型（含 actual 字段）
│   ├── pitfall.go                   #   装修踩坑日记模型
│   ├── material.go                  #   装修材料清单模型
│   ├── growth_diary.go              #   成长日记模型
│   ├── almanac.go                   #   老黄历日表模型
│   ├── tuyile_share.go              #   图一乐对话分享快照模型
│   └── json.go                      #   自定义 JSON 序列化类型（SQLite JSON 文本字段）
├── internal/                        # 🏗️ 业务核心（按领域分组）
│   ├── article/                     #   文章 CRUD、随机文章、点赞、阅读量
│   ├── album/                       #   相册 CRUD、密码验证、封面维护
│   ├── photo/                       #   照片 CRUD、点赞/踩、密码验证
│   ├── auth/                        #   登录、退出、鉴权检查
│   ├── upload/                      #   腾讯云 COS 上传（直传 + 分片）
│   ├── download/                    #   文件下载代理（白名单防 SSRF）
│   ├── guestbook/                   #   留言板 CRUD
│   ├── profile/                     #   博主个人资料
│   ├── changelog/                   #   更新日志查询
│   ├── theme/                       #   主题偏好（Redis 存储）
│   ├── pv/                          #   PV/UV 统计、访问日志、访客地理分布
│   ├── perf/                        #   前端性能指标上报与分析查询
│   ├── renovation/                  #   装修知识文章 CRUD
│   ├── budget/                      #   装修预算查询与批量编辑
│   ├── pitfall/                     #   装修踩坑日记分类/条目与批量编辑
│   ├── material/                    #   装修材料清单 CRUD 与排序
│   ├── growthdiary/                 #   成长日记按月查询与批量编辑
│   ├── wechat/                      #   微信公众号 JS-SDK 签名、草稿同步 worker、素材与日志
│   ├── almanac/                     #   老黄历查询（单日/月度，缺库时算法现算）
│   ├── tuyile/                      #   图一乐 AG-UI SSE 聊天、LLM 代理、黄历工具闭环、限流、分享
│   └── backup/                      #   SQLite 启动备份 + 每日 03:00 VACUUM INTO 定时备份
├── pkg/                             # 🔗 基础设施层
│   ├── db/                          #   db.go（SQLite WAL + 自动建表/兼容迁移）、changelog_seed.go、renovation_seed.go
│   ├── rds/rds.go                   #   Redis 客户端初始化
│   ├── middleware/middleware.go     #   Token 鉴权（30 天滚动续期）+ 内存 IP 限频
│   ├── ipgeo/ipgeo.go               #   ip2region 封装
│   ├── nginx/purge.go               #   Nginx 本地 SSR 缓存清理工具
│   └── almanac/                     #   农历/黄历算法与网页解析（algo.go / site.go）
└── docs/                            # 📚 工程与分领域 API 文档
```

## 领域模块注册表

| 模块 | 包路径 | 说明 | 关键依赖 |
|------|--------|------|----------|
| `article` | `internal/article/` | 文章 CRUD、随机文章、点赞、阅读量 | `pkg/db`、`pkg/nginx`、`internal/wechat` |
| `album` | `internal/album/` | 相册 CRUD、密码验证、封面维护 | `pkg/db` |
| `photo` | `internal/photo/` | 照片 CRUD、点赞/踩、密码验证 | `album.UpdateAlbumCover` |
| `auth` | `internal/auth/` | 登录、退出、鉴权检查 | `pkg/middleware`、`pkg/rds` |
| `upload` | `internal/upload/` | 腾讯云 COS 直传与分片上传、清理 | COS SDK |
| `download` | `internal/download/` | 文件下载代理（`DOWNLOAD_ALLOWED_HOSTS` 白名单校验） | `config.DownloadConfig` |
| `guestbook` | `internal/guestbook/` | 留言板 CRUD | `pkg/middleware` |
| `profile` | `internal/profile/` | 博主个人资料 | `pkg/db` |
| `changelog` | `internal/changelog/` | 更新日志查询 | `pkg/db` |
| `theme` | `internal/theme/` | 主题偏好（Redis） | `pkg/rds`、`upload.RandomString` |
| `pv` | `internal/pv/` | PV/UV 统计、日志、访客地理分布 | `pkg/ipgeo`、`pkg/rds`、`pkg/middleware` |
| `perf` | `internal/perf/` | 前端 Web Vitals 性能指标上报、概览、页面聚合、趋势、明细 | `pkg/db` |
| `renovation` | `internal/renovation/` | 装修知识文章 CRUD | `pkg/db` |
| `budget` | `internal/budget/` | 装修预算查询与批量编辑（含写前备份） | `pkg/db`、`internal/backup` |
| `pitfall` | `internal/pitfall/` | 装修踩坑日记分类/条目与批量编辑 | `pkg/db` |
| `material` | `internal/material/` | 装修材料清单 CRUD 与排序（全鉴权） | `pkg/db` |
| `growthdiary` | `internal/growthdiary/` | 成长日记按月查询与批量编辑 | `pkg/db` |
| `wechat` | `internal/wechat/` | 公众号 JS-SDK 签名、草稿同步 worker、素材与日志 | `pkg/rds`、`pkg/db`、`models` |
| `almanac` | `internal/almanac/` | 老黄历单日/月度查询、缺库算法兜底 | `pkg/almanac`、`pkg/db` |
| `tuyile` | `internal/tuyile/` | 图一乐 AG-UI SSE 聊天、黄历工具闭环、内存限流、分享快照 | `internal/almanac`、`pkg/almanac`、`pkg/middleware` |
| `backup` | `internal/backup/` | SQLite `VACUUM INTO` 快照备份（`/root/blog-backup`，每 trigger 保留 30 份，每日 03:00） | `pkg/db` |

## 路由注册表（与 `main.go` 完全对齐）

### 认证 (`/api/auth`)
| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| POST | `/api/auth/login` | ❌ | `auth.Login` | 账号密码登录，写 httpOnly Cookie |
| POST | `/api/auth/logout` | ❌ | `auth.Logout` | 退出登录，删除 Redis Token |
| GET | `/api/auth/check` | ✅ | `auth.AuthCheck` | 验证 Token 有效性 |

### 文章 (`/api/articles`)
| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| GET | `/api/articles` | ❌ | `article.GetArticles` | 文章列表（支持 `?title=` 搜索） |
| GET | `/api/articles/random` | ❌ | `article.GetRandomArticle` | 随机文章（须在 `/:id` 前注册） |
| GET | `/api/articles/like-status-batch` | ❌ | `article.GetArticleLikeStatusBatch` | 批量查询点赞状态（最多 100 篇） |
| GET | `/api/articles/:id` | ❌ | `article.GetArticle` | 文章详情 |
| GET | `/api/articles/:id/like-status` | ❌ | `article.GetArticleLikeStatus` | 单篇点赞状态 |
| POST | `/api/articles/:id/like` | ❌ | `article.LikeArticle` | 点赞/取消切换 |
| POST | `/api/articles/:id/view` | ❌ | `article.ViewArticle` | 阅读量统计上报 |
| POST | `/api/articles` | ✅ | `article.CreateArticle` | 创建文章（UUID 主键） |
| PUT | `/api/articles/:id` | ✅ | `article.UpdateArticle` | 更新文章 |
| DELETE | `/api/articles/:id` | ✅ | `article.DeleteArticle` | 删除文章 |

### 相册与照片 (`/api/albums`, `/api/photos`)
| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| GET | `/api/albums` | ❌ | `album.GetAlbums` | 相册列表（含 `photoCount`/`hasPassword`） |
| GET | `/api/albums/:id/photos` | ❌ | `album.GetAlbumPhotos` | 相册照片列表 |
| POST | `/api/albums/:id/verify-password` | ❌ | `album.VerifyAlbumPassword` | 验证相册密码 |
| POST / PUT / DELETE | `/api/albums` / `/:id` | ✅ | `album.{Create,Update,Delete}Album` | 相册管理 |
| POST | `/api/albums/:id/photos` | ✅ | `album.AddAlbumPhoto` | 添加照片并更新封面 |
| GET / POST | `/api/photos/:id/likes` | ❌ | `photo.{Get,Post}PhotoLike` | 照片点赞 |
| GET / POST | `/api/photos/:id/dislikes` | ❌ | `photo.{Get,Post}PhotoDislike` | 照片踩 |
| POST | `/api/photos/:id/verify-password` | ❌ | `photo.VerifyPhotoPassword` | 照片密码验证 |
| PUT / DELETE | `/api/photos/:id` | ✅ | `photo.{Update,Delete}Photo` | 照片管理 |

### 上传、下载、资料、留言、主题、更新日志
| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| POST | `/api/upload` | ✅ | `upload.Upload` | COS 小文件直传 |
| POST / POST / DELETE | `/api/upload/chunk`, `/merge`, `/chunk` | ✅ | `upload.{UploadChunk,MergeChunks,CancelChunkUpload}` | 分片上传与清理 |
| GET | `/api/download` | ❌ | `download.Proxy` | 白名单域名文件下载代理 |
| GET / PUT | `/api/profile` | ❌ / ✅ | `profile.{Get,Update}Profile` | 博主资料 |
| GET / POST / DELETE | `/api/messages`, `/messages/:id` | ❌ / 限频3/min / ✅ | `guestbook.{GetMessages,CreateMessage,DeleteMessage}` | 留言板 |
| GET / POST | `/api/theme` | ❌ | `theme.{Get,Save}Theme` | 主题偏好 |
| GET | `/api/changelog` | ❌ | `changelog.GetChangelog` | 更新日志列表 |

### PV/UV 统计 (`/api/pv`) 与性能监控 (`/api/perf`)
| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| POST | `/api/pv/record` | ❌ | `pv.RecordPV` | 上报访问记录（Redis 60s 防刷） |
| GET | `/api/pv/{trend,top-pages,logs,overview,geo}` | ✅ | `pv.{GetTrend,GetTopPages,GetLogs,GetOverview,GetGeoDistribution}` | 统计报表与地理分布 |
| POST | `/api/perf/report` | ❌ | `perf.Report` | 前端 Web Vitals 上报 |
| GET | `/api/perf/{overview,pages,trend,logs}` | ✅ | `perf.{GetOverview,GetPages,GetTrend,GetLogs}` | 性能监控报表 |

### 装修与成长日记 (`/api/renovation/*`, `/api/materials`, `/api/growth-diary`)
| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| GET | `/api/renovation/articles`, `/:id` | ❌ | `renovation.{GetArticles,GetArticle}` | 装修文章公开查询 |
| POST / PUT / DELETE | `/api/renovation/articles`, `/:id` | ✅ | `renovation.{Create,Update,Delete}Article` | 装修文章管理 |
| GET / POST | `/api/renovation/budget/items`, `/items/batch` | ❌ / ✅ | `budget.{GetItems,BatchEditItems}` | 装修预算查询与批量编辑 |
| GET / GET / POST | `/api/renovation/pitfall/categories`, `/items`, `/items/batch` | ❌ / ❌ / ✅ | `pitfall.{GetCategories,GetItems,BatchEditItems}` | 踩坑日记查询与批量编辑 |
| GET / GET / POST | `/api/growth-diary/items`, `/months`, `/items/batch` | ❌ / ❌ / ✅ | `growthdiary.{GetItems,GetMonths,BatchEditItems}` | 成长日记查询与批量编辑 |
| GET / POST / PUT / DELETE / PUT | `/api/materials`, `/:id`, `/:id/sort` | ✅（全部需鉴权） | `material.*` | 材料清单全量管理与排序 |

### 微信 (`/api/wechat`, `/api/admin/wechat`, `/api/admin/articles`)
| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| GET | `/api/wechat/jssdk-config` | ❌ | `wechat.GetJSSDKConfig` | 微信公众号 JS-SDK 签名配置 |
| GET | `/api/admin/wechat/{access-token-status,server-ip,sync-logs}` | ✅ | `wechat.{GetAccessTokenStatus,GetServerIP,ListAllSyncLogs}` | 微信 token 状态、出口 IP、全量同步日志 |
| POST / GET | `/api/admin/articles/:id/sync-wechat`, `/wechat-sync-logs` | ✅ | `wechat.{SyncWechatHandler,GetSyncLogsHandler}` | 单篇同步草稿箱与同步日志 |

### 玩具：图一乐 (`/api/toys/tuyile`，详见 `docs/api/tuyile.md`)
| 方法 | 路径 | 鉴权 | Handler | 说明 |
|------|------|------|---------|------|
| GET | `/api/toys/tuyile/almanac` | ❌ | `almanac.GetDay` | 单日老黄历（默认北京时间今天） |
| GET | `/api/toys/tuyile/almanac/month` | ❌ | `almanac.GetMonth` | 月度老黄历精简字段 |
| POST | `/api/toys/tuyile/chat` | ❌（IP 限流） | `tuyile.Chat` | AG-UI SSE 运势聊天（含 `query_almanac` 工具闭环） |
| POST / GET | `/api/toys/tuyile/share`, `/share/:hash` | ❌（创建限流） | `tuyile.{CreateShare,GetShare}` | 对话分享快照 |

## 数据库表注册表（共 19 张表，与 `pkg/db/db.go` 对齐）

| 表名 | Model 路径 | 主键 | 说明 |
|------|-----------|------|------|
| `articles` | `models/article.go` | `id` TEXT UUID | 文章主表（含 `views` 与 5 个微信同步列） |
| `article_likes` | `models/article.go` | `id` 自增 | 文章点赞记录 |
| `albums` | `models/album.go` | `id` 自增 | 相册表 |
| `photos` | `models/album.go` | `id` 自增 | 照片/视频表（含 `media_type`/`thumbnail_url`/`duration`） |
| `photo_likes` / `photo_dislikes` | `models/album.go` | `id` 自增 | 照片点赞与踩记录 |
| `profile` | `models/misc.go` | `id`（固定=1） | 博主资料单行表（含 `name` 列） |
| `messages` | `models/misc.go` | `id` 自增 | 留言板（`device_id` 唯一） |
| `changelogs` | `models/misc.go` | `id` 自增 | 更新日志（`version` 唯一，启动幂等播种） |
| `page_views` | `models/pv.go` | `id` 自增 | PV/UV 页面访问记录 |
| `perf_metrics` | `models/perf.go` | `id` 自增 | 前端 Web Vitals 性能指标 |
| `renovation_articles` | `models/renovation.go` | `id` 自增 | 装修知识文章（空库自动播种） |
| `budget_items` | `models/budget.go` | `id` 自增 | 装修成本预算（含 `actual` 实际支出列） |
| `pitfall_items` | `models/pitfall.go` | `id` 自增 | 装修踩坑日记 |
| `material_items` | `models/material.go` | `id` 自增 | 装修材料清单（含 `sort_order`） |
| `growth_diary_items` | `models/growth_diary.go` | `id` 自增 | 成长日记（`images`/`videos` JSON、`happened_at`） |
| `wechat_sync_logs` | `models/article.go` | `id` 自增 | 微信草稿同步日志 |
| `almanac_days` | `models/almanac.go` | `id` 自增（`date` 唯一） | 老黄历知识库 |
| `tuyile_shares` | `models/tuyile_share.go` | `id` 自增（`hash` 唯一） | 图一乐对话分享快照 |

## 鉴权与缓存策略

- **鉴权**：Cookie-based Token（`name=auth_token, httpOnly=true, path=/`），Redis 存储 `auth_token:{64位hex}` → `username`，TTL 30 天滚动续期。
- **Redis Key**：`auth_token:{token}`（30d）、`theme:{uid}`（30d）、`pv:{device_id}:{path}`（60s）、`wechat:access_token`（7000s）、`wechat:jsapi_ticket`（7000s）。
- **内存限流**：通用 `middleware.RateLimit` 与 `internal/tuyile/ratelimit.go` 均为进程内 Map，不占 Redis。

## 环境变量配置（与 `config/config.go` 对齐）

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `SERVER_PORT` | `8080` | 监听端口 |
| `CORS_ORIGIN` | `https://fatwill.cn,https://www.fatwill.cn` | 允许跨域来源（逗号分隔） |
| `SITE_URL` | `https://fatwill.cn` | 站点根 URL（微信同步/分享链接） |
| `IP2REGION_PATH` | `data/ip2region.xdb` | ip2region 离线库路径 |
| `DB_PATH` | `/root/blog-data/blog.db` | SQLite 数据库路径（MaxOpen=1, WAL） |
| `REDIS_HOST` / `REDIS_PORT` / `REDIS_PASSWORD` | `127.0.0.1` / `6379` / *(空)* | Redis 配置 |
| `UPLOAD_DIR` / `UPLOAD_TMP_DIR` | `/root/blog-uploads` / `.../tmp` | 上传与分片临时目录 |
| `COS_SECRET_ID` / `COS_SECRET_KEY` | *(必填)* | 腾讯云 COS 密钥 |
| `COS_BUCKET` / `COS_REGION` | `fatwill-cloud-1253664788` / `ap-guangzhou` | COS 存储桶与地域（勿随域名改动） |
| `COS_BASE_URL` | `https://fatwill-cloud-1253664788.cos.ap-guangzhou.myqcloud.com` | COS 原始域名 |
| `COS_CUSTOM_DOMAIN` / `COS_LEGACY_DOMAINS` | `https://assets.fatwill.cn` / `https://cdn.fatwill.cn` | 新写入域名与删除反解历史域名 |
| `DOWNLOAD_ALLOWED_HOSTS` | `assets.fatwill.cn,cdn.fatwill.cn,pic.fatwill.cn,fatwill-cloud-1253664788.cos.ap-guangzhou.myqcloud.com` | `/api/download` 代理白名单 |
| `WECHAT_APP_ID` / `WECHAT_APP_SECRET` | *(必填)* | 微信公众号配置 |
| `LLM_BASE_URL` / `LLM_API_KEY` / `LLM_MODEL` | *(空)* | 图一乐 OpenAI 兼容上游（生产智谱 `glm-5.3-flash`，禁止 woa.com） |
| `LLM_REASONING_EFFORT` / `LLM_THINKING_TYPE` / `LLM_EXTRA_BODY` / `LLM_DISABLE_TOOLS` | *(空)* | 思考控制与厂商扩展参数（GLM-5.3 生产用 `low`，勿传 `disabled`） |
| `LLM_TIMEOUT_SECONDS` / `LLM_HEADER_TIMEOUT_SECONDS` / `LLM_MAX_TOKENS` / `LLM_MAX_TOOL_ROUNDS` | `120` / `90` / `2000` / `3` | 流式超时与工具轮数上限 |
| `TUYILE_RATE_PER_MIN` / `TUYILE_RATE_PER_DAY` / `TUYILE_MAX_CONCURRENCY` | `10` / `100` / `8` | 图一乐 chat 限流与并发保护 |
| `TUYILE_SHARE_RATE_PER_MIN` / `TUYILE_SHARE_RATE_PER_DAY` | `5` / `30` | 图一乐 share 创建限流 |

> 生产环境变量清单模板见 `deploy/env.production.template`；**禁止读取或输出** VPS 上 `/etc/systemd/system/blog-backend.service` 的密钥取值。

## 部署信息

| 项目 | 值 |
|------|------|
| 服务器 | `ssh vps`（`203.195.213.129`，2 核 2GiB 1Mbps，已接入腾讯云 EdgeOne 个人版，见根 [AGENTS.md](../AGENTS.md)） |
| 部署目录 / 二进制 | `/root/blog-backend/blog-backend` |
| systemd 服务 | `blog-backend`（`systemctl restart blog-backend` / `journalctl -u blog-backend -n 100 --no-pager`） |
| 交叉编译 | `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o blog-backend .` |

## 开发与提交规范速查

- **新增接口流程**：`models/` 定义结构体 → `internal/<domain>/` 实现 handler → `main.go` 注册路由 → `go build ./... && go test ./...` → 更新 `docs/api/<domain>.md` 与本文件 `AGENTS.md`。
- **Git 提交前缀**：`feat(backend agent):` / `fix(backend agent):` / `docs(backend agent):` / `refactor(backend agent):`。
- **更新日志**：`version` 无 `v` 前缀（如 `2.17.0`），`date` 为 `YYYY-MM-DD`，`logs` JSON 数组每条 ≤ 20 字、最多 5 条，写入 `pkg/db/changelog_seed.go`。

## 变更日志

- 2026-11-20：全面对齐 `main.go`、`config.go`、`pkg/db/db.go` 更新 `backend/AGENTS.md`（补齐 `backup`/`budget`/`download`/`growthdiary`/`material`/`perf`/`pitfall`/`renovation` 模块、路由与 6 张遗漏数据表，`IP2REGION_PATH`、`articles/random`、`articles/:id/view`、`pv/geo`、`admin/wechat/sync-logs`）；新增 `AGENTS.md` 强制随改动同步更新铁律；关联根 `AGENTS.md` 的 EdgeOne 接入与两层缓存架构。
- 2026-10-13：**图一乐二期（v2.17.0）** — 切换智谱 `glm-5.3-flash`、`query_almanac` 工具闭环、思考状态事件与 `tuyile_shares` 分享接口。
- 2026-10-12：**图一乐一期（v2.16.0）** — 新增 `almanac_days` 表、`cmd/almanac-crawler` 爬虫与 AG-UI SSE 运势聊天。
- 2026-09-16：存量图片 URL 域名统一为 `assets.fatwill.cn`；新增 `pkg/db/changelog_seed.go` 启动幂等播种。
- 2026-09-08：主域切换为 `fatwill.cn`，`CORS_ORIGIN` 支持逗号分隔多值，新增 `COS_LEGACY_DOMAINS` 与 `DOWNLOAD_ALLOWED_HOSTS`。
- 2026-06-29：新增微信公众号草稿同步与 `wechat_sync_logs` 表；05-13 新增微信 JS-SDK 签名接口；04-05 迁移 SQLite + 腾讯云 COS。
