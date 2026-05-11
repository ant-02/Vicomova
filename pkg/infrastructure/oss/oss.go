package oss

import (
	"context"
	"time"
)

type UploadResult struct {
	URL string
}

// OSS 对象存储接口，用于前端直传场景
// 后端只负责生成上传凭证，不再负责实际上传
type OSS interface {
	// GetUploadToken 获取上传凭证和上传地址，前端用此 token 直传到七牛云
	GetUploadToken(ctx context.Context, key string, expire time.Duration) (token string, host string, domain string, err error)
}
