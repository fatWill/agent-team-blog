package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatWill/agent-team-blog/backend/config"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// Init 初始化 SQLite 连接
func Init(cfg *config.DBConfig) error {
	// 确保数据库文件所在目录存在
	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建数据库目录失败: %w", err)
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return fmt.Errorf("连接 SQLite 失败: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取底层 sql.DB 失败: %w", err)
	}

	// SQLite 连接池设置（单文件数据库，并发写入受限）
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	// SQLite 性能优化 PRAGMA
	DB.Exec("PRAGMA journal_mode=WAL")
	DB.Exec("PRAGMA synchronous=NORMAL")
	DB.Exec("PRAGMA cache_size=-8000")
	DB.Exec("PRAGMA busy_timeout=5000")
	DB.Exec("PRAGMA foreign_keys=ON")

	// 自动建表（首次初始化空库时）
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("自动建表失败: %w", err)
	}

	return nil
}

// autoMigrate 创建所有表（如果不存在）
func autoMigrate() error {
	ddl := []string{
		// 文章表
		`CREATE TABLE IF NOT EXISTS articles (
			id TEXT NOT NULL PRIMARY KEY,
			title TEXT NOT NULL DEFAULT '',
			summary TEXT NOT NULL DEFAULT '',
			content TEXT,
			cover_image TEXT NOT NULL DEFAULT '',
			like_count INTEGER NOT NULL DEFAULT 0,
			views INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,

		// 文章点赞记录
		`CREATE TABLE IF NOT EXISTS article_likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			article_id TEXT NOT NULL DEFAULT '',
			device_id TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS uk_article_device ON article_likes (article_id, device_id)`,

		// 相册表
		`CREATE TABLE IF NOT EXISTS albums (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL DEFAULT '',
			description TEXT,
			cover_url TEXT,
			password_hash TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,

		// 照片表
		`CREATE TABLE IF NOT EXISTS photos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			album_id INTEGER NOT NULL DEFAULT 0,
			url TEXT NOT NULL DEFAULT '',
			caption TEXT,
			media_type TEXT NOT NULL DEFAULT 'image',
			thumbnail_url TEXT,
			duration INTEGER,
			password_hash TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_photos_album_id ON photos (album_id)`,

		// 照片点赞记录
		`CREATE TABLE IF NOT EXISTS photo_likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			photo_id INTEGER NOT NULL DEFAULT 0,
			device_id TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS uk_photo_likes_device ON photo_likes (photo_id, device_id)`,

		// 照片踩记录
		`CREATE TABLE IF NOT EXISTS photo_dislikes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			photo_id INTEGER NOT NULL DEFAULT 0,
			device_id TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS uk_photo_dislikes_device ON photo_dislikes (photo_id, device_id)`,

		// 个人资料表（单行）
		`CREATE TABLE IF NOT EXISTS profile (
			id INTEGER PRIMARY KEY,
			avatar TEXT NOT NULL DEFAULT '',
			bio TEXT NOT NULL DEFAULT ''
		)`,
		// 确保有一行默认数据
		`INSERT OR IGNORE INTO profile (id, avatar, bio) VALUES (1, '', '')`,

		// 留言表
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_id TEXT NOT NULL DEFAULT '',
			nickname TEXT,
			content TEXT NOT NULL DEFAULT '',
			last_modified_date TEXT,
			ip TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS uk_messages_device_id ON messages (device_id)`,

		// 更新日志表
		`CREATE TABLE IF NOT EXISTS changelogs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version TEXT NOT NULL DEFAULT '',
			date TEXT NOT NULL DEFAULT '',
			logs TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS uk_changelogs_version ON changelogs (version)`,

		// 页面访问记录表
		`CREATE TABLE IF NOT EXISTS page_views (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL DEFAULT '',
			device_id TEXT NOT NULL DEFAULT '',
			ip TEXT NOT NULL DEFAULT '',
			user_agent TEXT,
			device_type TEXT,
			browser TEXT,
			os TEXT,
			referer TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_pv_created_at ON page_views (created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_pv_path ON page_views (path)`,
		`CREATE INDEX IF NOT EXISTS idx_pv_device_id ON page_views (device_id)`,

		// 前端性能指标表
		`CREATE TABLE IF NOT EXISTS perf_metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			page TEXT NOT NULL DEFAULT '',
			fcp REAL,
			lcp REAL,
			ttfb REAL,
			dom_ready REAL,
			load_time REAL,
			js_load REAL,
			device_type TEXT,
			ua TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_perf_page ON perf_metrics (page)`,
		`CREATE INDEX IF NOT EXISTS idx_perf_created ON perf_metrics (created_at)`,

		// 兼容已有数据库：为 articles 表新增 views 字段（如果不存在）
		// SQLite 不支持 IF NOT EXISTS 语法，用 SELECT 检测
	}

	for _, sql := range ddl {
		if err := DB.Exec(sql).Error; err != nil {
			return fmt.Errorf("执行 DDL 失败 [%s]: %w", sql[:60], err)
		}
	}

	// 兼容迁移：为已有 articles 表新增 views 列
	var viewsColCount int
	DB.Raw(`SELECT COUNT(*) FROM pragma_table_info('articles') WHERE name = 'views'`).Scan(&viewsColCount)
	if viewsColCount == 0 {
		if err := DB.Exec(`ALTER TABLE articles ADD COLUMN views INTEGER NOT NULL DEFAULT 0`).Error; err != nil {
			return fmt.Errorf("迁移 articles.views 失败: %w", err)
		}
	}

	// 兼容迁移：为已有 photos 表新增 media_type、thumbnail_url、duration 列
	var mediaTypeColCount int
	DB.Raw(`SELECT COUNT(*) FROM pragma_table_info('photos') WHERE name = 'media_type'`).Scan(&mediaTypeColCount)
	if mediaTypeColCount == 0 {
		DB.Exec(`ALTER TABLE photos ADD COLUMN media_type TEXT NOT NULL DEFAULT 'image'`)
		DB.Exec(`ALTER TABLE photos ADD COLUMN thumbnail_url TEXT`)
		DB.Exec(`ALTER TABLE photos ADD COLUMN duration INTEGER`)
	}

	return nil
}
