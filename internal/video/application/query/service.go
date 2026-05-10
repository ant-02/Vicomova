package query

import (
	"context"
	"fmt"
	"time"

	"vicomova/internal/video/domain/repository"
	"vicomova/internal/video/domain/service"
	"vicomova/pkg/infrastructure/oss"
)

type VideoQueryService struct {
	repo    repository.VideoRepository
	hotAlgo service.HotAlgorithm
	oss     oss.OSS
}

func NewVideoQueryService(repo repository.VideoRepository, hotAlgo service.HotAlgorithm, oss oss.OSS) *VideoQueryService {
	return &VideoQueryService{
		repo:    repo,
		hotAlgo: hotAlgo,
		oss:     oss,
	}
}

// GetUploadToken 获取上传凭证，供前端直传到 OSS
func (s *VideoQueryService) GetUploadToken(ctx context.Context, key string, expireSeconds int64) (*GetUploadTokenResult, error) {
	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}

	expire := time.Duration(expireSeconds) * time.Second
	if expire == 0 {
		expire = 3600 * time.Second // 默认 1 小时
	}

	token, err := s.oss.GetUploadToken(ctx, key, expire)
	if err != nil {
		return nil, fmt.Errorf("failed to get upload token: %w", err)
	}

	// 获取 domain（从 oss 接口获取 URL 然后提取 domain）
	url, _ := s.oss.GetURL(ctx, key)
	domain := url[:len(url)-len(key)-1] // 简单提取 domain，实际应该从配置获取

	return &GetUploadTokenResult{
		Token: token,
		Key:   key,
		Domain: domain,
	}, nil
}

// GetUploadTokenResult 上传凭证结果
type GetUploadTokenResult struct {
	Token  string
	Key    string
	Domain string
}
