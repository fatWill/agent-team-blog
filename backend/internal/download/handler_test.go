package download

import (
	"testing"

	"github.com/fatWill/agent-team-blog/backend/config"
)

func TestSetDownloadConfig(t *testing.T) {
	original := allowedHosts
	defer func() { allowedHosts = original }()

	tests := []struct {
		name     string
		cfg      *config.DownloadConfig
		allow    []string
		disallow []string
	}{
		{
			name:     "nil 配置保留默认白名单",
			cfg:      nil,
			allow:    []string{"assets.fatwill.cloud", "pic.fatwill.cloud"},
			disallow: []string{"evil.example.com"},
		},
		{
			name:     "空白名单保留默认值",
			cfg:      &config.DownloadConfig{AllowedHosts: []string{}},
			allow:    []string{"assets.fatwill.cloud"},
			disallow: []string{"evil.example.com"},
		},
		{
			name:     "自定义白名单覆盖默认值",
			cfg:      &config.DownloadConfig{AllowedHosts: []string{" Assets.Test.Com ", "cdn.test.com"}},
			allow:    []string{"assets.test.com", "cdn.test.com"},
			disallow: []string{"assets.fatwill.cloud"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowedHosts = original
			SetDownloadConfig(tt.cfg)
			for _, h := range tt.allow {
				if !allowedHosts[h] {
					t.Errorf("期望 %q 在白名单中", h)
				}
			}
			for _, h := range tt.disallow {
				if allowedHosts[h] {
					t.Errorf("期望 %q 不在白名单中", h)
				}
			}
		})
	}
}

func TestDefaultAllowedHostsFromConfig(t *testing.T) {
	original := allowedHosts
	defer func() { allowedHosts = original }()

	cfg := config.Load()
	allowedHosts = original
	SetDownloadConfig(&cfg.Download)

	for _, h := range []string{
		"assets.fatwill.cloud",
		"pic.fatwill.cloud",
		"fatwill-cloud-1253664788.cos.ap-guangzhou.myqcloud.com",
	} {
		if !allowedHosts[h] {
			t.Errorf("默认配置应包含 %q", h)
		}
	}
}
