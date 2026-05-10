package storage

import "context"

type ImageStorage interface {
	Upload(ctx context.Context, key string, reader io.Reader) (string, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) (string, error)
}
