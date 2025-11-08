package minio

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/repository/interfaces"
	"github.com/minio/minio-go/v7"
)

type storageRepository struct {
	client     *minio.Client
	bucketName string
	address    string
	useSSL     bool
}

func NewStorageRepository(cl *minio.Client, cfg *config.MinIOConfig) interfaces.FileStorage {
	return &storageRepository{
		client:     cl,
		bucketName: cfg.BucketName,
		address:    cfg.Address,
		useSSL:     cfg.UseSSL,
	}
}

func (m *storageRepository) Upload(filename string, data []byte, contentType string) (string, error) {
	ctx := context.Background()

	// Загружаем файл в MinIO
	_, err := m.client.PutObject(ctx, m.bucketName, filename, strings.NewReader(string(data)), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	// Генерируем URL для доступа к файлу
	return m.GetURL(filename), nil
}

func (m *storageRepository) Delete(filename string) error {
	ctx := context.Background()

	err := m.client.RemoveObject(ctx, m.bucketName, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	return nil
}

func (m *storageRepository) DeleteByURL(url string) error {
	// Извлекаем имя файла из URL
	filename := filepath.Base(url)
	return m.Delete(filename)
}

func (m *storageRepository) GetURL(filename string) string {
	protocol := "http"
	if m.useSSL {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s/%s/%s", protocol, m.address, m.bucketName, filename)
}

// PresignedURL генерирует URL с временным доступом (опционально)
func (m *storageRepository) PresignedURL(filename string, expiry time.Duration) (string, error) {
	ctx := context.Background()

	url, err := m.client.PresignedGetObject(ctx, m.bucketName, filename, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}
