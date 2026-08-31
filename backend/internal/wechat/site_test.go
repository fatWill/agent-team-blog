package wechat

import (
	"strings"
	"testing"
)

func TestArticleURLUsesConfiguredSiteURL(t *testing.T) {
	original := SiteURL()
	defer SetSiteURL(original)

	tests := []struct {
		name    string
		siteURL string
		want    string
	}{
		{"默认域名", "https://fatwill.cloud", "https://fatwill.cloud/articles/abc123"},
		{"自定义域名", "https://test.example.com", "https://test.example.com/articles/abc123"},
		{"结尾斜杠自动去除", "https://test.example.com/", "https://test.example.com/articles/abc123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetSiteURL(tt.siteURL)
			if got := articleURL("abc123"); got != tt.want {
				t.Fatalf("articleURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConverterEmbedsConfiguredSiteURL(t *testing.T) {
	original := SiteURL()
	defer SetSiteURL(original)

	SetSiteURL("https://test.example.com")

	content := []byte(`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"hello"}]}]}`)
	result, err := TiptapToWechatHTML(content, "abc123")
	if err != nil {
		t.Fatalf("TiptapToWechatHTML() error = %v", err)
	}
	if !strings.Contains(result.HTML, `href="https://test.example.com/articles/abc123"`) {
		t.Fatalf("HTML 未使用配置的 SiteURL:\n%s", result.HTML)
	}
	if strings.Contains(result.HTML, "fatwill.cloud") {
		t.Fatalf("HTML 仍包含硬编码域名:\n%s", result.HTML)
	}
}
