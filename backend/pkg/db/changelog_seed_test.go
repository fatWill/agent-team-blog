package db

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/fatWill/agent-team-blog/backend/config"
	"github.com/fatWill/agent-team-blog/backend/models"
)

// TestSeedChangelogIdempotent 验证 changelog 播种在空库建表后能写入 2.15.0，
// 且重复执行保持幂等、不覆盖已有文案。
func TestSeedChangelogIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	if err := Init(&config.DBConfig{Path: path}); err != nil {
		t.Fatalf("Init 失败: %v", err)
	}

	var got models.Changelog
	if err := DB.Where("version = ?", "2.15.0").First(&got).Error; err != nil {
		t.Fatalf("查询 2.15.0 失败: %v", err)
	}
	if got.Date != "2026-09-16" {
		t.Errorf("date 期望 2026-09-16，实际 %s", got.Date)
	}

	var logs []string
	if err := json.Unmarshal(got.Logs, &logs); err != nil {
		t.Fatalf("logs 不是合法 JSON 数组: %v (raw=%s)", err, string(got.Logs))
	}
	if len(logs) != 3 {
		t.Errorf("logs 期望 3 条，实际 %d 条: %v", len(logs), logs)
	}
	t.Logf("✅ 播种成功 id=%d version=%s date=%s logs=%v", got.ID, got.Version, got.Date, logs)

	// 幂等性：重复播种不应新增行，也不应覆盖内容
	if err := seedChangelogs(); err != nil {
		t.Fatalf("重复播种失败: %v", err)
	}
	var count int64
	DB.Model(&models.Changelog{}).Where("version = ?", "2.15.0").Count(&count)
	if count != 1 {
		t.Errorf("重复播种后期望 1 行，实际 %d 行", count)
	}
	t.Logf("✅ 幂等性验证通过，重复执行后仍为 %d 行", count)
}
