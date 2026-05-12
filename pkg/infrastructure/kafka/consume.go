package kafka

import (
	"context"
	"sync"
	"time"

	kafka "github.com/segmentio/kafka-go"

	"vicomova/pkg/log"
)

// ConsumeFunc 消息处理函数
type ConsumeFunc func(msg *Message) error

// Consumer 消费者封装
type Consumer struct {
	readers map[string][]*kafka.Reader
	chans   map[string]chan *Message
	mu      sync.RWMutex
}

// NewConsumer 创建一个消费者
func NewConsumer() *Consumer {
	return &Consumer{
		readers: make(map[string][]*kafka.Reader),
		chans:   make(map[string]chan *Message),
	}
}

// ConsumeTopic 消费指定 topic 的消息
func (c *Consumer) ConsumeTopic(ctx context.Context, brokers []string, topic, groupID string, consumerNum, chanCap int, handler ConsumeFunc) (<-chan *Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 幂等：同一个 topic 直接返回已有 channel
	if ch, ok := c.chans[topic]; ok {
		return ch, nil
	}

	cap := 100
	if chanCap > 0 {
		cap = chanCap
	}
	ch := make(chan *Message, cap)
	c.chans[topic] = ch

	// 创建 N 个 reader
	readers := make([]*kafka.Reader, 0, consumerNum)
	for i := 0; i < consumerNum; i++ {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     groupID,
			MaxBytes:    10 * 1024 * 1024, // 10MB
			MaxAttempts: 3,
			ErrorLogger: log.Error,
			Dialer:      getDialer(),
		})
		readers = append(readers, reader)

		// 启动消费者 goroutine
		go c.consume(ctx, reader, ch, handler)
	}
	c.readers[topic] = readers

	return ch, nil
}

func (c *Consumer) consume(ctx context.Context, r *kafka.Reader, ch chan *Message, handler ConsumeFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Error.Printf("read message from kafka failed: %v", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
			continue
		}

		kafkaMsg := &Message{Key: msg.Key, Value: msg.Value}
		if handler != nil {
			if err := handler(kafkaMsg); err != nil {
				continue
			}
		}

		select {
		case ch <- kafkaMsg:
			if err := r.CommitMessages(ctx, msg); err != nil {
				log.Error.Printf("commit kafka message failed: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for topic, readers := range c.readers {
		for _, r := range readers {
			if err := r.Close(); err != nil {
				log.Error.Printf("close reader[%s] failed: %v", topic, err)
			}
		}
	}

	for topic, ch := range c.chans {
		close(ch)
		log.Info.Printf("closed consumer channel: %s", topic)
	}

	c.readers = make(map[string][]*kafka.Reader)
	c.chans = make(map[string]chan *Message)

	return nil
}
