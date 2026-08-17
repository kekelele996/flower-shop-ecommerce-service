package util

import (
	"log/slog"
	"os"
)

// Logger 全局日志实例，全栈（handler/service/middleware）统一引用。
var Logger *slog.Logger

// InitLogger 初始化结构化日志，输出 JSON 便于采集。
func InitLogger(env string) *slog.Logger {
	level := slog.LevelInfo
	if env == "development" {
		level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	Logger = slog.New(handler)
	slog.SetDefault(Logger)
	return Logger
}

// GetLogger 返回全局 logger。
func GetLogger() *slog.Logger {
	if Logger == nil {
		return InitLogger("development")
	}
	return Logger
}
