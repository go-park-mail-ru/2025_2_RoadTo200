package interfaces

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
)

type PgxIface interface {
	pgx.Tx
	Begin(context.Context) (pgx.Tx, error)
	Close()
}

type MinioIface interface {
	PutObject(context.Context, string, string, io.Reader, int64, minio.PutObjectOptions) (minio.UploadInfo, error)
	RemoveObject(context.Context, string, string, minio.RemoveObjectOptions) error
	PresignedGetObject(context.Context, string, string, time.Duration, url.Values) (*url.URL, error)
}
