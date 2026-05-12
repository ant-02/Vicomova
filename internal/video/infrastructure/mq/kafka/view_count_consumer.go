package kafka

import (
	"context"
	"encoding/binary"
	"sync"
	"time"

	"vicomova/internal/video/domain/repository"
	pkgKafka "vicomova/pkg/infrastructure/kafka"
	"vicomova/pkg/log"
)

type ViewCountConsumer struct {
	consumer *pkgKafka.Consumer
	topic    string
	groupID  string
	repo     repository.VideoRepository

	// 窗口累积
	mu      sync.Mutex
	counter map[int64]int64 // videoID -> count

	// 控制
	ctx    context.Context
	cancel context.CancelFunc
}

func NewViewCountConsumer(topic, groupID string, repo repository.VideoRepository) *ViewCountConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &ViewCountConsumer{
		consumer: pkgKafka.NewConsumer(),
		topic:    topic,
		groupID:  groupID,
		repo:     repo,
		counter:  make(map[int64]int64),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start 启动消费
func (c *ViewCountConsumer) Start() error {
	ch, err := c.consumer.ConsumeTopic(
		c.ctx,
		pkgKafka.GetBrokers(),
		c.topic,
		c.groupID,
		1,    // consumerNum
		1000, // channel capacity
		c.onMessage,
	)
	if err != nil {
		return err
	}

	// 启动定时刷新
	go c.flushLoop()

	// 消费 channel（实际处理在上面 onMessage 已完成，这里只是防止 channel 溢出）
	go func() {
		for range ch {
		}
	}()

	log.Info.Printf("ViewCountConsumer started, topic=%s, group=%s", c.topic, c.groupID)
	return nil
}

// onMessage 消息回调，累积计数
func (c *ViewCountConsumer) onMessage(msg *pkgKafka.Message) error {
	videoID, n := binary.Varint(msg.Value)
	if n <= 0 {
		log.Error.Printf("ViewCountConsumer: decode failed, n=%d, data_len=%d", n, len(msg.Value))
		return nil // 不重试，跳过这条消息
	}

	c.mu.Lock()
	c.counter[videoID]++
	count := c.counter[videoID]
	shouldFlush := len(c.counter) >= FlushThreshold
	c.mu.Unlock()

	log.Debug.Printf("ViewCountConsumer: received videoID=%d, accumulative_count=%d", videoID, count)

	// 达到阈值立即刷新
	if shouldFlush {
		go c.flush()
	}

	return nil
}

// flushLoop 定时刷新
func (c *ViewCountConsumer) flushLoop() {
	ticker := time.NewTicker(FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			c.flush() // 最后一次刷新
			return
		case <-ticker.C:
			c.flush()
		}
	}
}

// flush 刷新累积到数据库
func (c *ViewCountConsumer) flush() {
	c.mu.Lock()
	if len(c.counter) == 0 {
		c.mu.Unlock()
		return
	}

	// 复制并清空
	batch := c.counter
	c.counter = make(map[int64]int64)
	c.mu.Unlock()

	log.Info.Printf("ViewCountConsumer: flushing %d videos to database", len(batch))

	// 批量更新数据库
	for videoID, count := range batch {
		if err := c.repo.IncrementViewBatch(context.Background(), videoID, count); err != nil {
			log.Error.Printf("ViewCountConsumer: increment failed videoID=%d: %v", videoID, err)
		}
	}

	log.Info.Printf("ViewCountConsumer: flushed successfully")
}

// Stop 停止消费
func (c *ViewCountConsumer) Stop() error {
	c.cancel()
	return c.consumer.Close()
}
