package kafka

import (
	"sync"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"

	"vicomova/pkg/log"
)

const (
	Timeout = 3 * time.Second
)

type Message struct {
	Key, Value []byte
}

var (
	initOnce     sync.Once
	brokers      []string
	saslUsername string
	saslPassword string

	// 单例
	transport     *kafka.Transport
	transportOnce sync.Once
	dialer        *kafka.Dialer
	dialerOnce    sync.Once
)

func Init(brokersList []string, username, password string) error {
	initOnce.Do(func() {
		brokers = brokersList
		saslUsername = username
		saslPassword = password
		log.Info.Printf("Kafka initialized: %v", brokers)
	})
	return nil
}

func GetBrokers() []string {
	return brokers
}

func getDialer() *kafka.Dialer {
	dialerOnce.Do(func() {
		dialer = &kafka.Dialer{
			Timeout:       Timeout,
			DualStack:     true,
			SASLMechanism: plain.Mechanism{Username: saslUsername, Password: saslPassword},
		}
	})
	return dialer
}

func getTransport() *kafka.Transport {
	transportOnce.Do(func() {
		transport = &kafka.Transport{
			SASL: plain.Mechanism{Username: saslUsername, Password: saslPassword},
		}
	})
	return transport
}
