package changelog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/fatWill/agent-team-blog/backend/config"
	"github.com/fatWill/agent-team-blog/backend/pkg/db"
	"github.com/gin-gonic/gin"
)

// TestGetChangelogReturnsSeeded215 端到端验证 GET /api/changelog 能返回被播种的 2.15.0 条目，
// 且 logs 以 JSON 数组（而非转义字符串）形式输出，保证前端渲染正常。
func TestGetChangelogReturnsSeeded215(t *testing.T) {
	if err := db.Init(&config.DBConfig{Path: filepath.Join(t.TempDir(), "test.db")}); err != nil {
		t.Fatalf("Init 失败: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/changelog", GetChangelog)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/changelog", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("状态码期望 200，实际 %d", w.Code)
	}

	var resp struct {
		Changelog []struct {
			Version string   `json:"version"`
			Date    string   `json:"date"`
			Logs    []string `json:"logs"`
		} `json:"changelog"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("响应体解析失败: %v\nbody=%s", err, w.Body.String())
	}

	var found bool
	for _, item := range resp.Changelog {
		if item.Version == "2.15.0" {
			found = true
			if item.Date != "2026-09-16" {
				t.Errorf("date 期望 2026-09-16，实际 %s", item.Date)
			}
			if len(item.Logs) != 3 {
				t.Errorf("logs 期望 3 条，实际 %d", len(item.Logs))
			}
			t.Logf("✅ 接口返回 2.15.0：date=%s logs=%v", item.Date, item.Logs)
		}
	}
	if !found {
		t.Fatalf("响应中未找到 2.15.0，body=%s", w.Body.String())
	}
}
