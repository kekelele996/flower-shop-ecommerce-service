package database

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/flowershop/backend/internal/config"
	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect 建立 PostgreSQL 连接。
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}

// Migrate 自动迁移全部表。
func Migrate(db *gorm.DB) (int, error) {
	models := []interface{}{
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.CartItem{},
		&model.Order{},
		&model.OrderItem{},
		&model.Payment{},
		&model.Logistics{},
		&model.Review{},
		&model.Coupon{},
		&model.AuditLog{},
	}
	if err := db.AutoMigrate(models...); err != nil {
		return 0, err
	}
	return len(models), nil
}

// LogDBStats 输出数据库连接日志。
func LogDBStats(logger *slog.Logger, cfg *config.Config) {
	logger.Info(constants.LogDBConnected, "host", cfg.DBHost, "port", cfg.DBPort, "db", cfg.DBName)
}
