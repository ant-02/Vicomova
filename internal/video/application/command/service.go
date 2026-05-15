package command

import (
	"context"

	"vicomova/internal/video/domain/repository"
	kafkapkg "vicomova/internal/video/infrastructure/mq/kafka"
	"vicomova/pkg/constants"
	"vicomova/pkg/infrastructure/oss"
	"vicomova/pkg/log"
)

type VideoCommandService struct {
	repo     repository.VideoRepository
	cache    repository.VideoCache
	oss      oss.OSS
	producer *kafkapkg.VideoIndexProducer
}

func NewVideoCommandService(repo repository.VideoRepository, cache repository.VideoCache, ossClient oss.OSS) *VideoCommandService {
	return &VideoCommandService{
		repo:  repo,
		cache: cache,
		oss:   ossClient,
	}
}

func (s *VideoCommandService) SetIndexProducer(producer *kafkapkg.VideoIndexProducer) {
	s.producer = producer
}

// sendIndexEvent 发送视频索引事件到 Kafka
func (s *VideoCommandService) sendIndexEvent(ctx context.Context, videoID int64, title, description string) {
	if s.producer == nil {
		return
	}
	event := &kafkapkg.VideoIndexEvent{
		VideoID:     videoID,
		Title:       title,
		Description: description,
	}
	if err := s.producer.SendIndexEvent(ctx, event); err != nil {
		log.Error.Printf("VideoCommandService: send index event failed: %v", err)
	} else {
		log.Debug.Printf("VideoCommandService: sent index event, videoID=%d", videoID)
	}
}

func NewVideoCommandServiceWithKafka(repo repository.VideoRepository, cache repository.VideoCache, ossClient oss.OSS, producer *kafkapkg.VideoIndexProducer) *VideoCommandService {
	svc := &VideoCommandService{
		repo:     repo,
		cache:    cache,
		oss:      ossClient,
		producer: producer,
	}
	return svc
}

// sendIndexEventToKafka 异步发送索引事件（不阻塞主流程）
func (s *VideoCommandService) sendIndexEventAsync(ctx context.Context, videoID int64, title, description string) {
	go func() {
		_, cancel := context.WithTimeout(context.Background(), constants.KafkaSendTimeout)
		defer cancel()
		s.sendIndexEvent(ctx, videoID, title, description)
	}()
}
