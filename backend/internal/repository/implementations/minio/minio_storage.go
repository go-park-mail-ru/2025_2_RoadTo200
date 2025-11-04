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
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioStorage struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	useSSL     bool
}

func NewMinIOStorage(cfg *config.MinIOConfig) (interfaces.FileStorage, error) {
	// Инициализация MinIO клиента
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	storage := &minioStorage{
		client:     client,
		bucketName: cfg.BucketName,
		endpoint:   cfg.Endpoint,
		useSSL:     cfg.UseSSL,
	}

	// Проверяем существование бакета
	exists, err := client.BucketExists(context.Background(), cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		// Создаем бакет если не существует
		err = client.MakeBucket(context.Background(), cfg.BucketName, minio.MakeBucketOptions{
			Region: cfg.Region,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}

		// Настраиваем политику доступа (публичный доступ для чтения)
		policy := `{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": "*",
                    "Action": ["s3:GetObject"],
                    "Resource": ["arn:aws:s3:::` + cfg.BucketName + `/*"]
                }
            ]
        }`
		err = client.SetBucketPolicy(context.Background(), cfg.BucketName, policy)
		if err != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %w", err)
		}
	}

	return storage, nil
}

func (m *minioStorage) Upload(filename string, data []byte, contentType string) (string, error) {
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

func (m *minioStorage) Delete(filename string) error {
	ctx := context.Background()

	err := m.client.RemoveObject(ctx, m.bucketName, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from MinIO: %w", err)
	}

	return nil
}

func (m *minioStorage) DeleteByURL(url string) error {
	// Извлекаем имя файла из URL
	filename := filepath.Base(url)
	return m.Delete(filename)
}

func (m *minioStorage) GetURL(filename string) string {
	protocol := "http"
	if m.useSSL {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s/%s/%s", protocol, m.endpoint, m.bucketName, filename)
}

// PresignedURL генерирует URL с временным доступом (опционально)
func (m *minioStorage) PresignedURL(filename string, expiry time.Duration) (string, error) {
	ctx := context.Background()

	url, err := m.client.PresignedGetObject(ctx, m.bucketName, filename, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}
