package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"vicomova/pkg/constants"
	"vicomova/pkg/etcd"

	"gopkg.in/yaml.v3"
)

var (
	globalConfig *Config
	configMu     sync.RWMutex
	logger       *log.Logger
	client       *etcd.Client
	closeChan    chan struct{}
	callbacks    []ConfigChangeCallback
)

func init() {
	logger = log.New(os.Stdout, "[config] ", log.LstdFlags)
}

// Init initializes config from etcd
func Init(serviceName string) {
	endpoints := os.Getenv("ETCD_ENDPOINTS")
	if endpoints == "" {
		panic("ETCD_ENDPOINTS is required")
	}

	username := os.Getenv("ETCD_USERNAME")
	password := os.Getenv("ETCD_PASSWORD")
	namespace := os.Getenv("ETCD_NAMESPACE")

	cli, err := etcd.NewClient(strings.Split(endpoints, ","), username, password, namespace)
	if err != nil {
		panic(fmt.Sprintf("failed to create etcd client: %v", err))
	}
	client = cli

	cfg, err := loadConfigFromEtcd(client, constants.EtcdConfigKey, serviceName)
	if err != nil {
		panic(fmt.Sprintf("load config failed: %v", err))
	}

	configMu.Lock()
	globalConfig = cfg
	configMu.Unlock()

	logger.Printf("config initialized for service %s", serviceName)

	closeChan = make(chan struct{})
	go watchEtcdConfig(client, constants.EtcdConfigKey, serviceName)
}

func loadConfigFromEtcd(client *etcd.Client, path string, serviceName string) (*Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("no config found at %s", path)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(resp.Kvs[0].Value, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Match service from services map
	if cfg.Services != nil {
		if s, ok := cfg.Services[serviceName]; ok {
			cfg.Service = s
		}
	}

	return cfg, nil
}

func watchEtcdConfig(client *etcd.Client, path string, serviceName string) {
	for {
		select {
		case <-closeChan:
			return
		default:
		}

		watchChan := client.Watch(context.Background(), path)

		for watchResp := range watchChan {
			if watchResp.Err() != nil {
				logger.Printf("watch error: %v, reconnecting...", watchResp.Err())
				time.Sleep(2 * time.Second)
				break
			}

			for _, event := range watchResp.Events {
				if event.IsModify() || event.Type == 1 { // EventTypePut = 1
					logger.Printf("config changed: %s", event.Kv.Key)

					cfg, err := loadConfigFromEtcd(client, path, serviceName)
					if err != nil {
						logger.Printf("failed to reload config: %v", err)
						continue
					}

					configMu.Lock()
					globalConfig = cfg
					configMu.Unlock()

					logger.Printf("config reloaded successfully")

					// Notify callbacks
					for _, cb := range callbacks {
						cb(cfg)
					}
				}
			}
		}
	}
}

// RegisterCallback registers a function to be called when config changes
func RegisterCallback(cb ConfigChangeCallback) {
	callbacks = append(callbacks, cb)
}

// Get returns the global config
func Get() *Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return globalConfig
}

// Close closes the etcd client and stops watching
func Close() {
	if closeChan != nil {
		close(closeChan)
	}
	if client != nil {
		_ = client.Close()
	}
}

// GetClient returns the underlying etcd client
func GetClient() *etcd.Client {
	return client
}
