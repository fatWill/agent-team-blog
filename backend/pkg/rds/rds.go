package rds

import (
	"context"
	"fmt"

	"github.com/fatWill/agent-team-blog/backend/config"
	"github.com/redis/go-redis/v9"
)

// RDB 全局 Redis 客户端
var RDB *redis.Client

// Init 初始化 Redis 连接
func Init(cfg *config.RedisConfig) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	// 测试连接
	if err := RDB.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("连接 Redis 失败: %w", err)
	}

	return nil
}
