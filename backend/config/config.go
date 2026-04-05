package config

import "os"

// Config 应用配置
type Config struct {
	Server ServerConfig
	DB     DBConfig
	Redis  RedisConfig
	Upload UploadConfig
	COS    COSConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port       string
	CORSOrigin string
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
	SecretID  string
	SecretKey string
	Bucket    string
	Region    string
	BaseURL   string // CDN 访问域名，如 https://assets.fatwill.cloud
}

// UploadConfig 上传配置
type UploadConfig struct {
	Dir    string // 上传文件存储目录
	TmpDir string // 分片临时目录
}

// Load 从环境变量加载配置，提供默认值
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:       getEnv("SERVER_PORT", "8080"),
			CORSOrigin: getEnv("CORS_ORIGIN", "https://fatwill.cloud"),
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
			SecretID:  getEnv("COS_SECRET_ID", ""),
			SecretKey: getEnv("COS_SECRET_KEY", ""),
			Bucket:    getEnv("COS_BUCKET", "fatwill-cloud-1253664788"),
			Region:    getEnv("COS_REGION", "ap-guangzhou"),
			BaseURL:   getEnv("COS_BASE_URL", "https://fatwill-cloud-1253664788.cos.ap-guangzhou.myqcloud.com"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}