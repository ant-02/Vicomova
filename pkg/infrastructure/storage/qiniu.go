package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Domain    string `mapstructure:"domain"`
}

type QiniuStorage struct {
	mac    *qbox.Mac
	bucket string
	domain string
}

func NewQiniuStorage(cfg *QiniuConfig) *QiniuStorage {
	mac := qbox.NewMac(cfg.AccessKey, cfg.SecretKey)
	return &QiniuStorage{
		mac:    mac,
		bucket: cfg.Bucket,
		domain: cfg.Domain,
	}
}

func (s *QiniuStorage) Upload(ctx context.Context, key string, reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	putPolicy := storage.PutPolicy{
		Scope: s.bucket,
	}
	upToken := putPolicy.UploadToken(s.mac)

	cfg := storage.Config{}
	uploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}
	err = uploader.Put(ctx, &ret, upToken, key, data, nil)
	if err != nil {
		return "", err
	}

	return s.GetURL(ctx, key)
}

func (s *QiniuStorage) Delete(ctx context.Context, key string) error {
	bucketManager := storage.NewBucketManager(s.mac, &storage.Config{})
	return bucketManager.Delete(s.bucket, key)
}

func (s *QiniuStorage) GetURL(ctx context.Context, key string) (string, error) {
	return fmt.Sprintf("http://%s/%s", s.domain, key), nil
}
