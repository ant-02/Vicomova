package storage

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type LocalDiskConfig struct {
	BasePath string `mapstructure:"base_path"`
	BaseURL  string `mapstructure:"base_url"`
}

type LocalDiskStorage struct {
	BasePath string
	BaseURL  string
}

func NewLocalDiskStorage(cfg *LocalDiskConfig) *LocalDiskStorage {
	return &LocalDiskStorage{
		BasePath: cfg.BasePath,
		BaseURL:  cfg.BaseURL,
	}
}

func (s *LocalDiskStorage) Upload(ctx context.Context, key string, reader io.Reader) (string, error) {
	fullPath := filepath.Join(s.BasePath, key)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}

	return s.GetURL(ctx, key)
}

func (s *LocalDiskStorage) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(s.BasePath, key)
	return os.Remove(fullPath)
}

func (s *LocalDiskStorage) GetURL(ctx context.Context, key string) (string, error) {
	base, err := url.Parse(s.BaseURL)
	if err != nil {
		return "", err
	}
	return base.JoinPath(key).String(), nil
}
