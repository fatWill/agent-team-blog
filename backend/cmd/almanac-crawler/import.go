package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/fatWill/agent-team-blog/backend/internal/almanac"
	"github.com/fatWill/agent-team-blog/backend/models"
)

// importFrom 从本地生成的 SQLite 只读导入 almanac_days 到目标库（按 date upsert，单事务）。
//
// 用途：本地跑完爬虫后把 almanac.db 上传到服务器，再导入生产库，
// 生产侧无需联网抓取，也不需要安装 sqlite3 CLI。
func importFrom(srcPath string, dst *almanac.Repo, start, end string) error {
	abs, err := filepath.Abs(srcPath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(abs); err != nil {
		return fmt.Errorf("源库不存在: %w", err)
	}
	src, err := gorm.Open(sqlite.Open("file:"+abs+"?mode=ro"), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
	if err != nil {
		return fmt.Errorf("打开源库失败: %w", err)
	}
	var rows []models.AlmanacDay
	if err := src.Where("date >= ? AND date <= ?", start, end).Order("date ASC").Find(&rows).Error; err != nil {
		return fmt.Errorf("读取源库失败: %w", err)
	}
	if len(rows) == 0 {
		return fmt.Errorf("源库 %s 在 %s~%s 内没有数据", abs, start, end)
	}
	n, err := dst.UpsertBatch(rows)
	if err != nil {
		return err
	}
	log.Printf("✅ 从 %s 导入 %d 天（%s ~ %s）", abs, n, rows[0].Date, rows[len(rows)-1].Date)
	return nil
}
