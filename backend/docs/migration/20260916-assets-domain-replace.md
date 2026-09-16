# 存量图片 URL 域名迁移（assets.fatwill.cloud → assets.fatwill.cn）

- **执行日期**：2026-09-16
- **脚本**：`backend/scripts/migrate_assets_domain_20260916.sh`
- **目标库**：`/root/blog-data/blog.db`（SQLite，WAL 模式）
- **背景**：主域已切换到 `fatwill.cn`，但 DB 存量图片 URL 仍指向 `assets.fatwill.cloud`，该域名走 EdgeOne CNAME，切换后已无 DNS 解析（服务器实测 `assets.fatwill.cloud 无解析`），导致存量文章/相册裂图。

---

## 一、全库扫描结果

扫描方式：从生产库 scp 一份 3.3MB 副本到本地，`PRAGMA wal_checkpoint(TRUNCATE)` 合并 WAL 后，动态枚举 `sqlite_master` 全部 17 张业务表的每一列（SQLite 动态类型，一律 `CAST(col AS TEXT)` 参与 `LIKE`），逐列统计命中。未采用「只查已知字段」的方式，避免漏字段。

### 1.1 `assets.fatwill.cloud` 命中清单（本次替换范围）

| 表名 | 字段名 | 命中行数 | 命中样例 |
|---|---|---|---|
| `albums` | `cover_url` | 1 | `https://assets.fatwill.cloud/upload/1775486161007-9828ba4f.jpg` |
| `articles` | `cover_image` | 15 | `https://assets.fatwill.cloud/blog-images/20260411/hermes-cover-16x9.png` |
| `articles` | `content` | 9 | `https://assets.fatwill.cloud/blog-images/20260411/block-org-compare.png`（Tiptap JSON 内 `image.attrs.src`） |
| `photos` | `url` | 3 | `https://assets.fatwill.cloud/upload/1775486946980-44c94592.mp4` |
| `growth_diary_items` | `images` | 1 | `https://assets.fatwill.cloud/upload/1785739289522-0b8b5873.jpg`（JSON 数组） |
| `material_items` | `attachments` | 2 | `https://assets.fatwill.cloud/materials/1776610270969-5b4c76ee.jpg`（JSON 数组） |
| **合计** | | **31** | |

### 1.2 其他旧域名扫描结果

| 域名 | 全库命中 | 处置 |
|---|---|---|
| `assets.fatwill.cloud` | 31 行 | ✅ 替换为 `assets.fatwill.cn` |
| `cdn.fatwill.cloud` | 0 | 脚本保留语句（幂等无副作用），本次 0 变更 |
| `pic.fatwill.cloud` | 0 | 同上 |
| `img.fatwill.cloud` / `static.fatwill.cloud` | 0 | 不纳入 |
| `fatwill.cloud`（裸主域） | 115 处引用 | ❌ **不替换**（按约束） |
| `www.fatwill.cloud` / `api.fatwill.cloud` | 各 2 处（均在 `page_views.referer`） | ❌ 不替换 |

`http://` 协议头的旧域名 URL 命中数为 0，全部为 `https://`。

### 1.3 扫描为 0、故未纳入的字符串字段

`photos.thumbnail_url`、`profile.avatar`、`growth_diary_items.videos`、`renovation_articles.content`、`material_items.tags`、`pitfall_items.*`、`budget_items.remark`、`messages.content`、`perf_metrics.*`、`wechat_sync_logs.*` — 均为 0 命中，不纳入 UPDATE，避免无意义写放大（2 核 2G 机器上减少一次全表扫描与页重写）。

---

## 二、⚠️ 阻断性发现：`assets.fatwill.cn` 当前返回 403

**替换脚本本身已验证通过，但仅执行替换不足以修复裂图。** 实测：

```
GET https://assets.fatwill.cn/blog-images/20260411/hermes-cover-16x9.png
→ 403 AccessDenied  (COS XML 错误体)

# 服务器侧绕过 nginx 直连 COS 源站
GET https://fatwill-cloud-1253664788.cos.ap-guangzhou.myqcloud.com/upload/...
→ 403 AccessDenied
```

**根因定位**：403 来自 COS 存储桶自身权限，不是 nginx 反代配置问题。证据链：

1. `assets.fatwill.cn` DNS 是 A 记录直连本机，nginx 有对应 `server` 块反代到 COS 源站，链路通（`assets_access.log` 有 403 记录，说明请求到达并被回源）；
2. 在服务器本地直连 COS 源站域名同样 403，绕过了 nginx；
3. 响应头含 `x-cos-request-id`，错误体是 COS 的 `<Code>AccessDenied</Code>`。

推断：旧的 `assets.fatwill.cloud` 走 EdgeOne，由 CDN 用回源鉴权读取私有桶；`assets.fatwill.cn` 改成 nginx 直接反代后是匿名 GET，而桶的读权限并非「公有读」，因此 403。

**需要你在腾讯云控制台确认并处理（二选一）**：

- 方案 A：将桶 `fatwill-cloud-1253664788` 的访问权限设为「公有读私有写」（最省事，与原 CDN 对外可读的效果一致）；
- 方案 B：保持私有桶，在 nginx 反代层加 COS 签名（需要在 nginx 侧引入 SecretId/SecretKey 做 Authorization 计算，2 核 2G 上有额外 CPU 开销，且密钥落到 nginx 配置，不推荐）。

**建议顺序**：先修 COS 权限并 `curl` 验证 `assets.fatwill.cn` 返回 200，再执行本次 DB 替换。若先替换，URL 会从「域名不解析」变成「403」，裂图现象不变。

> 附带发现（本次不处理，仅告知）：`articles` 里另有 25 行图片 URL 是 `https://fatwill.cloud/uploads/...` 形式（`cover_image` 13 行 + `content` 12 行）。这些是更早期的本地 uploads 路径，走裸主域。实测 `fatwill.cloud/uploads/...` → 301 → `fatwill.cn/uploads/...` → **200 正常**，所以不裂图，只是多一跳。按你的约束裸主域不替换，这批保持原样。

---

## 三、交付给运维的执行说明

### 3.1 前置条件

- [ ] COS 桶权限已修复，`curl -I https://assets.fatwill.cn/upload/1775486161007-9828ba4f.jpg` 返回 **200**（未满足时执行替换无法解决裂图）
- [ ] 脚本已上传到服务器，例如 `/root/migrate_assets_domain_20260916.sh`
- [ ] 服务器已安装 `sqlite3` CLI（`yum install -y sqlite`；后端用的是纯 Go 驱动，机器上不一定有 CLI）

### 3.2 部署顺序

```bash
# 1) 停后端（必须，否则 SQLite 写锁冲突 database is locked）
systemctl stop blog-backend
systemctl is-active blog-backend          # 期望输出 inactive

# 2) 备份 + 替换 + 校验（脚本内部自动完成三件事，无需手工 cp）
bash /root/migrate_assets_domain_20260916.sh
echo "exit=$?"                            # 期望 0；非 0 时不要启动服务

# 3) 校验（脚本已自动跑，如需手工复核）
sqlite3 /root/blog-data/blog.db \
  "SELECT COUNT(*) FROM articles WHERE CAST(cover_image AS TEXT) LIKE '%assets.fatwill.cloud%';"   # 期望 0
sqlite3 /root/blog-data/blog.db \
  "SELECT COUNT(*) FROM articles WHERE CAST(cover_image AS TEXT) LIKE '%assets.fatwill.cn%';"      # 期望 15

# 4) 启动后端
systemctl start blog-backend
systemctl is-active blog-backend          # 期望 active

# 5) 接口抽查
curl -s "https://fatwill.cn/api/articles?page=1&page_size=5" | head -c 500
```

脚本会自动执行：服务状态校验（非 inactive 直接 `exit 1`）→ 写锁探测 → WAL checkpoint → 备份到 `blog.db.bak.migrate-assets-domain-20260916` → 单事务替换 → 残留校验 + 裸域守卫校验 + `integrity_check` → 打印摘要。

### 3.3 抽样验证的具体 URL

替换后逐条 `curl`，全部应返回 `200` 且 `Content-Type: image/*`：

| 用途 | 文章/资源 | 验证 URL |
|---|---|---|
| 文章封面（blog-images 路径） | `d6be77db-…` Hermes | `https://assets.fatwill.cn/blog-images/20260411/hermes-cover-16x9.png` |
| 文章正文内嵌图（Tiptap JSON） | `b302ea9f-…` Claude Code Token | `https://assets.fatwill.cn/blog-images/20260414/claude-code-token-infographic.png` |
| 相册封面 + 照片（upload 路径） | album `id=2` / photo `id=15` | `https://assets.fatwill.cn/upload/1775486161007-9828ba4f.jpg` |
| 材料清单附件（materials 路径） | `material_items id=1` | `https://assets.fatwill.cn/materials/1776610270969-5b4c76ee.jpg` |

页面级抽查（浏览器打开，确认封面与正文图片均正常）：

```
https://fatwill.cn/articles/d6be77db-ec3b-4bbc-b099-3187e87bc5e3
https://fatwill.cn/articles/b302ea9f-54d3-450c-9656-9347127539fa
https://fatwill.cn/life/album           # 相册封面
```

> 注意：前端 `toThumbUrl` / `toWebpUrl` 会给 `assets` 域名的图片追加 `?imageMogr2/...` 数据万象参数。验证时请连带参数一起 curl 一次，确认数据万象在新域名下也正常（数据万象与桶绑定，通常不受域名影响，但值得确认一次）。

### 3.4 回滚

```bash
systemctl stop blog-backend
cp -p /root/blog-data/blog.db.bak.migrate-assets-domain-20260916 /root/blog-data/blog.db
rm -f /root/blog-data/blog.db-wal /root/blog-data/blog.db-shm   # 清理残留 WAL，避免与旧库文件不匹配
systemctl start blog-backend
```

本地已实测回滚：恢复后 `assets.fatwill.cloud` 命中数回到 31，`integrity_check = ok`。

### 3.5 预期停服时间

| 阶段 | 耗时 |
|---|---|
| `systemctl stop` | ~1s |
| 备份 3.3MB 库文件 | <1s |
| 6 条 UPDATE（31 行命中，最大表 `page_views` 未参与写） | <1s |
| 校验查询（含 `integrity_check`） | 1~2s |
| `systemctl start` + 预热 | 2~3s |
| **合计** | **约 5~10 秒** |

建议选低峰时段执行。停服期间前端 SSR 取不到 API 数据会报错，但时间窗极短。

---

## 四、脚本自测记录（本地生产库副本，非生产环境）

| 用例 | 结果 |
|---|---|
| 首次执行 | ✅ 31 行替换，6 字段 old 全部 15/9/3/1/1/2 → 0，new 0 → 对应值，`exit 0` |
| 第二次执行（幂等） | ✅ 所有 `changes()` = 0，无报错，`exit 0` |
| 服务未停 / 无 systemctl | ✅ `exit 1`，未触碰数据库 |
| 数据库文件不存在 | ✅ `exit 1` |
| 裸主域守卫 | ✅ `fatwill.cloud` 引用行数 89 → 89，未被误伤 |
| JSON 完整性 | ✅ `articles.content` / `changelogs.logs` / `growth_diary_items.images` / `material_items.attachments` 的 `json_valid()` 非法条数均为 0 |
| `integrity_check` | ✅ ok |
| 回滚 | ✅ 从备份恢复后命中数回到 31，integrity ok |

**未在生产服务器上执行任何 sqlite3 写操作**；对生产仅做了只读操作（`ls`、`systemctl is-active`、日志 `tail`、`curl` 探测）与一次 `scp` 只读下载。本次不涉及任何 Go 代码变更。

---

## 五、设计说明

**为什么拉副本到本地而不是让运维回传扫描结果**：生产库仅 3.3MB，`scp` 一份是纯读操作，零风险。拿到真实数据后才能（1）动态枚举全部表列而非依赖 `docs/sqlite.md`（实测文档已滞后，缺 `perf_metrics`、`material_items`、`growth_diary_items`、`budget_items`、`pitfall_items`、`renovation_articles` 六张表）；（2）在真实数据上跑通首次执行 + 幂等 + 回滚三条路径。让运维回传结果需要多轮往返，且无法验证脚本。

**为什么幂等靠 `WHERE ... LIKE` 而不是只靠 `REPLACE`**：`REPLACE` 本身对已替换数据是 no-op，但不带 `WHERE` 会重写全表所有行（`updated_at` 若有触发器会被带动，且 SQLite 会重写数据页）。加 `WHERE LIKE` 后第二次执行 `changes()` 精确为 0，既是幂等证明，也避免 31 行的操作退化成全表写。

**为什么加「裸主域守卫」**：本次最大风险是把 `fatwill.cloud` 裸域一起换掉。`REPLACE(x,'assets.fatwill.cloud','assets.fatwill.cn')` 在字符串层面不会匹配裸域，但为防脚本被后人误改扩大范围，Phase 2 记录裸域引用基线、Phase 4 比对，不一致即判失败。

**为什么用 `BEGIN IMMEDIATE` 而非 `BEGIN`**：`BEGIN IMMEDIATE` 立刻获取写锁，若后端仍在运行会当场失败而不是等到第一条 UPDATE 才报 `database is locked`，失败更早、更明确。Phase 0 也用 `BEGIN IMMEDIATE; ROLLBACK;` 做了一次无副作用的锁探测。

**为什么备份要连带 `-wal` / `-shm`**：WAL 模式下已提交但未 checkpoint 的事务还在 `-wal` 里，只 `cp` 主库文件可能丢数据。脚本先 `wal_checkpoint(TRUNCATE)` 把 WAL 落盘再备份，并额外拷一份 `-wal` / `-shm` 兜底。
