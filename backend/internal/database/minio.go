package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/flowershop/backend/internal/config"
	"github.com/flowershop/backend/internal/constants"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ConnectMinIO 建立 MinIO 连接并确保 bucket 存在。
func ConnectMinIO(cfg *config.Config) (*minio.Client, error) {
	client, err := minio.New(cfg.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure: cfg.MinIOUseSSL,
	})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	exists, err := client.BucketExists(ctx, cfg.MinIOBucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.MinIOBucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}
	// 设置公开读策略，便于前端通过 nginx /minio 代理直接展示图片
	policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, cfg.MinIOBucket)
	_ = client.SetBucketPolicy(ctx, cfg.MinIOBucket, policy)
	return client, nil
}

// LogMinIO 输出 MinIO 连接日志。
func LogMinIO(logger *slog.Logger, cfg *config.Config) {
	logger.Info(constants.LogMinIOConnected, "endpoint", cfg.MinIOEndpoint, "bucket", cfg.MinIOBucket)
}
