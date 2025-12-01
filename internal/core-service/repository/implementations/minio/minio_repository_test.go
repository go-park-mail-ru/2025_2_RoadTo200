package minio

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/url"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
)

// MockMinioClient реализует мок для minio.Client
type MockMinioClient struct {
	PutObjectFunc          func(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error)
	RemoveObjectFunc       func(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error
	PresignedGetObjectFunc func(ctx context.Context, bucketName, objectName string, expiry time.Duration, reqParams url.Values) (*url.URL, error)
}

func (m *MockMinioClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	if m.PutObjectFunc != nil {
		return m.PutObjectFunc(ctx, bucketName, objectName, reader, objectSize, opts)
	}
	return minio.UploadInfo{}, nil
}

func (m *MockMinioClient) RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	if m.RemoveObjectFunc != nil {
		return m.RemoveObjectFunc(ctx, bucketName, objectName, opts)
	}
	return nil
}

func (m *MockMinioClient) PresignedGetObject(ctx context.Context, bucketName, objectName string, expiry time.Duration, reqParams url.Values) (*url.URL, error) {
	if m.PresignedGetObjectFunc != nil {
		return m.PresignedGetObjectFunc(ctx, bucketName, objectName, expiry, reqParams)
	}
	return &url.URL{}, nil
}

func TestStorageRepository_Upload(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		data        []byte
		contentType string
		setupMock   func(*MockMinioClient)
		wantURL     string
		wantErr     bool
	}{
		{
			name:        "successful upload",
			filename:    "test.jpg",
			data:        []byte("test data"),
			contentType: "image/jpeg",
			setupMock: func(m *MockMinioClient) {
				m.PutObjectFunc = func(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
					assert.Equal(t, "test-bucket", bucketName)
					assert.Equal(t, "test.jpg", objectName)
					assert.Equal(t, int64(9), objectSize)
					assert.Equal(t, "image/jpeg", opts.ContentType)
					return minio.UploadInfo{}, nil
				}
			},
			wantURL: "http://localhost:9000/test-bucket/test.jpg",
			wantErr: false,
		},
		{
			name:        "upload error",
			filename:    "test.jpg",
			data:        []byte("test data"),
			contentType: "image/jpeg",
			setupMock: func(m *MockMinioClient) {
				m.PutObjectFunc = func(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
					return minio.UploadInfo{}, errors.New("upload failed")
				}
			},
			wantURL: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockMinioClient{}
			tt.setupMock(mockClient)

			cfg := &config.MinIOConfig{
				BucketName: "test-bucket",
				Address:    "localhost:9000",
				UseSSL:     false,
			}

			repo := &StorageRepository{
				client:     mockClient,
				bucketName: cfg.BucketName,
				address:    cfg.Address,
				useSSL:     cfg.UseSSL,
			}

			ctx := context.Background()
			gotURL, err := repo.Upload(ctx, tt.filename, tt.data, tt.contentType)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to upload file to MinIO")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantURL, gotURL)
			}
		})
	}
}

func TestStorageRepository_Delete(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		setupMock func(*MockMinioClient)
		wantErr   bool
	}{
		{
			name:     "successful delete",
			filename: "test.jpg",
			setupMock: func(m *MockMinioClient) {
				m.RemoveObjectFunc = func(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
					assert.Equal(t, "test-bucket", bucketName)
					assert.Equal(t, "test.jpg", objectName)
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:     "delete error",
			filename: "test.jpg",
			setupMock: func(m *MockMinioClient) {
				m.RemoveObjectFunc = func(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
					return errors.New("delete failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockMinioClient{}
			tt.setupMock(mockClient)

			cfg := &config.MinIOConfig{
				BucketName: "test-bucket",
				Address:    "localhost:9000",
				UseSSL:     false,
			}

			repo := &StorageRepository{
				client:     mockClient,
				bucketName: cfg.BucketName,
				address:    cfg.Address,
				useSSL:     cfg.UseSSL,
			}

			ctx := context.Background()
			err := repo.Delete(ctx, tt.filename)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete file from MinIO")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStorageRepository_DeleteByURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		setupMock func(*MockMinioClient)
		wantErr   bool
	}{
		{
			name: "successful delete by URL",
			url:  "http://localhost:9000/test-bucket/test.jpg",
			setupMock: func(m *MockMinioClient) {
				m.RemoveObjectFunc = func(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
					assert.Equal(t, "test-bucket", bucketName)
					assert.Equal(t, "test.jpg", objectName)
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "delete by URL with path",
			url:  "http://localhost:9000/test-bucket/path/to/test.jpg",
			setupMock: func(m *MockMinioClient) {
				m.RemoveObjectFunc = func(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
					assert.Equal(t, "test-bucket", bucketName)
					assert.Equal(t, "test.jpg", objectName)
					return nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockMinioClient{}
			tt.setupMock(mockClient)

			cfg := &config.MinIOConfig{
				BucketName: "test-bucket",
				Address:    "localhost:9000",
				UseSSL:     false,
			}

			repo := &StorageRepository{
				client:     mockClient,
				bucketName: cfg.BucketName,
				address:    cfg.Address,
				useSSL:     cfg.UseSSL,
			}

			ctx := context.Background()
			err := repo.DeleteByURL(ctx, tt.url)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStorageRepository_GetURL(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		useSSL   bool
		wantURL  string
	}{
		{
			name:     "HTTP URL",
			filename: "test.jpg",
			useSSL:   false,
			wantURL:  "http://localhost:9000/test-bucket/test.jpg",
		},
		{
			name:     "HTTPS URL",
			filename: "test.jpg",
			useSSL:   true,
			wantURL:  "https://localhost:9000/test-bucket/test.jpg",
		},
		{
			name:     "with path",
			filename: "path/to/test.jpg",
			useSSL:   false,
			wantURL:  "http://localhost:9000/test-bucket/path/to/test.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.MinIOConfig{
				BucketName: "test-bucket",
				Address:    "localhost:9000",
				UseSSL:     tt.useSSL,
			}

			repo := &StorageRepository{
				client:     &MockMinioClient{},
				bucketName: cfg.BucketName,
				address:    cfg.Address,
				useSSL:     cfg.UseSSL,
			}

			ctx := context.Background()
			gotURL := repo.GetURL(ctx, tt.filename)

			assert.Equal(t, tt.wantURL, gotURL)
		})
	}
}

func TestStorageRepository_PresignedURL(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		expiry    time.Duration
		setupMock func(*MockMinioClient)
		wantErr   bool
	}{
		{
			name:     "successful presigned URL",
			filename: "test.jpg",
			expiry:   time.Hour,
			setupMock: func(m *MockMinioClient) {
				m.PresignedGetObjectFunc = func(ctx context.Context, bucketName, objectName string, expiry time.Duration, reqParams url.Values) (*url.URL, error) {
					assert.Equal(t, "test-bucket", bucketName)
					assert.Equal(t, "test.jpg", objectName)
					assert.Equal(t, time.Hour, expiry)
					return url.Parse("http://localhost:9000/test-bucket/test.jpg?X-Amz-Expires=3600")
				}
			},
			wantErr: false,
		},
		{
			name:     "presigned URL error",
			filename: "test.jpg",
			expiry:   time.Hour,
			setupMock: func(m *MockMinioClient) {
				m.PresignedGetObjectFunc = func(ctx context.Context, bucketName, objectName string, expiry time.Duration, reqParams url.Values) (*url.URL, error) {
					return nil, errors.New("presigned URL failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockMinioClient{}
			tt.setupMock(mockClient)

			cfg := &config.MinIOConfig{
				BucketName: "test-bucket",
				Address:    "localhost:9000",
				UseSSL:     false,
			}

			repo := &StorageRepository{
				client:     mockClient,
				bucketName: cfg.BucketName,
				address:    cfg.Address,
				useSSL:     cfg.UseSSL,
			}

			ctx := context.Background()
			gotURL, err := repo.PresignedURL(ctx, tt.filename, tt.expiry)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to generate presigned URL")
			} else {
				assert.NoError(t, err)
				assert.Contains(t, gotURL, "http://localhost:9000/test-bucket/test.jpg")
				assert.Contains(t, gotURL, "X-Amz-Expires=3600")
			}
		})
	}
}

func TestNewStorageRepository(t *testing.T) {
	mockClient := &MockMinioClient{}
	cfg := &config.MinIOConfig{
		BucketName: "test-bucket",
		Address:    "localhost:9000",
		UseSSL:     true,
	}

	repo := NewStorageRepository(mockClient, cfg)

	assert.NotNil(t, repo)
	assert.Equal(t, mockClient, repo.client)
	assert.Equal(t, "test-bucket", repo.bucketName)
	assert.Equal(t, "localhost:9000", repo.address)
	assert.True(t, repo.useSSL)
}

// Вспомогательная функция для создания reader из байтов
type bytesReader struct {
	*bytes.Reader
}

func (b bytesReader) Read(p []byte) (n int, err error) {
	return b.Reader.Read(p)
}

var _ io.Reader = (*bytesReader)(nil)
