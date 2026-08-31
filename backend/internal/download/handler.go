package download

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fatWill/agent-team-blog/backend/config"
	"github.com/gin-gonic/gin"
)

// allowedHosts 允许代理下载的域名白名单，由 main 在启动时通过 SetDownloadConfig 注入
// 默认值保持线上现状，避免未配置环境变量时行为变化
var allowedHosts = map[string]bool{
	"assets.fatwill.cloud": true,
	"pic.fatwill.cloud":    true,
	"fatwill-cloud-1253664788.cos.ap-guangzhou.myqcloud.com": true,
}

// SetDownloadConfig 设置下载代理配置（域名白名单），空白名单则保留默认值
func SetDownloadConfig(cfg *config.DownloadConfig) {
	if cfg == nil || len(cfg.AllowedHosts) == 0 {
		return
	}
	hosts := make(map[string]bool, len(cfg.AllowedHosts))
	for _, h := range cfg.AllowedHosts {
		if h = strings.ToLower(strings.TrimSpace(h)); h != "" {
			hosts[h] = true
		}
	}
	if len(hosts) > 0 {
		allowedHosts = hosts
	}
}

// 代理请求超时时间
var proxyClient = &http.Client{
	Timeout: 120 * time.Second,
}

// Proxy GET /api/download — 代理下载远程文件，解决前端 CORS 问题
func Proxy(c *gin.Context) {
	rawURL := c.Query("url")
	filename := c.Query("name")

	if rawURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "缺少 url 参数"})
		return
	}

	// 解析并校验 URL
	parsed, err := url.Parse(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "statusCode": 400, "statusMessage": "无效的 URL，仅支持 http/https 协议"})
		return
	}

	// 域名白名单校验
	host := strings.ToLower(parsed.Hostname())
	if !allowedHosts[host] {
		c.JSON(http.StatusForbidden, gin.H{"error": true, "statusCode": 403, "statusMessage": "该域名不在允许下载的白名单中"})
		return
	}

	// 发起远程请求
	resp, err := proxyClient.Get(rawURL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": true, "statusCode": 502, "statusMessage": "请求远程资源失败"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": true, "statusCode": resp.StatusCode, "statusMessage": "远程服务器返回错误"})
		return
	}

	// 设置下载文件名（默认从 URL 路径中提取）
	if filename == "" {
		parts := strings.Split(parsed.Path, "/")
		if len(parts) > 0 {
			filename = parts[len(parts)-1]
		}
		if filename == "" {
			filename = "download"
		}
	}

	// 设置响应头
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)

	// 透传 Content-Type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		c.Header("Content-Type", contentType)
	} else {
		c.Header("Content-Type", "application/octet-stream")
	}

	// 透传 Content-Length
	if cl := resp.Header.Get("Content-Length"); cl != "" {
		c.Header("Content-Length", cl)
	}

	// 流式透传响应体
	c.Status(http.StatusOK)
	io.Copy(c.Writer, resp.Body)
}
