package kafka

import (
	"vicomova/pkg/log"
)

// Consumer 消费者接口
type Consumer interface {
	// Start 启动消费
	Start() error
	// Stop 停止消费
	Stop() error
}

// ConsumerManager 消费者管理器
type ConsumerManager struct {
	consumers []Consumer
}

func NewConsumerManager(consumers ...Consumer) *ConsumerManager {
	return &ConsumerManager{
		consumers: consumers,
	}
}

// Register 注册消费者（可变参数）
func (m *ConsumerManager) Register(consumers ...Consumer) {
	m.consumers = append(m.consumers, consumers...)
}

// StartAll 启动所有消费者
func (m *ConsumerManager) StartAll() error {
	for _, c := range m.consumers {
		if err := c.Start(); err != nil {
			log.Error.Printf("ConsumerManager: failed to start consumer: %v", err)
			return err
		}
	}
	log.Info.Printf("ConsumerManager: started %d consumers", len(m.consumers))
	return nil
}

// StopAll 停止所有消费者
func (m *ConsumerManager) StopAll() error {
	for _, c := range m.consumers {
		if err := c.Stop(); err != nil {
			log.Error.Printf("ConsumerManager: failed to stop consumer: %v", err)
		}
	}
	log.Info.Printf("ConsumerManager: stopped all consumers")
	return nil
}
