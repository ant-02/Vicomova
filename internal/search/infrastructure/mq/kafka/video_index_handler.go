package kafka

import (
	"context"
	"encoding/json"
	"log"

	"vicomova/internal/search/application/query"
	pkgKafka "vicomova/pkg/infrastructure/kafka"
)

type VideoIndexHandler struct {
	searchSvc *query.SearchService
	consumer  *pkgKafka.Consumer
}

func NewVideoIndexHandler(searchSvc *query.SearchService) *VideoIndexHandler {
	return &VideoIndexHandler{
		searchSvc: searchSvc,
		consumer:  pkgKafka.NewConsumer(),
	}
}

func (h *VideoIndexHandler) Start(ctx context.Context, brokers []string, topic, groupID string) error {
	ch, err := h.consumer.ConsumeTopic(ctx, brokers, topic, groupID, 2, 100, func(msg *pkgKafka.Message) error {
		var event VideoIndexEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("VideoIndexHandler: unmarshal failed: %v", err)
			return nil
		}

		q := &query.IndexVideoQuery{
			VideoID:     event.VideoID,
			Title:       event.Title,
			Description: event.Description,
		}

		if err := h.searchSvc.IndexVideo(ctx, q); err != nil {
			log.Printf("VideoIndexHandler: index video failed: %v", err)
			return err
		}

		log.Printf("VideoIndexHandler: indexed video, videoID=%d", event.VideoID)
		return nil
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = h.consumer.Close()
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				log.Printf("VideoIndexHandler: received message, key=%s", string(msg.Key))
			}
		}
	}()

	log.Printf("VideoIndexHandler: started, topic=%s, group=%s", topic, groupID)
	return nil
}

func (h *VideoIndexHandler) Stop() error {
	return h.consumer.Close()
}
