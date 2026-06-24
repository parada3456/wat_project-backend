package storage

import (
	"context"
	"io"
)

type StoragePort interface {
	UploadFile(ctx context.Context, bucket, key string, data io.Reader, contentType string) (url string, err error)
	DeleteFile(ctx context.Context, bucket, key string) error
}
