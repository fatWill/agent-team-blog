#!/usr/bin/env bash
# ============================================================================
# migrate_assets_domain_20260916.sh
#
# 用途：一次停服窗口内，把 SQLite 存量数据中的旧域名图片 URL 全部收敛到 fatwill.cn。
#
# 分两组替换（同一个 BEGIN IMMEDIATE 事务内完成，只停服一次）：
#
#   【组 A】资源子域替换（原始需求）
#     assets.fatwill.cloud -> assets.fatwill.cn
#     cdn.fatwill.cloud    -> cdn.fatwill.cn     （扫描 0 命中，保留求幂等）
#     pic.fatwill.cloud    -> pic.fatwill.cn     （扫描 0 命中，保留求幂等）
#   作用于 6 表 6 字段（2026-09-16 全库扫描定稿，31 行 / 37 处命中）：
#     albums.cover_url             1 行 /  1 处
#     articles.cover_image        15 行 / 15 处
#     articles.content             9 行 / 11 处
#     photos.url                   3 行 /  3 处
#     growth_diary_items.images    1 行 /  4 处
#     material_items.attachments   2 行 /  3 处
#
#   【组 B】早期本地 uploads 路径替换（2026-09-16 扩展）
#     https://fatwill.cloud/uploads/ -> https://fatwill.cn/uploads/
#   仅作用于 articles 两个字段（全库扫描确认只有这两列命中，25 行 / 40 处）：
#     articles.cover_image        13 行 / 13 处
#     articles.content            12 行 / 27 处
#
#   合计：56 行 / 77 处，涉及 6 张表 6 个字段。
#
# ⚠️ 组 B 为什么必须用完整前缀 https://fatwill.cloud/uploads/ 而不是裸域：
#   裸主域 fatwill.cloud 全库共 156 处引用，其中 79 处是「从 fatwill.cloud 迁移
#   过来」这类 Markdown / Tiptap 正文文案、站内链接与 page_views.referer 埋点，
#   按约束一律不动。只有带 /uploads/ 前缀的 40 处才是真正的图片 URL。
#
# 明确不替换（Phase 4 有守卫强制校验）：
#   - 裸主域 fatwill.cloud（changelogs.logs 文案、articles.content 站内链接、
#     page_views.referer 埋点）
#   - 扫描命中为 0 的字段（photos.thumbnail_url / profile.avatar /
#     growth_diary_items.videos 等）不纳入 UPDATE，避免无意义写放大
#
# 安全设计：
#   1. 前置校验 blog-backend 必须为 inactive（避免 SQLite database is locked）
#   2. 变体形式前置断言（http:// / www. / 转义 \/ / 无协议头），防止漏替
#   3. 自动备份 blog.db（含 -wal / -shm，WAL 模式下缺一不可）
#   4. 单事务 BEGIN IMMEDIATE ... COMMIT，任一语句失败整体 ROLLBACK
#   5. 幂等：REPLACE 对已替换数据为 no-op；WHERE LIKE 保证第二次执行 0 变更
#   6. 裸域守卫用「出现次数级不变量」而非行数：
#        prose = occ(fatwill.cloud) - occ(assets.fatwill.cloud)
#                                  - occ(https://fatwill.cloud/uploads/)
#      该值在本次替换前后必须严格相等（当前基线 79 处）。
#      注意：旧版脚本用的行数守卫公式已失效——它把 uploads 行也算进去了，
#      而本次正是要替换这些行；且 page_views.referer 会随访问持续增长，
#      行数基线不具备可比性。出现次数级不变量对两者都免疫。
#   7. 末尾自校验：assets 残留 0、uploads 残留 0、prose 不变、integrity_check ok
#
# 用法：
#   bash migrate_assets_domain_20260916.sh                 # 默认生产库路径
#   bash migrate_assets_domain_20260916.sh /path/to/x.db    # 指定库（本地演练）
#   SKIP_SERVICE_CHECK=1 bash migrate_assets_domain_20260916.sh /tmp/t.db
#                                                          # 本地测试跳过服务检查
# ============================================================================
set -uo pipefail

DB_PATH="${1:-/root/blog-data/blog.db}"
SERVICE_NAME="${SERVICE_NAME:-blog-backend}"
BACKUP_SUFFIX="bak.migrate-assets-domain-20260916"
BACKUP_PATH="${DB_PATH}.${BACKUP_SUFFIX}"

# --- 组 A：资源子域映射 old|new ------------------------------------------------
MAPPINGS_DOMAIN=(
  "assets.fatwill.cloud|assets.fatwill.cn"
  "cdn.fatwill.cloud|cdn.fatwill.cn"
  "pic.fatwill.cloud|pic.fatwill.cn"
)

# --- 组 A：目标字段 table|column ----------------------------------------------
TARGETS_DOMAIN=(
  "albums|cover_url"
  "articles|cover_image"
  "articles|content"
  "photos|url"
  "growth_diary_items|images"
  "material_items|attachments"
)

# --- 组 B：uploads 前缀映射（精确前缀，绝不裸域） -------------------------------
UPLOADS_OLD="https://fatwill.cloud/uploads/"
UPLOADS_NEW="https://fatwill.cn/uploads/"

# --- 组 B：目标字段（全库扫描确认仅 articles 两列命中） ------------------------
TARGETS_UPLOADS=(
  "articles|cover_image"
  "articles|content"
)

# --- 裸域守卫扫描范围（全库唯一含 fatwill.cloud 的 8 列，经动态全表枚举确认） ---
BARE_COLS=(
  "albums|cover_url"
  "articles|cover_image"
  "articles|content"
  "changelogs|logs"
  "growth_diary_items|images"
  "material_items|attachments"
  "page_views|referer"
  "photos|url"
)

log()  { printf '%b\n' "$*"; }
step() { printf '\n\033[1m==> %s\033[0m\n' "$*"; }
die()  { printf '\033[31m[FATAL] %s\033[0m\n' "$*" >&2; exit 1; }

sq() { sqlite3 -noheader "$DB_PATH" "$1"; }

# occ_expr <column> <pattern> —— 生成「某列中 pattern 出现次数」的 SQL 片段
# 原理：(总长 - 抠掉 pattern 后的长度) / pattern 长度，SQLite 整除且必然整除
occ_expr() {
  local col="$1" pat="$2"
  printf "COALESCE(SUM((length(CAST(\"%s\" AS TEXT)) - length(REPLACE(CAST(\"%s\" AS TEXT),'%s',''))) / %s),0)" \
    "$col" "$col" "$pat" "${#pat}"
}

# 计算全库「纯文案裸域出现次数」= 裸域总数 - assets 前缀数 - uploads 前缀数
prose_occurrences() {
  local sql="SELECT 0"
  local tc t c
  for tc in "${BARE_COLS[@]}"; do
    t="${tc%%|*}"; c="${tc##*|}"
    sql="$sql + (SELECT $(occ_expr "$c" "fatwill.cloud") FROM \"$t\")"
    sql="$sql - (SELECT $(occ_expr "$c" "assets.fatwill.cloud") FROM \"$t\")"
    sql="$sql - (SELECT $(occ_expr "$c" "$UPLOADS_OLD") FROM \"$t\")"
  done
  sq "$sql;"
}

# 计算全库某 pattern 的总出现次数
total_occurrences() {
  local pat="$1" sql="SELECT 0" tc t c
  for tc in "${BARE_COLS[@]}"; do
    t="${tc%%|*}"; c="${tc##*|}"
    sql="$sql + (SELECT $(occ_expr "$c" "$pat") FROM \"$t\")"
  done
  sq "$sql;"
}

# ---------------------------------------------------------------------------
# Phase 0：前置校验
# ---------------------------------------------------------------------------
step "Phase 0 前置校验"

command -v sqlite3 >/dev/null 2>&1 || die "未找到 sqlite3 命令，请先安装：yum install -y sqlite"
[ -f "$DB_PATH" ] || die "数据库文件不存在：$DB_PATH"

if [ "${SKIP_SERVICE_CHECK:-0}" = "1" ]; then
  log "  [skip] 已通过 SKIP_SERVICE_CHECK=1 跳过服务状态检查（仅允许本地演练使用）"
else
  if ! command -v systemctl >/dev/null 2>&1; then
    die "未找到 systemctl，无法确认 ${SERVICE_NAME} 已停止，拒绝执行"
  fi
  SVC_STATE="$(systemctl is-active "$SERVICE_NAME" 2>/dev/null || true)"
  log "  ${SERVICE_NAME} 当前状态：${SVC_STATE:-unknown}"
  if [ "$SVC_STATE" != "inactive" ] && [ "$SVC_STATE" != "failed" ]; then
    die "${SERVICE_NAME} 未停止（当前 ${SVC_STATE:-unknown}）。请先执行：systemctl stop ${SERVICE_NAME}"
  fi
fi

# 确认库可写且当前无其他进程持锁
sq "PRAGMA quick_check;" >/dev/null 2>&1 || die "数据库无法读取或已损坏：$DB_PATH"
if ! sqlite3 "$DB_PATH" "BEGIN IMMEDIATE; ROLLBACK;" >/dev/null 2>&1; then
  die "数据库当前被其他进程持有写锁（database is locked），请确认后端已完全停止"
fi
log "  数据库可写、无写锁占用 ✓"

# WAL 落盘，保证备份文件自包含
sq "PRAGMA wal_checkpoint(TRUNCATE);" >/dev/null 2>&1 || true
log "  WAL 已 checkpoint ✓"

# --- 变体形式断言：确保 uploads URL 只有 https:// 一种写法，不会漏替 ---
log "  变体形式断言（期望全部为 0，否则说明存在脚本覆盖不到的写法）："
VARIANT_FAIL=0
for v in "http://fatwill.cloud/uploads/" "www.fatwill.cloud/uploads/" "fatwill.cloud\\/uploads"; do
  n=0
  for tc in "${BARE_COLS[@]}"; do
    t="${tc%%|*}"; c="${tc##*|}"
    k=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%$v%';")
    n=$((n + k))
  done
  printf '    %-34s %s\n' "$v" "$n"
  [ "$n" != "0" ] && VARIANT_FAIL=1
done
# 无协议头的 fatwill.cloud/uploads/（即不带 https:// 前缀的裸写法）
n=0
for tc in "${BARE_COLS[@]}"; do
  t="${tc%%|*}"; c="${tc##*|}"
  k=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%fatwill.cloud/uploads/%' AND CAST(\"$c\" AS TEXT) NOT LIKE '%${UPLOADS_OLD}%';")
  n=$((n + k))
done
printf '    %-34s %s\n' "无 https:// 前缀的 uploads" "$n"
[ "$n" != "0" ] && VARIANT_FAIL=1
if [ "$VARIANT_FAIL" != "0" ]; then
  die "检测到脚本未覆盖的 uploads URL 变体写法，请先扩充映射规则再执行，避免替换不彻底"
fi
log "  变体形式断言通过 ✓"

# ---------------------------------------------------------------------------
# Phase 1：强制备份
# ---------------------------------------------------------------------------
step "Phase 1 备份"

if [ -f "$BACKUP_PATH" ]; then
  TS="$(date +%Y%m%d%H%M%S)"
  log "  备份已存在（说明脚本曾执行过），另存为 ${BACKUP_PATH}.rerun-${TS}"
  cp -p "$DB_PATH" "${BACKUP_PATH}.rerun-${TS}" || die "备份失败"
  log "  备份完成：${BACKUP_PATH}.rerun-${TS}"
else
  cp -p "$DB_PATH" "$BACKUP_PATH" || die "备份失败"
  log "  备份完成：$BACKUP_PATH"
fi
# WAL / SHM 一并备份（checkpoint 后通常为空，仍保留以求完整）
[ -f "${DB_PATH}-wal" ] && cp -p "${DB_PATH}-wal" "${BACKUP_PATH}-wal" 2>/dev/null || true
[ -f "${DB_PATH}-shm" ] && cp -p "${DB_PATH}-shm" "${BACKUP_PATH}-shm" 2>/dev/null || true
ls -la "${DB_PATH}"* | sed 's/^/    /'

# ---------------------------------------------------------------------------
# Phase 2：替换前计数（before 基线）
# ---------------------------------------------------------------------------
step "Phase 2 替换前计数"

log "  【组 A】assets.fatwill.cloud"
A_LABELS=(); A_BEFORE_OLD=(); A_BEFORE_NEW=()
idx=0
for tgt in "${TARGETS_DOMAIN[@]}"; do
  t="${tgt%%|*}"; c="${tgt##*|}"
  n_old=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%assets.fatwill.cloud%';")
  n_new=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%assets.fatwill.cn%';")
  A_LABELS[$idx]="$t.$c"; A_BEFORE_OLD[$idx]="$n_old"; A_BEFORE_NEW[$idx]="$n_new"
  printf '    %-34s old=%-4s new=%-4s\n' "$t.$c" "$n_old" "$n_new"
  idx=$((idx+1))
done

log "  【组 B】${UPLOADS_OLD}"
B_LABELS=(); B_BEFORE_OLD=(); B_BEFORE_NEW=()
idx=0
for tgt in "${TARGETS_UPLOADS[@]}"; do
  t="${tgt%%|*}"; c="${tgt##*|}"
  n_old=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%${UPLOADS_OLD}%';")
  n_new=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%${UPLOADS_NEW}%';")
  B_LABELS[$idx]="$t.$c"; B_BEFORE_OLD[$idx]="$n_old"; B_BEFORE_NEW[$idx]="$n_new"
  printf '    %-34s old=%-4s new=%-4s\n' "$t.$c" "$n_old" "$n_new"
  idx=$((idx+1))
done

OCC_BARE_BEFORE=$(total_occurrences "fatwill.cloud")
OCC_ASSETS_BEFORE=$(total_occurrences "assets.fatwill.cloud")
OCC_UPLOADS_BEFORE=$(total_occurrences "$UPLOADS_OLD")
PROSE_BEFORE=$(prose_occurrences)

log ""
log "  [出现次数级基线]"
log "    裸域 fatwill.cloud 总出现次数        ：$OCC_BARE_BEFORE"
log "    其中 assets.fatwill.cloud            ：$OCC_ASSETS_BEFORE   （本次替换）"
log "    其中 ${UPLOADS_OLD}                  ：$OCC_UPLOADS_BEFORE   （本次替换）"
log "    \033[1m纯文案裸域（守卫不变量，必须保持不变）：$PROSE_BEFORE\033[0m"

# ---------------------------------------------------------------------------
# Phase 3：事务内执行替换（组 A + 组 B 同一事务，只停服一次）
# ---------------------------------------------------------------------------
step "Phase 3 执行替换（单事务，组 A + 组 B）"

SQL_FILE="$(mktemp -t migrate_assets_domain.XXXXXX.sql)"
trap 'rm -f "$SQL_FILE"' EXIT

{
  echo ".bail on"
  echo "PRAGMA foreign_keys=ON;"
  echo "BEGIN IMMEDIATE;"

  echo "SELECT '--- 组 A：资源子域 ---';"
  for tgt in "${TARGETS_DOMAIN[@]}"; do
    t="${tgt%%|*}"; c="${tgt##*|}"
    for m in "${MAPPINGS_DOMAIN[@]}"; do
      old="${m%%|*}"; new="${m##*|}"
      echo "SELECT '  [$t.$c] $old -> $new';"
      # WHERE LIKE 限定命中行，保证幂等（第二次执行 changes()=0）
      echo "UPDATE \"$t\" SET \"$c\" = REPLACE(CAST(\"$c\" AS TEXT), '$old', '$new') WHERE CAST(\"$c\" AS TEXT) LIKE '%$old%';"
      echo "SELECT '      changed rows = ' || changes();"
    done
  done

  echo "SELECT '--- 组 B：uploads 前缀 ---';"
  for tgt in "${TARGETS_UPLOADS[@]}"; do
    t="${tgt%%|*}"; c="${tgt##*|}"
    echo "SELECT '  [$t.$c] ${UPLOADS_OLD} -> ${UPLOADS_NEW}';"
    echo "UPDATE \"$t\" SET \"$c\" = REPLACE(CAST(\"$c\" AS TEXT), '${UPLOADS_OLD}', '${UPLOADS_NEW}') WHERE CAST(\"$c\" AS TEXT) LIKE '%${UPLOADS_OLD}%';"
    echo "SELECT '      changed rows = ' || changes();"
  done

  echo "COMMIT;"
} > "$SQL_FILE"

# .bail on：任一语句报错立即中止，未 COMMIT 的事务由 sqlite3 退出时自动 ROLLBACK
if ! sqlite3 "$DB_PATH" < "$SQL_FILE"; then
  log ""
  die "替换过程出错，事务已自动 ROLLBACK，数据库未被修改。可用备份恢复：$BACKUP_PATH"
fi

log "  事务已提交 ✓"
sq "PRAGMA wal_checkpoint(TRUNCATE);" >/dev/null 2>&1 || true

# ---------------------------------------------------------------------------
# Phase 4：验证
# ---------------------------------------------------------------------------
step "Phase 4 验证"

FAIL=0
A_REPLACED=0
B_REPLACED=0

log "  【组 A】assets.fatwill.cloud 残留应为 0"
printf '    %-34s %-16s %-16s %s\n' "TABLE.COLUMN" "OLD before→after" "NEW before→after" "RESULT"
idx=0
for tgt in "${TARGETS_DOMAIN[@]}"; do
  t="${tgt%%|*}"; c="${tgt##*|}"
  a_old=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%assets.fatwill.cloud%';")
  a_new=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%assets.fatwill.cn%';")
  res="OK"
  if [ "$a_old" != "0" ]; then res="FAIL(残留 $a_old)"; FAIL=1; fi
  printf '    %-34s %-16s %-16s %s\n' "$t.$c" "${A_BEFORE_OLD[$idx]}→$a_old" "${A_BEFORE_NEW[$idx]}→$a_new" "$res"
  A_REPLACED=$((A_REPLACED + A_BEFORE_OLD[idx]))
  idx=$((idx+1))
done

log ""
log "  【组 B】${UPLOADS_OLD} 残留应为 0（articles.cover_image / articles.content 二次扫描）"
printf '    %-34s %-16s %-16s %s\n' "TABLE.COLUMN" "OLD before→after" "NEW before→after" "RESULT"
idx=0
for tgt in "${TARGETS_UPLOADS[@]}"; do
  t="${tgt%%|*}"; c="${tgt##*|}"
  b_old=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%${UPLOADS_OLD}%';")
  b_new=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%${UPLOADS_NEW}%';")
  res="OK"
  if [ "$b_old" != "0" ]; then res="FAIL(残留 $b_old)"; FAIL=1; fi
  printf '    %-34s %-16s %-16s %s\n' "$t.$c" "${B_BEFORE_OLD[$idx]}→$b_old" "${B_BEFORE_NEW[$idx]}→$b_new" "$res"
  B_REPLACED=$((B_REPLACED + B_BEFORE_OLD[idx]))
  idx=$((idx+1))
done

# 全库兜底：任意列都不应再有 fatwill.cloud/uploads/
log ""
left_uploads=0
for tc in "${BARE_COLS[@]}"; do
  t="${tc%%|*}"; c="${tc##*|}"
  n=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%fatwill.cloud/uploads/%';")
  left_uploads=$((left_uploads + n))
done
log "  全库 fatwill.cloud/uploads/ 残留（含任意协议头写法）：$left_uploads"
[ "$left_uploads" != "0" ] && FAIL=1

for p in cdn pic; do
  left=0
  for tgt in "${TARGETS_DOMAIN[@]}"; do
    t="${tgt%%|*}"; c="${tgt##*|}"
    n=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%$p.fatwill.cloud%';")
    left=$((left+n))
  done
  log "  $p.fatwill.cloud 残留：$left"
  [ "$left" != "0" ] && FAIL=1
done

# --- 裸域守卫：出现次数级不变量 ---
OCC_BARE_AFTER=$(total_occurrences "fatwill.cloud")
OCC_ASSETS_AFTER=$(total_occurrences "assets.fatwill.cloud")
OCC_UPLOADS_AFTER=$(total_occurrences "$UPLOADS_OLD")
PROSE_AFTER=$(prose_occurrences)

log ""
log "  [裸域守卫校验]"
log "    裸域总出现次数    ：${OCC_BARE_BEFORE} → ${OCC_BARE_AFTER}（预期减少 $((OCC_ASSETS_BEFORE + OCC_UPLOADS_BEFORE))）"
log "    assets 前缀       ：${OCC_ASSETS_BEFORE} → ${OCC_ASSETS_AFTER}"
log "    uploads 前缀      ：${OCC_UPLOADS_BEFORE} → ${OCC_UPLOADS_AFTER}"
log "    纯文案裸域(不变量)：${PROSE_BEFORE} → ${PROSE_AFTER}"

if [ "$PROSE_BEFORE" != "$PROSE_AFTER" ]; then
  log "  \033[31m纯文案裸域出现次数发生变化，说明误伤了不该替换的文案内容！\033[0m"
  FAIL=1
fi
EXPECT_BARE_AFTER=$((OCC_BARE_BEFORE - OCC_ASSETS_BEFORE - OCC_UPLOADS_BEFORE))
if [ "$OCC_BARE_AFTER" != "$EXPECT_BARE_AFTER" ]; then
  log "  \033[31m裸域总出现次数不等于预期值 ${EXPECT_BARE_AFTER}，替换范围与预期不符！\033[0m"
  FAIL=1
fi

# --- JSON 完整性：Tiptap content 与 JSON 数组字段替换后必须仍是合法 JSON ---
log ""
BAD_JSON=$(sq "SELECT (SELECT COUNT(*) FROM articles WHERE CAST(content AS TEXT) LIKE '{%' AND json_valid(CAST(content AS TEXT)) = 0) + (SELECT COUNT(*) FROM growth_diary_items WHERE CAST(images AS TEXT) LIKE '[%' AND json_valid(CAST(images AS TEXT)) = 0) + (SELECT COUNT(*) FROM material_items WHERE CAST(attachments AS TEXT) LIKE '[%' AND json_valid(CAST(attachments AS TEXT)) = 0);" 2>/dev/null || echo "skip")
if [ "$BAD_JSON" = "skip" ]; then
  log "  JSON 校验：当前 sqlite3 未编译 JSON1 扩展，跳过"
else
  log "  JSON 非法条数：$BAD_JSON"
  [ "$BAD_JSON" != "0" ] && FAIL=1
fi

if ! sq "PRAGMA integrity_check;" | grep -q '^ok$'; then
  log "  integrity_check 未通过"
  FAIL=1
else
  log "  integrity_check ✓"
fi

# ---------------------------------------------------------------------------
# Phase 5：摘要
# ---------------------------------------------------------------------------
step "Phase 5 执行摘要"
log "  数据库            ：${DB_PATH}"
log "  备份文件          ：${BACKUP_PATH}"
log "  组 A 替换行数     ：${A_REPLACED}（assets/cdn/pic 子域，6 表 6 字段）"
log "  组 B 替换行数     ：${B_REPLACED}（uploads 前缀，articles 2 字段）"
log "  本次替换行数合计  ：$((A_REPLACED + B_REPLACED))（重复执行时为 0，属正常幂等表现）"
log "  裸主域残留        ：${OCC_BARE_AFTER} 处（均为文案 / 站内链接 / referer 埋点，按约束保留）"
log "  新域名 assets 条数：$(sq "SELECT (SELECT COUNT(*) FROM albums WHERE CAST(cover_url AS TEXT) LIKE '%assets.fatwill.cn%') + (SELECT COUNT(*) FROM articles WHERE CAST(cover_image AS TEXT) LIKE '%assets.fatwill.cn%') + (SELECT COUNT(*) FROM articles WHERE CAST(content AS TEXT) LIKE '%assets.fatwill.cn%') + (SELECT COUNT(*) FROM photos WHERE CAST(url AS TEXT) LIKE '%assets.fatwill.cn%') + (SELECT COUNT(*) FROM growth_diary_items WHERE CAST(images AS TEXT) LIKE '%assets.fatwill.cn%') + (SELECT COUNT(*) FROM material_items WHERE CAST(attachments AS TEXT) LIKE '%assets.fatwill.cn%');")"
log "  新域名 uploads 条数：$(sq "SELECT (SELECT COUNT(*) FROM articles WHERE CAST(cover_image AS TEXT) LIKE '%${UPLOADS_NEW}%') + (SELECT COUNT(*) FROM articles WHERE CAST(content AS TEXT) LIKE '%${UPLOADS_NEW}%');")"

if [ "$FAIL" != "0" ]; then
  log ""
  die "校验未通过，请勿启动服务，先用备份回滚：cp -p $BACKUP_PATH $DB_PATH"
fi

log ""
log "  \033[32m全部校验通过。现在可以执行：systemctl start ${SERVICE_NAME}\033[0m"
exit 0
