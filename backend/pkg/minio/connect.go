package minio

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// TODO: пределать под енв

func NewMinioPool(cfg *config.MinIOConfig) (*minio.Client, error) {
	// Инициализация MinIO клиента
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
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

	return client, nil
}
