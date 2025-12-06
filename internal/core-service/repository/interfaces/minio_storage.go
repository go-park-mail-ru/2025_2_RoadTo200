package interfaces

import "context"

type StorageRepository interface {
	Upload(ctx context.Context, filename string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, filename string) error
	DeleteByURL(ctx context.Context, url string) error
	GetURL(ctx context.Context, filename string) string
}
