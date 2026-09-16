#!/usr/bin/env bash
# ============================================================================
# migrate_assets_domain_20260916.sh
#
# 用途：将 SQLite 存量数据中的旧资源域名 assets.fatwill.cloud 全量替换为
#       assets.fatwill.cn（旧域名 EdgeOne 已失效，导致存量文章/相册裂图）。
#
# 替换范围（经 2026-09-16 全库扫描定稿，共 6 表 6 字段 / 31 行命中）：
#   albums.cover_url              1
#   articles.cover_image         15
#   articles.content              9
#   photos.url                    3
#   growth_diary_items.images     1
#   material_items.attachments    2
#
# 顺带处理（本次扫描命中数均为 0，保留语句以求一次性处理干净且幂等）：
#   cdn.fatwill.cloud -> cdn.fatwill.cn
#   pic.fatwill.cloud -> pic.fatwill.cn
#
# 明确不替换：
#   - 裸主域 fatwill.cloud（changelogs 文案、articles 站内链接、page_views.referer）
#   - photos.thumbnail_url / profile.avatar / growth_diary_items.videos 等
#     扫描命中为 0 的字段不纳入，避免无意义写放大
#
# 安全设计：
#   1. 前置校验 blog-backend 必须为 inactive（避免 SQLite database is locked）
#   2. 自动备份 blog.db（含 -wal / -shm，WAL 模式下缺一不可）
#   3. 单事务 BEGIN IMMEDIATE ... COMMIT，任一语句失败整体 ROLLBACK
#   4. 幂等：REPLACE 对已替换数据为 no-op；WHERE LIKE 保证第二次执行 0 变更
#   5. 末尾自校验：旧域名残留必须为 0，否则退出码非 0
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

# 替换映射：old|new
MAPPINGS=(
  "assets.fatwill.cloud|assets.fatwill.cn"
  "cdn.fatwill.cloud|cdn.fatwill.cn"
  "pic.fatwill.cloud|pic.fatwill.cn"
)

# 目标字段：table|column
TARGETS=(
  "albums|cover_url"
  "articles|cover_image"
  "articles|content"
  "photos|url"
  "growth_diary_items|images"
  "material_items|attachments"
)

log()  { printf '%b\n' "$*"; }
step() { printf '\n\033[1m==> %s\033[0m\n' "$*"; }
die()  { printf '\033[31m[FATAL] %s\033[0m\n' "$*" >&2; exit 1; }

sq() { sqlite3 -noheader "$DB_PATH" "$1"; }

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

declare -a BEFORE_OLD BEFORE_NEW LABELS
idx=0
for tgt in "${TARGETS[@]}"; do
  t="${tgt%%|*}"; c="${tgt##*|}"
  n_old=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%assets.fatwill.cloud%';")
  n_new=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%assets.fatwill.cn%';")
  LABELS[$idx]="$t.$c"
  BEFORE_OLD[$idx]="$n_old"
  BEFORE_NEW[$idx]="$n_new"
  printf '  %-34s old=%-4s new=%-4s\n' "$t.$c" "$n_old" "$n_new"
  idx=$((idx+1))
done

BARE_BEFORE=$(sq "SELECT (SELECT COUNT(*) FROM changelogs WHERE CAST(logs AS TEXT) LIKE '%fatwill.cloud%') + (SELECT COUNT(*) FROM page_views WHERE CAST(referer AS TEXT) LIKE '%fatwill.cloud%') + (SELECT COUNT(*) FROM articles WHERE CAST(content AS TEXT) LIKE '%//fatwill.cloud/%');")
log "  [守卫基线] 裸域 fatwill.cloud 引用行数（本次必须保持不变）：$BARE_BEFORE"

# ---------------------------------------------------------------------------
# Phase 3：事务内执行替换
# ---------------------------------------------------------------------------
step "Phase 3 执行替换（单事务）"

SQL_FILE="$(mktemp -t migrate_assets_domain.XXXXXX.sql)"
trap 'rm -f "$SQL_FILE"' EXIT

{
  echo ".bail on"
  echo "PRAGMA foreign_keys=ON;"
  echo "BEGIN IMMEDIATE;"
  for tgt in "${TARGETS[@]}"; do
    t="${tgt%%|*}"; c="${tgt##*|}"
    for m in "${MAPPINGS[@]}"; do
      old="${m%%|*}"; new="${m##*|}"
      echo "SELECT '  [$t.$c] $old -> $new';"
      # WHERE LIKE 限定命中行，保证幂等（第二次执行 changes()=0）
      echo "UPDATE \"$t\" SET \"$c\" = REPLACE(CAST(\"$c\" AS TEXT), '$old', '$new') WHERE CAST(\"$c\" AS TEXT) LIKE '%$old%';"
      echo "SELECT '      changed rows = ' || changes();"
    done
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
TOTAL_REPLACED=0
printf '  %-34s %-14s %-14s %s\n' "TABLE.COLUMN" "OLD before→after" "NEW before→after" "RESULT"
idx=0
for tgt in "${TARGETS[@]}"; do
  t="${tgt%%|*}"; c="${tgt##*|}"
  a_old=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%assets.fatwill.cloud%';")
  a_new=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%assets.fatwill.cn%';")
  res="OK"
  if [ "$a_old" != "0" ]; then res="FAIL(残留 $a_old)"; FAIL=1; fi
  printf '  %-34s %-14s %-14s %s\n' "$t.$c" "${BEFORE_OLD[$idx]}→$a_old" "${BEFORE_NEW[$idx]}→$a_new" "$res"
  TOTAL_REPLACED=$((TOTAL_REPLACED + BEFORE_OLD[idx]))
  idx=$((idx+1))
done

log ""
for p in cdn pic; do
  left=0
  for tgt in "${TARGETS[@]}"; do
    t="${tgt%%|*}"; c="${tgt##*|}"
    n=$(sq "SELECT COUNT(*) FROM \"$t\" WHERE CAST(\"$c\" AS TEXT) LIKE '%$p.fatwill.cloud%';")
    left=$((left+n))
  done
  log "  $p.fatwill.cloud 残留：$left"
  [ "$left" != "0" ] && FAIL=1
done

BARE_AFTER=$(sq "SELECT (SELECT COUNT(*) FROM changelogs WHERE CAST(logs AS TEXT) LIKE '%fatwill.cloud%') + (SELECT COUNT(*) FROM page_views WHERE CAST(referer AS TEXT) LIKE '%fatwill.cloud%') + (SELECT COUNT(*) FROM articles WHERE CAST(content AS TEXT) LIKE '%//fatwill.cloud/%');")
log "  [守卫校验] 裸域 fatwill.cloud 引用行数：$BARE_BEFORE → $BARE_AFTER"
if [ "$BARE_BEFORE" != "$BARE_AFTER" ]; then
  log "  \033[31m裸主域引用数发生变化，说明误伤了不该替换的内容！\033[0m"
  FAIL=1
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
log "  数据库        ：${DB_PATH}"
log "  备份文件      ：${BACKUP_PATH}"
log "  本次替换行数  ：${TOTAL_REPLACED}（重复执行时为 0，属正常幂等表现）"
log "  新域名总条数  ：$(sq "SELECT (SELECT COUNT(*) FROM albums WHERE CAST(cover_url AS TEXT) LIKE '%assets.fatwill.cn%') + (SELECT COUNT(*) FROM articles WHERE CAST(cover_image AS TEXT) LIKE '%assets.fatwill.cn%') + (SELECT COUNT(*) FROM articles WHERE CAST(content AS TEXT) LIKE '%assets.fatwill.cn%') + (SELECT COUNT(*) FROM photos WHERE CAST(url AS TEXT) LIKE '%assets.fatwill.cn%') + (SELECT COUNT(*) FROM growth_diary_items WHERE CAST(images AS TEXT) LIKE '%assets.fatwill.cn%') + (SELECT COUNT(*) FROM material_items WHERE CAST(attachments AS TEXT) LIKE '%assets.fatwill.cn%');")"

if [ "$FAIL" != "0" ]; then
  log ""
  die "校验未通过，请勿启动服务，先用备份回滚：cp -p $BACKUP_PATH $DB_PATH"
fi

log ""
log "  \033[32m全部校验通过。现在可以执行：systemctl start ${SERVICE_NAME}\033[0m"
exit 0
