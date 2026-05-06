package config

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"vicomova/internal/shared/pkg/log"
	"vicomova/pkg/etcd"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ChangeCallback func(oldVal, newVal interface{})

type ConfigWatcher struct {
	etcdClient  *etcd.Client
	config      *Config
	subscribers map[string][]ChangeCallback
	mu          sync.RWMutex
}

func NewConfigWatcher(etcdClient *etcd.Client, cfg *Config) *ConfigWatcher {
	return &ConfigWatcher{
		etcdClient:  etcdClient,
		config:      cfg,
		subscribers: make(map[string][]ChangeCallback),
	}
}

// Watch registers a callback for config changes
func (w *ConfigWatcher) Watch(key string, callback ChangeCallback) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.subscribers[key] = append(w.subscribers[key], callback)
	log.Info.Printf("Registered watcher for config key: %s", key)
}

// Get returns current value for a key (e.g., "jwt.secret")
func (w *ConfigWatcher) Get(key string) interface{} {
	w.mu.RLock()
	defer w.mu.RUnlock()

	parts := strings.Split(key, ".")
	if len(parts) < 2 {
		return nil
	}

	switch parts[0] {
	case "jwt":
		if parts[1] == "secret" {
			return w.config.JWT.Secret
		}
	case "redis":
		if parts[1] == "host" {
			return w.config.Redis.Host
		}
		if parts[1] == "port" {
			return w.config.Redis.Port
		}
	case "database":
		return w.config.Database
	case "server":
		return w.config.Server
	}
	return nil
}

// Start begins watching etcd for config changes
func (w *ConfigWatcher) Start(ctx context.Context) {
	if w.etcdClient == nil {
		log.Warn.Println("ConfigWatcher: etcd client not available, hot reload disabled")
		return
	}

	go w.watchLoop(ctx)
	log.Info.Println("ConfigWatcher started")
}

func (w *ConfigWatcher) watchLoop(ctx context.Context) {
	watcher := clientv3.NewWatcher(w.etcdClient.CLI())
	configPrefix := w.etcdClient.Key("config")

	ch := watcher.Watch(ctx, configPrefix, clientv3.WithPrefix())

	for {
		select {
		case resp, ok := <-ch:
			if !ok {
				log.Info.Println("ConfigWatcher: watch channel closed")
				return
			}
			if resp.Err() != nil {
				log.Error.Printf("ConfigWatcher: watch error: %v", resp.Err())
				continue
			}

			for _, event := range resp.Events {
				if event.Type == clientv3.EventTypePut {
					w.handleConfigChange(event.Kv.Key, event.Kv.Value)
				}
			}

		case <-ctx.Done():
			watcher.Close()
			return
		}
	}
}

func (w *ConfigWatcher) handleConfigChange(key []byte, value []byte) {
	configKey := w.extractKey(string(key))
	if configKey == "" {
		return
	}

	w.mu.RLock()
	callbacks := w.subscribers[configKey]
	w.mu.RUnlock()

	if len(callbacks) == 0 {
		return
	}

	log.Info.Printf("Config changed: %s", configKey)

	// Parse new value
	var newVal interface{}
	if err := json.Unmarshal(value, &newVal); err != nil {
		log.Error.Printf("Failed to unmarshal config value for %s: %v", configKey, err)
		return
	}

	// Notify all subscribers
	for _, cb := range callbacks {
		oldVal := w.Get(configKey)
		cb(oldVal, newVal)
	}

	// Update local config
	w.updateLocalConfig(configKey, value)
}

func (w *ConfigWatcher) extractKey(etcdKey string) string {
	// etcd key is like "/vicomova/config/jwt.secret"
	// extract "jwt.secret"
	prefix := w.etcdClient.Key("config")
	if !strings.HasPrefix(etcdKey, prefix) {
		return ""
	}
	key := strings.TrimPrefix(etcdKey, prefix)
	return strings.TrimPrefix(key, "/")
}

func (w *ConfigWatcher) updateLocalConfig(key string, value []byte) {
	w.mu.Lock()
	defer w.mu.Unlock()

	parts := strings.Split(key, ".")
	if len(parts) < 2 {
		return
	}

	switch parts[0] {
	case "jwt":
		if parts[1] == "secret" {
			w.config.JWT.Secret = string(value)
			log.Info.Printf("Updated local config: jwt.secret")
		}
	case "redis":
		if parts[1] == "host" {
			w.config.Redis.Host = string(value)
		} else if parts[1] == "port" {
			fmt.Sscanf(string(value), "%d", &w.config.Redis.Port)
		}
	}
}

// InitDefaultConfig writes default config to etcd if not exists
func InitDefaultConfig(ctx context.Context, etcdClient *etcd.Client, cfg *Config) error {
	if etcdClient == nil {
		return nil
	}

	// Initialize jwt.secret if not exists
	resp, err := etcdClient.Get(ctx, "config/jwt/secret")
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		_, err = etcdClient.Put(ctx, "config/jwt/secret", cfg.JWT.Secret)
		if err != nil {
			return fmt.Errorf("failed to init jwt.secret: %w", err)
		}
		log.Info.Println("Initialized jwt.secret in etcd")
	}

	// Initialize rate limit config if not exists
	resp, err = etcdClient.Get(ctx, "config/rate_limit")
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		defaultRateLimit := map[string]interface{}{
			"enabled":             true,
			"requests_per_second": 100,
		}
		data, _ := json.Marshal(defaultRateLimit)
		_, err = etcdClient.Put(ctx, "config/rate_limit", string(data))
		if err != nil {
			return fmt.Errorf("failed to init rate_limit: %w", err)
		}
		log.Info.Println("Initialized rate_limit in etcd")
	}

	return nil
}
