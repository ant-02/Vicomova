package oss

import (
	"context"
	"fmt"
	"time"

	"vicomova/pkg/config"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
)

type OSSConfig = config.QiniuOSS

type OSSQiniu struct {
	credentials *credentials.Credentials
	bucket      string
	domain      string
}

func NewOSSQiniu(cfg *OSSConfig) *OSSQiniu {
	creds := credentials.NewCredentials(cfg.AccessKey, cfg.SecretKey)
	return &OSSQiniu{
		credentials: creds,
		bucket:      cfg.Bucket,
		domain:      cfg.Domain,
	}
}

// GetUploadToken 获取上传凭证，前端使用此 token 直传到七牛云
// 同时返回 domain，供前端拼接完整 URL
func (s *OSSQiniu) GetUploadToken(ctx context.Context, key string, expire time.Duration) (string, string, error) {
	putPolicy, err := uptoken.NewPutPolicyWithKey(s.bucket, key, time.Now().Add(expire))
	if err != nil {
		return "", "", fmt.Errorf("failed to create put policy: %w", err)
	}
	token, err := uptoken.NewSigner(putPolicy, s.credentials).GetUpToken(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to get upload token: %w", err)
	}
	return token, s.domain, nil
}
