-- =============================================================================
-- changelog 2.15.0 补数脚本（备用方案 / Plan B）
-- =============================================================================
-- 正常情况下【不需要执行本脚本】。
-- 主方案：新版二进制启动时会由 pkg/db/changelog_seed.go 自动幂等播种 2.15.0。
--
-- 仅在以下情况使用本脚本：
--   1. 暂不发布新二进制，但需要让 2.15.0 立刻出现在 changelog 页面；
--   2. 新二进制已上线，但页面仍查不到 2.15.0，需人工兜底核查。
--
-- 目标库：SQLite（不是 MySQL），生产路径 /root/blog-data/blog.db
-- 表结构对齐：models/misc.go → models.Changelog，DDL 见 pkg/db/db.go
--   id INTEGER PK AUTOINCREMENT / version TEXT (uk_changelogs_version 唯一索引)
--   date TEXT / logs TEXT(JSON 数组) / created_at DATETIME / updated_at DATETIME
--
-- ⚠️ 执行前务必先备份，且强烈建议先停服务再执行：
--    数据库为 SQLite WAL 模式且服务端连接池 MaxOpenConns=1，
--    服务运行中由外部进程写入可能触发 database is locked。
--
--    systemctl stop blog-backend
--    cp /root/blog-data/blog.db /root/blog-data/blog.db.bak.$(date +%F-%H%M%S)
--    sqlite3 /root/blog-data/blog.db < seed_changelog_2.15.0.sql
--    systemctl start blog-backend
-- =============================================================================

-- INSERT OR IGNORE：依赖 version 唯一索引保证幂等，重复执行不会产生重复行，
-- 也不会覆盖线上已人工修订过的同版本文案。
INSERT OR IGNORE INTO changelogs (version, date, logs, created_at, updated_at)
VALUES (
    '2.15.0',
    '2026-09-16',
    '["🌐 站点主域切换到 fatwill.cn","🔁 旧域名 fatwill.cloud 自动 301 跳转到新域","🖼️ 图片与静态资源域名同步切换为 assets.fatwill.cn / cdn.fatwill.cn"]',
    '2026-09-16 00:00:00',
    '2026-09-16 00:00:00'
);

-- 校验：应返回 1 行，且 logs 为合法 JSON 数组
SELECT id, version, date, logs FROM changelogs WHERE version = '2.15.0';
