package config

import "os"

// Config 应用配置
type Config struct {
	Server ServerConfig
	DB     DBConfig
	Redis  RedisConfig
	Upload UploadConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port       string
	CORSOrigin string
}

// DBConfig MySQL 配置
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
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
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "blog"),
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
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
