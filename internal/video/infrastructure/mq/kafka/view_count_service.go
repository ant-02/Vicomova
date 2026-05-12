package kafka

import (
	"context"

	"vicomova/internal/video/domain/service"
	"vicomova/pkg/log"
)

type viewCountServiceImpl struct {
	producer *ViewCountProducer
}

func NewViewCountService(producer *ViewCountProducer) service.ViewCountService {
	return &viewCountServiceImpl{producer: producer}
}

func (s *viewCountServiceImpl) Record(ctx context.Context, videoID int64) error {
	if s.producer == nil {
		return nil
	}
	if err := s.producer.SendIncrement(ctx, videoID); err != nil {
		log.Error.Printf("ViewCountService: send failed, videoID=%d: %v", videoID, err)
		return err
	}
	log.Debug.Printf("ViewCountService: recorded videoID=%d", videoID)
	return nil
}