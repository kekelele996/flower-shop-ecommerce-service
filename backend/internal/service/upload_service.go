package service

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/util"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// UploadService 对象存储服务：上传商品/评价/头像图片到 MinIO。
type UploadService struct {
	client *minio.Client
	bucket string
	logger *slog.Logger
}

func NewUploadService(client *minio.Client, bucket string, logger *slog.Logger) *UploadService {
	return &UploadService{client: client, bucket: bucket, logger: logger}
}

// Upload 上传文件并返回访问 URL。
func (s *UploadService) Upload(ctx context.Context, file *multipart.FileHeader, kind string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", util.WrapAppError(constants.CodeUploadFailed, "open upload file failed", err)
	}
	defer src.Close()

	ext := strings.ToLower(path.Ext(file.Filename))
	if ext == "" {
		ext = ".jpg"
	}
	objectName := fmt.Sprintf("%s/%s%s", kind, time.Now().Format("20060102")+"/"+uuid.NewString(), ext)

	info, err := s.client.PutObject(ctx, s.bucket, objectName, src, file.Size, minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})
	if err != nil {
		return "", util.WrapAppError(constants.CodeUploadFailed, "put object to minio failed", err)
	}
	s.logger.Info(constants.LogUploadSucceeded, "key", objectName, "size", info.Size)
	return fmt.Sprintf("/minio/%s/%s", s.bucket, objectName), nil
}
