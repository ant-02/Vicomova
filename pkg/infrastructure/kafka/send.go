package kafka

import (
	"context"
	"sync"

	"vicomova/pkg/log"

	kafka "github.com/segmentio/kafka-go"
)

// Sender 发送者封装
type Sender struct {
	writers map[string]*kafka.Writer
	mu      sync.RWMutex
}

// NewSender 创建发送者
func NewSender() *Sender {
	return &Sender{
		writers: make(map[string]*kafka.Writer),
	}
}

// GetWriter 获取指定 topic 的 writer（缓存复用）
func (s *Sender) GetWriter(brokers []string, topic string) *kafka.Writer {
	s.mu.RLock()
	w, ok := s.writers[topic]
	s.mu.RUnlock()

	if ok {
		return w
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// double-check
	if w, ok = s.writers[topic]; ok {
		return w
	}

	w = &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.RoundRobin{},
		MaxAttempts:            3,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: false,
		ErrorLogger:            log.Error,
		Transport:              getTransport(),
	}
	s.writers[topic] = w

	return w
}

// Send 发送单条消息
func (s *Sender) Send(ctx context.Context, brokers []string, topic string, msg *Message) error {
	w := s.GetWriter(brokers, topic)

	kafkaMsg := kafka.Message{
		Key:   msg.Key,
		Value: msg.Value,
	}

	return w.WriteMessages(ctx, kafkaMsg)
}

// SendBatch 发送批量消息
func (s *Sender) SendBatch(ctx context.Context, brokers []string, topic string, messages []*Message) error {
	w := s.GetWriter(brokers, topic)

	kafkaMsgs := make([]kafka.Message, len(messages))
	for i, m := range messages {
		kafkaMsgs[i] = kafka.Message{Key: m.Key, Value: m.Value}
	}

	return w.WriteMessages(ctx, kafkaMsgs...)
}

// CloseWriter 关闭指定 topic 的 writer
func (s *Sender) CloseWriter(topic string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if w, ok := s.writers[topic]; ok {
		delete(s.writers, topic)
		return w.Close()
	}

	return nil
}

var (
	senderOnce sync.Once
	sender     *Sender
)

func GetSender() *Sender {
	senderOnce.Do(func() {
		sender = NewSender()
	})
	return sender
}

// CloseAllWriter 关闭所有 writer
func (s *Sender) CloseAllWriter() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for topic, w := range s.writers {
		if err := w.Close(); err != nil {
			log.Error.Printf("close writer[%s] failed: %v", topic, err)
		}
	}

	s.writers = make(map[string]*kafka.Writer)

	return nil
}
