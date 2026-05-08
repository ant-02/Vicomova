package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const configKey = "/vicomova/config"

var (
	globalConfig *Config
	configMu     sync.RWMutex
	logger       *log.Logger
	etcdClient   *clientv3.Client
	closeChan    chan struct{}
)

// ConfigChangeCallback is called when config changes
type ConfigChangeCallback func(*Config)

// callbacks for config change notifications
var callbacks []ConfigChangeCallback

func init() {
	logger = log.New(os.Stdout, "[config] ", log.LstdFlags)
}

// Init initializes config from etcd
func Init(serviceName string) {
	etcdAddr := os.Getenv("ETCD_ADDR")
	if etcdAddr == "" {
		logger.Fatal("ETCD_ADDR is required")
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdAddr, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logger.Fatalf("connect to etcd failed: %v", err)
	}
	etcdClient = client

	cfg, err := loadConfigFromEtcd(client, configKey, serviceName)
	if err != nil {
		logger.Fatalf("load config failed: %v", err)
	}

	configMu.Lock()
	globalConfig = cfg
	configMu.Unlock()

	logger.Printf("config initialized for service %s", serviceName)

	closeChan = make(chan struct{})
	go watchEtcdConfig(client, configKey, serviceName)
}

func loadConfigFromEtcd(client *clientv3.Client, path string, serviceName string) (*Config, error) {
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
	if err := json.Unmarshal(resp.Kvs[0].Value, cfg); err != nil {
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

func watchEtcdConfig(client *clientv3.Client, path string, serviceName string) {
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
				if event.IsModify() || event.Type == clientv3.EventTypePut {
					logger.Printf("config changed: %s", event.Kv.Key)

					cfg, err := loadConfigFromEtcd(client, path, serviceName)
					if err != nil {
						logger.Printf("failed to reload config: %v", err)
						continue
					}

					configMu.Lock()
					oldConfig := globalConfig
					globalConfig = cfg
					configMu.Unlock()

					logger.Printf("config reloaded successfully")

					// Notify callbacks
					for _, cb := range callbacks {
						cb(cfg)
					}

					_ = oldConfig // silence unused variable
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
	if etcdClient != nil {
		etcdClient.Close()
	}
}
