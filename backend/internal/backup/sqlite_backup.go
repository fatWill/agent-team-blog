package backup

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatWill/agent-team-blog/backend/pkg/db"
)

const (
	// BackupDir 备份目录
	BackupDir = "/root/blog-backup"
	// MaxKeep 每个 trigger 类型最多保留文件数
	MaxKeep = 30
)

// BackupNow 立即执行一次 SQLite 备份
// trigger: startup / prebudget / daily
func BackupNow(trigger string) error {
	// 确保备份目录存在
	if err := os.MkdirAll(BackupDir, 0700); err != nil {
		log.Printf("❌ 创建备份目录失败: %v", err)
		return err
	}

	// 生成备份文件路径
	ts := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("blog-%s-%s.db", trigger, ts)
	destPath := filepath.Join(BackupDir, filename)

	// 使用 VACUUM INTO 生成一致性快照
	sql := fmt.Sprintf(`VACUUM INTO '%s'`, destPath)
	if err := db.DB.Exec(sql).Error; err != nil {
		log.Printf("❌ sqlite backup failed trigger=%s err=%v", trigger, err)
		return err
	}

	// 获取文件大小
	info, _ := os.Stat(destPath)
	size := int64(0)
	if info != nil {
		size = info.Size()
	}

	log.Printf("✅ sqlite backup created path=%s size=%d trigger=%s", destPath, size, trigger)

	// 清理同 trigger 类型的旧备份
	cleanOldBackups(trigger)

	return nil
}

// StartDailyBackup 启动每天 03:00 的定时备份 goroutine
func StartDailyBackup() {
	go func() {
		for {
			now := time.Now()
			// 计算下一个 03:00
			next := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			sleepDur := next.Sub(now)
			time.Sleep(sleepDur)

			if err := BackupNow("daily"); err != nil {
				log.Printf("⚠️ daily backup failed: %v", err)
			}
		}
	}()
	log.Println("✅ SQLite 每日备份调度已启动 (03:00)")
}

// cleanOldBackups 清理指定 trigger 类型超过 MaxKeep 数量的旧备份
func cleanOldBackups(trigger string) {
	prefix := fmt.Sprintf("blog-%s-", trigger)

	entries, err := os.ReadDir(BackupDir)
	if err != nil {
		return
	}

	// 筛选出同 trigger 的文件
	var matched []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), prefix) && strings.HasSuffix(e.Name(), ".db") {
			matched = append(matched, e.Name())
		}
	}

	if len(matched) <= MaxKeep {
		return
	}

	// 按文件名排序（时间戳在文件名中，字典序即时间序）
	sort.Strings(matched)

	// 删除最旧的
	toDelete := matched[:len(matched)-MaxKeep]
	for _, name := range toDelete {
		path := filepath.Join(BackupDir, name)
		if err := os.Remove(path); err != nil {
			log.Printf("⚠️ 删除旧备份失败: %s err=%v", path, err)
		} else {
			log.Printf("🗑️ 已清理旧备份: %s", path)
		}
	}
}
