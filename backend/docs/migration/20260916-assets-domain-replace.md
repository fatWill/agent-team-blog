# 存量图片 URL 域名迁移（旧域名 → fatwill.cn）

- **执行日期**：2026-09-16
- **脚本**：`backend/scripts/migrate_assets_domain_20260916.sh`
- **目标库**：`/root/blog-data/blog.db`（SQLite，WAL 模式）
- **替换范围**：**56 行 / 77 处**，6 张表 6 个字段
  - 组 A `assets.fatwill.cloud` → `assets.fatwill.cn`：31 行 / 37 处
  - 组 B `https://fatwill.cloud/uploads/` → `https://fatwill.cn/uploads/`：25 行 / 40 处
- **背景**：主域已切换到 `fatwill.cn`。存量图片 URL 分两批遗留：
  1. `assets.fatwill.cloud` 走 EdgeOne CNAME，切换后已无 DNS 解析（服务器实测 `assets.fatwill.cloud 无解析`），导致存量文章/相册**裂图**；
  2. 更早期的 `https://fatwill.cloud/uploads/...` 走裸主域，目前靠 301 跳转到 `fatwill.cn/uploads/...` 仍能出图，**不裂图但多一跳**。按「不留历史尾巴」的要求一并清掉。

---

## 〇、组 B 为什么必须用完整前缀而不是裸主域

`fatwill.cloud` 裸主域全库共 **156 处**引用，其中只有 40 处是 `/uploads/` 图片 URL，其余 **79 处是不能动的内容**：

| 类别 | 出现次数 | 示例 |
|---|---|---|
| `page_views.referer` 访问埋点 | 75 | 历史访问来源，属事实记录，改了就是伪造数据 |
| `articles.content` 正文文案 / 站内链接 | 3 | `href="https://fatwill.cloud/articles/b302ea9f-…"`、`href="https://fatwill.cloud"` |
| `changelogs.logs` 更新日志文案 | 1 | 「从 fatwill.cloud 迁移」这类说明文字 |

所以替换串必须是 `https://fatwill.cloud/uploads/` 这个**完整精确前缀**——它在字符串层面无法匹配上表任何一条，天然不误伤。脚本另有出现次数级守卫做强制兜底（见 §5）。

---


## 一、全库扫描结果

扫描方式：从生产库 scp 一份 3.3MB 副本到本地，`PRAGMA wal_checkpoint(TRUNCATE)` 合并 WAL 后，动态枚举 `sqlite_master` 全部 17 张业务表的每一列（SQLite 动态类型，一律 `CAST(col AS TEXT)` 参与 `LIKE`），逐列统计命中。未采用「只查已知字段」的方式，避免漏字段。

### 1.1 组 A：`assets.fatwill.cloud` 命中清单

| 表名 | 字段名 | 命中行数 | 出现次数 | 命中样例 |
|---|---|---|---|---|
| `albums` | `cover_url` | 1 | 1 | `https://assets.fatwill.cloud/upload/1775486161007-9828ba4f.jpg` |
| `articles` | `cover_image` | 15 | 15 | `https://assets.fatwill.cloud/blog-images/20260411/hermes-cover-16x9.png` |
| `articles` | `content` | 9 | 11 | `https://assets.fatwill.cloud/blog-images/20260411/block-org-compare.png`（Tiptap JSON 内 `image.attrs.src`） |
| `photos` | `url` | 3 | 3 | `https://assets.fatwill.cloud/upload/1775486946980-44c94592.mp4` |
| `growth_diary_items` | `images` | 1 | 4 | `https://assets.fatwill.cloud/upload/1785739289522-0b8b5873.jpg`（JSON 数组，单行多图） |
| `material_items` | `attachments` | 2 | 3 | `https://assets.fatwill.cloud/materials/1776610270969-5b4c76ee.jpg`（JSON 数组） |
| **合计** | | **31** | **37** | |

### 1.2 组 B：`https://fatwill.cloud/uploads/` 命中清单（本轮新增）

全库动态枚举确认：**仅 `articles` 两列命中**，其他 15 张表 0 命中。

| 表名 | 字段名 | 命中行数 | 出现次数 | 命中样例 |
|---|---|---|---|---|
| `articles` | `cover_image` | 13 | 13 | `https://fatwill.cloud/uploads/blog-images/20260415/hero.png` |
| `articles` | `content` | 12 | 27 | `https://fatwill.cloud/uploads/blog-images/20260718-finch/hero.png`（Tiptap JSON，单篇多图） |
| **合计** | | **25** | **40** | |

变体写法排查（全部 0 命中，故单条映射即可覆盖，脚本 Phase 0 有强制断言）：

| 变体 | 命中 |
|---|---|
| `http://fatwill.cloud/uploads/`（明文 http） | 0 |
| `www.fatwill.cloud/uploads/` | 0 |
| `fatwill.cloud\/uploads`（JSON 转义斜杠） | 0 |
| `fatwill.cloud/uploads/` 但无 `https://` 前缀 | 0 |

组 A 与组 B 的命中行**无交集**（`cover_image` / `content` 上两者同时命中的行数均为 0），所以 31 + 25 = 56 行不存在重复计数。

### 1.3 其他旧域名扫描结果

| 域名 | 全库命中 | 处置 |
|---|---|---|
| `assets.fatwill.cloud` | 31 行 / 37 处 | ✅ 替换为 `assets.fatwill.cn` |
| `https://fatwill.cloud/uploads/` | 25 行 / 40 处 | ✅ 替换为 `https://fatwill.cn/uploads/` |
| `cdn.fatwill.cloud` | 0 | 脚本保留语句（幂等无副作用），本次 0 变更 |
| `pic.fatwill.cloud` | 0 | 同上 |
| `img.fatwill.cloud` / `static.fatwill.cloud` | 0 | 不纳入 |
| `fatwill.cloud`（裸主域，纯文案部分） | **79 处** | ❌ **不替换**（按约束，见 §0） |
| `www.fatwill.cloud` / `api.fatwill.cloud` | 各 2 处（均在 `page_views.referer`） | ❌ 不替换 |

`http://` 协议头的旧域名 URL 命中数为 0，全部为 `https://`。

> **⚠️ 基线口径修正**：上一轮文档记录的「裸主域 115 处」与守卫基线「89」均已失效，原因有两点：
> 1. 上一轮统计的是**行数**，本轮改为**出现次数**（单行 JSON 可含多个 URL，例如 `growth_diary_items.images` 1 行含 4 处），出现次数才是能精确校验的口径；
> 2. `page_views.referer` 会随访问持续增长，行数基线天生不可比。
>
> 本轮实测（2026-09-16 23:29 快照）：裸域总计 **156 处** = assets 37 + uploads 40 + 纯文案 **79**。替换后应剩 79 处，全部为文案 / 站内链接 / referer 埋点。

---

### 1.4 扫描为 0、故未纳入的字符串字段

`photos.thumbnail_url`、`profile.avatar`、`growth_diary_items.videos`、`renovation_articles.content`、`material_items.tags`、`pitfall_items.*`、`budget_items.remark`、`messages.content`、`perf_metrics.*`、`wechat_sync_logs.*` — 均为 0 命中，不纳入 UPDATE，避免无意义写放大（2 核 2G 机器上减少一次全表扫描与页重写）。

---

## 二、⚠️ 阻断性发现：`assets.fatwill.cn` 当前返回 403（仅影响组 A）

**组 A 的替换脚本本身已验证通过，但仅执行替换不足以修复 assets 类裂图。** 实测：

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

**建议顺序**：先修 COS 权限并 `curl` 验证 `assets.fatwill.cn` 返回 200，再执行本次 DB 替换。若先替换，assets 类 URL 会从「域名不解析」变成「403」，裂图现象不变。

> **组 B（`uploads/` 路径）不受此阻断影响**：`fatwill.cn/uploads/` 由本机 nginx 直接提供本地静态文件，与 COS 无关。已对全部 13 个去重目标 URL 逐条 `curl` 实测，**13/13 返回 200 `image/png`**。所以即使 COS 权限暂未修复，组 B 的替换也是纯收益（省掉一次 301 跳转）。
>
> 换句话说：本脚本可以随时执行，组 B 立即生效；组 A 的可见效果取决于 COS 权限何时修好。

---

## 三、交付给运维的执行说明

### 3.1 前置条件

- [ ] （仅影响组 A 效果）COS 桶权限已修复，`curl -I https://assets.fatwill.cn/upload/1775486161007-9828ba4f.jpg` 返回 **200**；未修复时组 A 替换仍可执行，只是 assets 图片仍 403
- [ ] 脚本已上传到服务器，例如 `/root/migrate_assets_domain_20260916.sh`
- [ ] 服务器已安装 `sqlite3` CLI（已确认在位：`/usr/bin/sqlite3` 3.26.0；后端用的是纯 Go 驱动，机器上不一定有 CLI）

### 3.2 部署顺序

```bash
# 1) 停后端（必须，否则 SQLite 写锁冲突 database is locked）
systemctl stop blog-backend
systemctl is-active blog-backend          # 期望输出 inactive

# 2) 备份 + 替换 + 校验（脚本内部自动完成三件事，无需手工 cp）
bash /root/migrate_assets_domain_20260916.sh
echo "exit=$?"                            # 期望 0；非 0 时不要启动服务

# 3) 校验（脚本已自动跑，如需手工复核）
# 组 A
sqlite3 /root/blog-data/blog.db \
  "SELECT COUNT(*) FROM articles WHERE CAST(cover_image AS TEXT) LIKE '%assets.fatwill.cloud%';"   # 期望 0
sqlite3 /root/blog-data/blog.db \
  "SELECT COUNT(*) FROM articles WHERE CAST(cover_image AS TEXT) LIKE '%assets.fatwill.cn%';"      # 期望 15
# 组 B
sqlite3 /root/blog-data/blog.db \
  "SELECT COUNT(*) FROM articles WHERE CAST(cover_image AS TEXT) LIKE '%fatwill.cloud/uploads/%';" # 期望 0
sqlite3 /root/blog-data/blog.db \
  "SELECT COUNT(*) FROM articles WHERE CAST(content AS TEXT) LIKE '%fatwill.cloud/uploads/%';"     # 期望 0
sqlite3 /root/blog-data/blog.db \
  "SELECT COUNT(*) FROM articles WHERE CAST(cover_image AS TEXT) LIKE '%fatwill.cn/uploads/%';"    # 期望 13
# 裸域残留（出现次数级，期望 79，全部为文案 / 站内链接 / referer）
sqlite3 /root/blog-data/blog.db \
  "SELECT (SELECT COALESCE(SUM((length(CAST(referer AS TEXT))-length(REPLACE(CAST(referer AS TEXT),'fatwill.cloud','')))/13),0) FROM page_views)
        + (SELECT COALESCE(SUM((length(CAST(content AS TEXT))-length(REPLACE(CAST(content AS TEXT),'fatwill.cloud','')))/13),0) FROM articles)
        + (SELECT COALESCE(SUM((length(CAST(logs AS TEXT))-length(REPLACE(CAST(logs AS TEXT),'fatwill.cloud','')))/13),0) FROM changelogs);"

# 4) 启动后端
systemctl start blog-backend
systemctl is-active blog-backend          # 期望 active

# 5) 接口抽查
curl -s "https://fatwill.cn/api/articles?page=1&page_size=5" | head -c 500
```

脚本会自动执行：服务状态校验（非 inactive 直接 `exit 1`）→ 写锁探测 → **uploads 变体写法断言** → WAL checkpoint → 备份到 `blog.db.bak.migrate-assets-domain-20260916` → **单事务替换（组 A + 组 B）** → 组 A 残留校验 + **组 B 二次扫描** + 裸域出现次数守卫 + JSON 合法性 + `integrity_check` → 打印摘要。

### 3.3 抽样验证的具体 URL

**组 A（assets 子域）** —— COS 权限修好后应返回 `200` + `Content-Type: image/*`；当前实测 403：

| 用途 | 文章/资源 | 验证 URL |
|---|---|---|
| 文章封面（blog-images 路径） | `d6be77db-…` Hermes | `https://assets.fatwill.cn/blog-images/20260411/hermes-cover-16x9.png` |
| 文章正文内嵌图（Tiptap JSON） | `b302ea9f-…` Claude Code Token | `https://assets.fatwill.cn/blog-images/20260414/claude-code-token-infographic.png` |
| 相册封面 + 照片（upload 路径） | album `id=2` / photo `id=15` | `https://assets.fatwill.cn/upload/1775486161007-9828ba4f.jpg` |
| 材料清单附件（materials 路径） | `material_items id=1` | `https://assets.fatwill.cn/materials/1776610270969-5b4c76ee.jpg` |

**组 B（uploads 本地静态路径）** —— 已全量实测 **13/13 返回 200 `image/png`**：

| 用途 | 文章 id | 验证 URL |
|---|---|---|
| 最早期封面（2025 年目录） | `02b3ff00-…` | `https://fatwill.cn/uploads/blog-images/20250715/hero.png` |
| 封面（扁平日期目录） | `b3c9acb3-…` | `https://fatwill.cn/uploads/blog-images/20260415/hero.png` |
| 封面（带 slug 后缀目录） | `e5736e0b-…` Finch | `https://fatwill.cn/uploads/blog-images/20260718-finch/hero.png` |
| 封面（带 slug 后缀目录） | `922f62ff-…` Kimi K3 | `https://fatwill.cn/uploads/blog-images/20260719-kimi-k3/hero.png` |
| 封面 + 正文内嵌图 | `07fb937d-…` Hermes Agent | `https://fatwill.cn/uploads/blog-images/20260830-hermes-agent/hero.png` |
| 封面 + 正文内嵌图 | `1d246d4d-…` Agentic RAG | `https://fatwill.cn/uploads/blog-images/20260831-agentic-rag/hero.png` |
| 封面（spec-kit 特殊命名） | `97b2752a-…` | `https://fatwill.cn/uploads/blog-images/20260506/spec-kit-hero.png` |

组 B 命中的 13 篇文章 id（`cover_image`，与 `content` 命中的 12 篇高度重合）：

```
02b3ff00  0398a0b9  07fb937d  1d246d4d  32e943f9  7f9070f7  83d2a091
922f62ff  97b2752a  b3c9acb3  cba7b80a  dc4cafff  e5736e0b
```

页面级抽查（浏览器打开，确认封面与正文图片均正常）：

```
# 组 A（assets 子域，依赖 COS 权限修复）
https://fatwill.cn/articles/d6be77db-ec3b-4bbc-b099-3187e87bc5e3
https://fatwill.cn/articles/b302ea9f-54d3-450c-9656-9347127539fa
https://fatwill.cn/life/album           # 相册封面

# 组 B（uploads 路径，本轮新增，不依赖 COS）
https://fatwill.cn/articles/07fb937d-e976-49d9-87f1-1a9678bd7c85   # Hermes Agent，封面+正文图
https://fatwill.cn/articles/1d246d4d-12b2-47bf-815b-e6a500a3b856   # Agentic RAG，封面+正文图
https://fatwill.cn/articles/e5736e0b-73bd-4451-9bf1-f1f8bdc3fa7c   # Finch，封面+正文图
https://fatwill.cn/articles/02b3ff00-5452-46c3-9b35-03ddb089ba94   # 2025 最早期目录
https://fatwill.cn/                                                # 首页列表页封面（13 篇混排）
```

> 注意：前端 `toThumbUrl` / `toWebpUrl` 会给 `assets` 域名的图片追加 `?imageMogr2/...` 数据万象参数。验证时请连带参数一起 curl 一次，确认数据万象在新域名下也正常（数据万象与桶绑定，通常不受域名影响，但值得确认一次）。`uploads/` 路径是本地静态文件，不走数据万象，无此问题。

### 3.4 回滚

```bash
systemctl stop blog-backend
cp -p /root/blog-data/blog.db.bak.migrate-assets-domain-20260916 /root/blog-data/blog.db
rm -f /root/blog-data/blog.db-wal /root/blog-data/blog.db-shm   # 清理残留 WAL，避免与旧库文件不匹配
systemctl start blog-backend
```

本地已实测回滚：恢复后 `assets.fatwill.cloud` 命中数回到 31、`fatwill.cloud/uploads/` 回到 25，`integrity_check = ok`，且 `.dump` 与原始库**逐字节一致**。

### 3.5 预期停服时间

| 阶段 | 耗时 |
|---|---|
| `systemctl stop` | ~1s |
| 备份 3.3MB 库文件 | <1s |
| Phase 0 变体断言（4 组 × 8 列 COUNT） | <1s |
| 20 条 UPDATE（组 A 18 条 + 组 B 2 条，56 行命中，最大表 `page_views` 未参与写） | <1s |
| 校验查询（含出现次数守卫、JSON 校验、`integrity_check`） | 2~3s |
| `systemctl start` + 预热 | 2~3s |
| **合计** | **约 6~12 秒** |

组 A 与组 B 合并在**同一个事务**内，只需一次停服窗口。建议选低峰时段执行。停服期间前端 SSR 取不到 API 数据会报错，但时间窗极短。

---

## 四、脚本自测记录（本地生产库副本 + 生产环境干跑，均非生产库）

### 4.1 本地三条路径（macOS，sqlite 3.51）

| 用例 | 结果 |
|---|---|
| **首次执行** | ✅ 共 **56 行**替换（组 A 31 + 组 B 25），`exit 0`；组 A 六字段 old 15/9/3/1/1/2 → 0；组 B `cover_image` 13→0、`content` 12→0 |
| **第二次执行（幂等）** | ✅ 全部 20 条 UPDATE 的 `changes()` 均为 0，替换行数合计 0，无报错，`exit 0` |
| **回滚** | ✅ 从备份恢复后命中数回到 31 / 25，`integrity_check = ok`，`.dump` 与原始库逐字节一致 |
| 服务未停 / 无 systemctl | ✅ `exit 1`，未触碰数据库 |
| 数据库文件不存在 | ✅ `exit 1` |
| **裸域守卫（正向）** | ✅ 纯文案裸域出现次数 **79 → 79** 不变；裸域总数 156 → 79，等于预期值 |
| **裸域守卫（反向验证）** | ✅ 故意把映射改成裸域 `fatwill.cloud→fatwill.cn` 后，守卫检出 79→76 并报「误伤了不该替换的文案内容」，`exit 1`，证明守卫不是空跑 |
| 变体写法断言 | ✅ `http://` / `www.` / 转义 `\/` / 无协议头 四类均为 0，断言通过 |
| JSON 完整性 | ✅ `articles.content` / `growth_diary_items.images` / `material_items.attachments` 的 `json_valid()` 非法条数为 0 |
| `integrity_check` | ✅ ok |

### 4.2 生产环境干跑（OpenCloudOS，bash 4.4，sqlite **3.26.0**）

本地 sqlite 是 3.51，生产是 3.26，版本差 25 个小版本，故额外做了兼容性验证：

| 验证项 | 结果 |
|---|---|
| `occ()` 出现次数算式在 3.26 下的正确性 | ✅ 合成测试库验证：裸域 4 / assets 1 / uploads 2，与预期完全一致；NULL 与整数列不报错 |
| `mktemp -t xxx.XXXXXX.sql` GNU 版兼容 | ✅ 正常生成 `/tmp/migrate_assets_domain.2xFHYz.sql`（BSD/GNU 行为差异已排除） |
| `json_valid()` 可用性 | ✅ 3.26 已带 JSON1，返回 1（脚本对不可用场景也有 skip 兜底） |
| **完整脚本干跑**（对生产库 `.backup` 出的一致性快照执行） | ✅ 输出与本地**完全一致**：56 行替换、守卫 79→79、`integrity_check` ok、`exit 0` |
| 干跑后生产库状态复核 | ✅ 生产库仍为 `assets.cloud=15 / uploads.cloud=13`（未被改动），`blog-backend` 仍 `active`，临时文件已清理 |

**未在生产服务器上执行任何 sqlite3 写操作**；对生产库仅做了只读操作（`ls`、`systemctl is-active`、`SELECT`、`.backup` 只读快照）与一次 `scp` 只读下载。干跑的写操作全部发生在 `/tmp` 下的快照副本上，且事后已删除。本次不涉及任何 Go 代码变更。

---

## 五、设计说明

**为什么扩展原脚本而不是新建一个独立脚本**：两组替换的强制约束完全相同（停服、备份、WAL checkpoint、单事务、幂等、校验），而且都作用在 `articles.cover_image` / `articles.content` 这两列上。拆成两个脚本会导致：停服两次（或一次停服跑两个脚本、备份文件互相覆盖）、`articles` 表被重写两遍、两份守卫基线各自失效风险翻倍。合并成「单脚本 + 组 A/组 B 两组映射 + 同一个 `BEGIN IMMEDIATE` 事务」后，要么全成功要么全回滚，只有一个备份点和一个回滚指令，运维心智负担最低。组 A/组 B 在代码里用独立的 `MAPPINGS_DOMAIN` / `UPLOADS_OLD` 常量和独立的计数数组分开统计，日志里也分节输出，可读性没有损失。

**为什么守卫口径从「行数」改成「出现次数」**（本轮最关键的修正）：上一轮的守卫公式是

```sql
COUNT(*) FROM changelogs WHERE logs LIKE '%fatwill.cloud%'
+ COUNT(*) FROM page_views WHERE referer LIKE '%fatwill.cloud%'
+ COUNT(*) FROM articles WHERE content LIKE '%//fatwill.cloud/%'
```

它在本轮**必然失效**，有两个独立的原因：

1. 第三项 `content LIKE '%//fatwill.cloud/%'` 会把组 B 要替换的 uploads 行也算进来（实测该项当前为 13）。本轮正是要替换这些行，守卫会误判为「被误伤」而拦截正常执行；
2. `page_views.referer` 随访问量持续增长，行数基线天生不可比。上一轮记的 115 / 89，本轮实测已变成 156 处 / 79 处（`referer` 命中从当时的值涨到 75 行）。**任何以行数为基线的守卫都会随时间自然漂移**。

改用 `occ(s) = (length(x) - length(replace(x,s,''))) / length(s)` 计算出现次数后，守卫变成一个真正的**不变量**：

```
prose = occ(fatwill.cloud) - occ(assets.fatwill.cloud) - occ(https://fatwill.cloud/uploads/)
```

替换 `assets.fatwill.cloud → assets.fatwill.cn` 会让 `occ(fatwill.cloud)` 和 `occ(assets.fatwill.cloud)` **同步减 1**，uploads 同理，所以 `prose` 在正确替换下恒定不变（实测 79 → 79）；一旦有人误改成裸域替换，`occ(fatwill.cloud)` 单独下降而两个减项不动，`prose` 立刻变小被检出（反向验证实测 79 → 76，`exit 1`）。这个口径对 `page_views` 增长也免疫——新增的 referer 同时进入被减数，不影响等式成立。脚本另外还断言了「裸域总数减少量恰好等于 assets + uploads 出现次数」，做双重约束。

**为什么组 B 必须用 `https://fatwill.cloud/uploads/` 完整前缀**：见 §0。补充一点工程细节——脚本 Phase 0 增加了**变体写法断言**，主动扫 `http://`、`www.`、JSON 转义 `\/`、无协议头四种可能漏替的写法，全为 0 才继续。这是防「精确前缀」反过来变成漏替风险：前缀越精确，越怕数据里存在另一种写法。断言把这个风险从「祈祷」变成「验证」。

**为什么幂等靠 `WHERE ... LIKE` 而不是只靠 `REPLACE`**：`REPLACE` 本身对已替换数据是 no-op，但不带 `WHERE` 会重写全表所有行（`updated_at` 若有触发器会被带动，且 SQLite 会重写数据页）。加 `WHERE LIKE` 后第二次执行 `changes()` 精确为 0，既是幂等证明，也避免 56 行的操作退化成全表写。在 2 核 2G 机器上，`articles` 全表重写与 25 行定点更新的 IO 差异是量级级别的。

**为什么用 `BEGIN IMMEDIATE` 而非 `BEGIN`**：`BEGIN IMMEDIATE` 立刻获取写锁，若后端仍在运行会当场失败而不是等到第一条 UPDATE 才报 `database is locked`，失败更早、更明确。Phase 0 也用 `BEGIN IMMEDIATE; ROLLBACK;` 做了一次无副作用的锁探测。

**为什么备份要连带 `-wal` / `-shm`**：WAL 模式下已提交但未 checkpoint 的事务还在 `-wal` 里，只 `cp` 主库文件可能丢数据。脚本先 `wal_checkpoint(TRUNCATE)` 把 WAL 落盘再备份，并额外拷一份 `-wal` / `-shm` 兜底。本轮下载副本时生产 `-wal` 有 4MB（比主库还大），印证了这个设计的必要性。

**为什么要在生产环境干跑一遍**：生产 sqlite 是 **3.26.0**（2018 年），本地是 3.51。`occ()` 算式依赖整数除法与 `length()` 对多字节字符的行为，`mktemp -t` 在 BSD/GNU 下模板语义还不同。这些都不是「应该没问题」能糊过去的，所以用 `.backup` 从生产库做了一份一致性快照放到 `/tmp`，在真实的 bash 4.4 + sqlite 3.26 上把脚本完整跑通，确认输出与本地逐项一致后才交付，事后删除快照。全过程对生产库只有读。

**为什么拉副本到本地而不是让运维回传扫描结果**：生产库仅 3.3MB，`scp` 一份是纯读操作，零风险。拿到真实数据后才能（1）动态枚举全部表列而非依赖 `docs/sqlite.md`（实测文档已滞后，缺 `perf_metrics`、`material_items`、`growth_diary_items`、`budget_items`、`pitfall_items`、`renovation_articles` 六张表）；（2）在真实数据上跑通首次执行 + 幂等 + 回滚三条路径。让运维回传结果需要多轮往返，且无法验证脚本。
