package kafka

import (
	"context"
	"encoding/binary"

	pkgKafka "vicomova/pkg/infrastructure/kafka"
	"vicomova/pkg/log"
)

type ViewCountProducer struct {
	sender *pkgKafka.Sender
	topic  string
}

func NewViewCountProducer(sender *pkgKafka.Sender, topic string) *ViewCountProducer {
	return &ViewCountProducer{
		sender: sender,
		topic:  topic,
	}
}

func (p *ViewCountProducer) Record(ctx context.Context, videoID int64) error {
	data := make([]byte, 8)
	n := binary.PutVarint(data, videoID)
	data = data[:n]

	msg := &pkgKafka.Message{
		Key:   nil,
		Value: data,
	}
	if err := p.sender.Send(ctx, pkgKafka.GetBrokers(), p.topic, msg); err != nil {
		log.Error.Printf("ViewCountProducer: send failed, videoID=%d: %v", videoID, err)
		return err
	}
	log.Debug.Printf("ViewCountProducer: sent message, videoID=%d, topic=%s", videoID, p.topic)
	return nil
}

func (p *ViewCountProducer) RecordBatch(ctx context.Context, videoIDs []int64) error {
	if len(videoIDs) == 0 {
		return nil
	}

	messages := make([]*pkgKafka.Message, len(videoIDs))
	for i, id := range videoIDs {
		data := make([]byte, 8)
		n := binary.PutVarint(data, id)
		messages[i] = &pkgKafka.Message{Key: nil, Value: data[:n]}
	}

	return p.sender.SendBatch(ctx, pkgKafka.GetBrokers(), p.topic, messages)
}
