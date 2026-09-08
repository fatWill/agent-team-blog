package config

import (
	"os"
	"strings"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	DB       DBConfig
	Redis    RedisConfig
	Upload   UploadConfig
	COS      COSConfig
	Download DownloadConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port          string
	CORSOrigins   []string // 允许的跨域来源列表（逗号分隔注入，域名迁移期可新旧并存）
	SiteURL       string   // 站点根 URL（用于拼接对外可访问的绝对页面链接）
	IP2RegionPath string   // ip2region.xdb 数据文件路径
}

// DBConfig SQLite 配置
type DBConfig struct {
	Path string // SQLite 数据库文件路径
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

// COSConfig 腾讯云 COS 配置
type COSConfig struct {
	SecretID      string
	SecretKey     string
	Bucket        string
	Region        string
	BaseURL       string   // COS 原始域名（SDK 内部使用）
	CustomDomain  string   // 自定义域名（返回给前端的图片 URL）
	LegacyDomains []string // 历史自定义域名（仅用于解析存量 URL，不用于生成新 URL）
}

// UploadConfig 上传配置
type UploadConfig struct {
	Dir    string // 上传文件存储目录
	TmpDir string // 分片临时目录
}

// DownloadConfig 下载代理配置
type DownloadConfig struct {
	AllowedHosts []string // 允许代理下载的域名白名单
}

// Load 从环境变量加载配置，提供默认值
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			// 域名迁移期（fatwill.cloud → fatwill.cn）默认同时放行新旧域名，待旧域名下线后可精简
			CORSOrigins: getEnvList("CORS_ORIGIN",
				"https://fatwill.cn,https://www.fatwill.cn,https://fatwill.cloud,https://www.fatwill.cloud"),
			SiteURL:       getEnv("SITE_URL", "https://fatwill.cn"),
			IP2RegionPath: getEnv("IP2REGION_PATH", "data/ip2region.xdb"),
		},
		DB: DBConfig{
			Path: getEnv("DB_PATH", "/root/blog-data/blog.db"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "127.0.0.1"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Upload: UploadConfig{
			Dir:    getEnv("UPLOAD_DIR", "/root/blog-uploads"),
			TmpDir: getEnv("UPLOAD_TMP_DIR", "/root/blog-uploads/tmp"),
		},
		COS: COSConfig{
			SecretID:     getEnv("COS_SECRET_ID", ""),
			SecretKey:    getEnv("COS_SECRET_KEY", ""),
			Bucket:       getEnv("COS_BUCKET", "fatwill-cloud-1253664788"),
			Region:       getEnv("COS_REGION", "ap-guangzhou"),
			BaseURL:      getEnv("COS_BASE_URL", "https://fatwill-cloud-1253664788.cos.ap-guangzhou.myqcloud.com"),
			CustomDomain: getEnv("COS_CUSTOM_DOMAIN", "https://assets.fatwill.cn"),
			// 存量文章/相册中的图片 URL 仍是旧域名，删除时需能反解出 COS key
			LegacyDomains: getEnvList("COS_LEGACY_DOMAINS",
				"https://assets.fatwill.cloud,https://cdn.fatwill.cloud,https://cdn.fatwill.cn"),
		},
		Download: DownloadConfig{
			// 旧域名保留在白名单中，保证存量文章内的图片仍可通过 /api/download 代理下载
			AllowedHosts: getEnvList("DOWNLOAD_ALLOWED_HOSTS",
				"assets.fatwill.cn,cdn.fatwill.cn,pic.fatwill.cn,assets.fatwill.cloud,cdn.fatwill.cloud,pic.fatwill.cloud,fatwill-cloud-1253664788.cos.ap-guangzhou.myqcloud.com"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getEnvList 读取逗号分隔的环境变量并切分为字符串切片，自动去除空白项
func getEnvList(key, fallback string) []string {
	raw := getEnv(key, fallback)
	parts := strings.Split(raw, ",")
	list := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			list = append(list, v)
		}
	}
	return list
}
