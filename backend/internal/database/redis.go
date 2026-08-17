package database

import (
	"context"
	"log/slog"
	"time"

	"github.com/flowershop/backend/internal/config"
	"github.com/flowershop/backend/internal/constants"
	"github.com/redis/go-redis/v9"
)

// ConnectRedis 建立 Redis 连接。
func ConnectRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}

// LogRedis 输出 Redis 连接日志。
func LogRedis(logger *slog.Logger, cfg *config.Config) {
	logger.Info(constants.LogRedisConnected, "addr", cfg.RedisAddr)
}
