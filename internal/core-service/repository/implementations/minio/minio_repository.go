package minio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/core-service/repository/interfaces"
	bdIface "github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/minio/minio-go/v7"
)

var _ interfaces.FileStorage = (*StorageRepository)(nil)

type StorageRepository struct {
	client     bdIface.MinioIface
	bucketName string
	address    string
	useSSL     bool
}

func NewStorageRepository(cl bdIface.MinioIface, cfg *config.MinIOConfig) *StorageRepository {
	addr := cfg.Address
	if os.Getenv(cfg.Address) != "" {
		addr = os.Getenv(cfg.Address)
	}
	return &StorageRepository{
		client:     cl,
		bucketName: cfg.BucketName,
		address:    addr,
		useSSL:     cfg.UseSSL,
	}
}

func (m *StorageRepository) Upload(ctx context.Context, filename string, data []byte, contentType string) (string, error) {
	// Загружаем файл в MinIO
	_, err := m.client.PutObject(ctx, m.bucketName, filename, strings.NewReader(string(data)), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	// Генерируем URL для доступа к файлу
	return m.GetURL(ctx, filename), nil
}

func (m *StorageRepository) Delete(ctx context.Context, filename string) error {
	err := m.client.RemoveObject(ctx, m.bucketName, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	return nil
}

func (m *StorageRepository) DeleteByURL(ctx context.Context, url string) error {
	// Извлекаем имя файла из URL
	filename := filepath.Base(url)
	return m.Delete(ctx, filename)
}

func (m *StorageRepository) GetURL(ctx context.Context, filename string) string {
	protocol := "http"
	if m.useSSL {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s/%s/%s", protocol, m.address, m.bucketName, filename)
}

// PresignedURL генерирует URL с временным доступом (опционально)
func (m *StorageRepository) PresignedURL(ctx context.Context, filename string, expiry time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, m.bucketName, filename, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}
