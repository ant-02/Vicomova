package storage

import (
    "context"
    "io"
)

type VideoStorage interface {
    Upload(ctx context.Context, key string, reader io.Reader) (string, error)
    Delete(ctx context.Context, key string) error
    GetURL(ctx context.Context, key string) (string, error)
}
