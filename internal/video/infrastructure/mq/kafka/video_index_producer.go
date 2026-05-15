package kafka

import (
	"context"
	"encoding/json"

	"vicomova/pkg/infrastructure/kafka"
	"vicomova/pkg/log"
)

type VideoIndexEvent struct {
	VideoID     int64  `json:"video_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type VideoIndexProducer struct {
	sender *kafka.Sender
	topic  string
}

func NewVideoIndexProducer(sender *kafka.Sender, topic string) *VideoIndexProducer {
	return &VideoIndexProducer{
		sender: sender,
		topic:  topic,
	}
}

func (p *VideoIndexProducer) SendIndexEvent(ctx context.Context, event *VideoIndexEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &kafka.Message{
		Key:   nil,
		Value: data,
	}
	if err := p.sender.Send(ctx, kafka.GetBrokers(), p.topic, msg); err != nil {
		log.Error.Printf("VideoIndexProducer: send failed, videoID=%d: %v", event.VideoID, err)
		return err
	}
	log.Debug.Printf("VideoIndexProducer: sent index event, videoID=%d, topic=%s", event.VideoID, p.topic)
	return nil
}
