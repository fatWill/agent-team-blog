# AGENTS.md（agent-team-blog 仓库总纲）

个人博客项目 [https://fatwill.cn](https://fatwill.cn) 的仓库级 AI 协作与运维总纲。前后端细分索引见：

| 目录 | 规范文档 | 说明 |
| --- | --- | --- |
| `frontend/` | [frontend/AGENTS.md](frontend/AGENTS.md) | Nuxt 3 + Vue 3 + TypeScript + Tailwind CSS + Pinia 前端 SSR |
| `backend/` | [backend/AGENTS.md](backend/AGENTS.md) | Go + Gin + GORM + SQLite + Redis 后端 API |
| `blogs/` | — | 博客文章源文件与配图素材 |
| `scripts/` | — | 缓存清理脚本（`purge-blog-cache.sh`、`purge-nginx-cache.sh`） |

---

## ⚠️ 强制规范：AGENTS.md 必须随代码与配置同步更新

> **任何涉及项目结构、功能、配置或线上架构的变更，必须在同一次改动中同步更新对应的 `AGENTS.md`，禁止只改代码/线上配置而不更新文档。**

触发同步更新的场景（包括但不限于）：
- **根 `AGENTS.md`**：域名解析、腾讯云 EdgeOne 规则/证书/套餐、VPS Nginx 路由、缓存策略、部署流程或跨模块约定变更。
- **`frontend/AGENTS.md`**：新增/删除/重命名页面路由、Feature 模块、组件、Composable、Plugin、Middleware、Nitro API 代理、环境变量或核心依赖。
- **`backend/AGENTS.md`**：新增/修改 HTTP API 路由、领域模块（`internal/`）、基础设施包（`pkg/`）、SQLite 数据表或字段迁移、环境变量或部署配置。

更新原则：与原条目冲突时**直接修改原条目**，不堆砌过时描述；并在对应文档的「变更日志」顶部补录一行简要记录。

### Git 提交铁律
- **凡在本仓库完成代码或文档改动并验证通过后，直接 `git commit && git push`，无需再询问用户**（只暂存本次改动涉及的文件，不误提交工作区中其他无关未提交文件）。

---

## 线上架构与基础设施

### 1. 服务器概况
- **主机**：`ssh vps`（`203.195.213.129`，OpenCloudOS 8，2 核 2GiB 内存 / 实际可用 1.7GiB + 2G Swap，**1Mbps 源站带宽**，50G 磁盘）。
- **同机服务**：
  - `blog-frontend`（systemd，Node 22 + Nuxt `.output/server/index.mjs`，监听 `127.0.0.1:3000`）
  - `blog-backend`（systemd，`/root/blog-backend/blog-backend`，监听 `:8080`，SQLite `/root/blog-data/blog.db`，自动备份至 `/root/blog-backup/`）
  - `redis`（`127.0.0.1:6379`）
  - `nginx`（监听 80 / 443，主站点配置 `/etc/nginx/conf.d/blog.conf`）
  - 同机亦托管微信小游戏后端（如 `wegame-taikongtiaodong`，监听 `127.0.0.1:8091`，反代 `/wegame/taikongtiaodong/`），**1Mbps 带宽共享，必须保护源站带宽**。

### 2. 域名、DNS 与腾讯云 EdgeOne 架构
- **DNS 解析**：`fatwill.cn` 托管在腾讯云 DNSPod（`greg.dnspod.net` / `fox.dnspod.net`），已通过 **DNSPod 托管 / CNAME 接入腾讯云 EdgeOne（个人版）**。
- **EdgeOne 个人版配额**：1 个站点（`fatwill.cn`）、最多 200 个子域名、**50 GB/月流量**、**300 万次/月请求数**（超额转后付费，需关注用量告警）、1 条自定义速率限制规则。

| 加速域名 | EdgeOne CNAME | 源站与回源配置 | HTTPS | EdgeOne 规则引擎 / 用途 |
| --- | --- | --- | --- | --- |
| `fatwill.cn` | `fatwill.cn.eo.dnse2.com` | `203.195.213.129`，协议跟随（80/443），回源 Host `fatwill.cn` | EO 免费 DV 证书（已部署） | 博客主站。规则 **`blog-cache`**：`IF1`（上方）`HOST = fatwill.cn` → 节点缓存 TTL 强制缓存（5~10 分钟，挡住访客与爬虫、释放 VPS 1Mbps 带宽）；`IF2`（下方覆盖上方）`HOST = fatwill.cn` 且 `URL Path` 匹配 `/api/`、`/admin`、`/login`、`/wegame/` → 节点缓存 TTL **不缓存** |
| `cdn.fatwill.cn` | `cdn.fatwill.cn.eo.dnse2.com` | `203.195.213.129`（VPS Nginx） | EO 免费 DV 证书（已部署） | Nuxt 构建产物（`/_nuxt/`、`public/` 等带 Hash 静态资源，长缓存） |
| `assets.fatwill.cn` | `assets.fatwill.cn.eo.dnse2.com` | 腾讯云 COS `fatwill-cloud-1253664788`（直连 COS，**不经 VPS**） | EO 免费 DV 证书（已部署） | 博客图片、相册与静态素材 |
| `game.fatwill.cn` | `game.fatwill.cn.eo.dnse2.com` | `203.195.213.129`，HTTPS 443，**回源 Host `fatwill.cn`**（复用 VPS 现有证书与 server 块） | EO 免费 DV 证书（已部署） | 微信小游戏专用加速域（API / WebSocket）。规则 **`wegame-ws-nocache`**：`HOST = game.fatwill.cn` → 开启 **WebSocket（超时 300s）** + 节点缓存 TTL **不缓存**；速率限制规则 **`wegame`**：`HOST = game.fatwill.cn`，单客户端 IP 10 秒 > 60 次 → 阻断拦截 60 秒 |

### 3. EdgeOne 规则引擎关键注意事项（防踩坑）
- **执行优先级**：规则引擎**从上到下逐条执行，下方规则覆盖上方规则**；因此通用缓存规则必须放上面（`IF1`），不缓存例外规则放下面（`IF2`）。
- **节点缓存 TTL 匹配条件限制**：仅支持 `HOST`、`URL Path`、`URL Full`、文件后缀等请求特征，**不支持 `请求方法` 或 `Cookie`**（EdgeOne 默认只缓存 GET/HEAD，POST/PUT/DELETE 天然不缓存）。
- **WebSocket 限制**：EdgeOne 仅支持基于 **HTTP/1.1** 的 WebSocket，最大空闲超时 **300 秒**，客户端需定时发送心跳包保活。

### 4. 两层缓存与刷新流程
博客公开 HTML 页面当前有两层缓存：
1. **边缘层（EdgeOne）**：规则 `blog-cache` 短缓存（5~10 分钟），默认挡住绝大部分访客与爬虫。
2. **源站层（VPS Nginx）**：`/etc/nginx/conf.d/blog.conf` 中 `location /` 的 24 小时 SSR 缓存（`/var/cache/nginx/blog`）。

**发版或发布/修改文章后的清缓存方式**：
- 清理 VPS Nginx 缓存：运行 `/usr/local/bin/purge-blog-cache.sh`（源码见 [scripts/purge-blog-cache.sh](scripts/purge-blog-cache.sh)）；后端在文章/相册等写入时也会通过 `pkg/nginx/purge.go` 触发本地清理。
- 边缘层 EdgeOne：不操作时最多等待 5~10 分钟自动过期；如需秒生效，到 EdgeOne 控制台 **站点加速 → 缓存管理 → 缓存刷新** 刷新对应 URL 或全站。

### 5. 运维与安全红线
- 远程运维统一使用 `remote-server`、`nginx-config`（及按需使用 `frontend-deploy`、`go-deploy`、`sftp-upload`、`remote-sql`）；修改 Nginx 前必须备份，经 `nginx -t` 通过后再 reload。
- 查线上 SQLite 必须只读（`sqlite3 -readonly`），禁止对 WAL 运行中的数据库直接写入。
- **绝不读取、打印或转述** `/etc/systemd/system/blog-backend.service` 中的环境变量密钥，或 `~/.zshrc` 中的数据库/Redis 密码。
- 重启服务、修改配置等中危操作前说明影响；删除数据、修改防火墙、系统升级等高危操作必须先获得用户确认。

---

## 变更日志
- **2026-11-20**：新建仓库根 `AGENTS.md`；记录 `fatwill.cn`（博客主站 `blog-cache`）与 `game.fatwill.cn`（微信小游戏 `wegame-ws-nocache` + `wegame` 速率限制）接入腾讯云 EdgeOne 个人版的完整配置与两层缓存说明；确立前后端与根目录 `AGENTS.md` 强制随改动同步更新的规范。
