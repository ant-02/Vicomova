package kafka

import (
	"fmt"

	"vicomova/internal/shared/pkg/constants"
	"vicomova/internal/shared/pkg/log"
	"vicomova/pkg/config"

	"github.com/IBM/sarama"
)

var Producer sarama.SyncProducer

func InitProducer(cfg *config.KafkaConfig) error {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = constants.KafkaProducerRetryMax
	saramaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		return fmt.Errorf("failed to create kafka producer: %w", err)
	}

	Producer = producer
	log.Info.Printf("Kafka producer connected: %v", cfg.Brokers)
	return nil
}

func GetProducer() sarama.SyncProducer {
	return Producer
}

func SendMessage(topic string, key string, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err := Producer.SendMessage(msg)
	return err
}

func Close() error {
	if Producer != nil {
		return Producer.Close()
	}
	return nil
}
