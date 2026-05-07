package config

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"vicomova/pkg/etcd"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdConfigLoader etcd 配置加载器
type EtcdConfigLoader struct {
	client      *clientv3.Client
	namespace   string // 配置前缀，如 "/vicomova/config"
	configMu   sync.RWMutex
	config      *Config
	onChange   func(*Config) // 配置变更回调
}

var (
	etcdLoader *EtcdConfigLoader
	etcdOnce   sync.Once
)

// LoadFromEtcd 从 etcd 加载配置
// namespace: 配置前缀，如 "/vicomova/config"
// onChange: 配置变更时的回调
func LoadFromEtcd(namespace string, onChange func(*Config)) (*Config, error) {
	var initErr error

	etcdOnce.Do(func() {
		// 从本地加载 etcd 连接配置
		localCfg, err := Load("")
		if err != nil {
			initErr = fmt.Errorf("failed to load local config for etcd: %w", err)
			return
		}

		// 连接 etcd
		client, err := etcd.NewClient(&localCfg.Etcd)
		if err != nil {
			initErr = fmt.Errorf("failed to connect to etcd: %w", err)
			return
		}

		// 创建加载器
		etcdLoader = &EtcdConfigLoader{
			client:    client,
			namespace: namespace,
			onChange:  onChange,
		}

		// 加载初始配置
		if err := etcdLoader.load(); err != nil {
			initErr = fmt.Errorf("failed to load config from etcd: %w", err)
			return
		}

		// 启动 watch
		go etcdLoader.watch()
	})

	if initErr != nil {
		return nil, initErr
	}

	return etcdLoader.config, nil
}

// load 从 etcd 加载配置
func (l *EtcdConfigLoader) load() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取所有配置
	resp, err := l.client.Get(ctx, l.namespace, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	if len(resp.Kvs) == 0 {
		return fmt.Errorf("no config found in etcd at %s", l.namespace)
	}

	// 构建配置 map
	configMap := make(map[string]interface{})
	for _, kv := range resp.Kvs {
		key := strings.TrimPrefix(string(kv.Key), l.namespace+"/")
		if key == "" {
			continue
		}

		// 尝试解析 JSON
		var value interface{}
		if err := json.Unmarshal(kv.Value, &value); err != nil {
			value = string(kv.Value)
		}
		setNestedValue(configMap, key, value)
	}

	// 转换为 JSON 再解析为 Config
	jsonData, err := json.Marshal(configMap)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	l.configMu.Lock()
	l.config = &cfg
	l.configMu.Unlock()

	return nil
}

// watch 监听配置变化
func (l *EtcdConfigLoader) watch() {
	for {
		ctx, cancel := context.WithCancel(context.Background())
		watchChan := l.client.Watch(ctx, l.namespace, clientv3.WithPrefix())

		for watchResp := range watchChan {
			if watchResp.Err() != nil {
				time.Sleep(2 * time.Second)
				break
			}

			for _, event := range watchResp.Events {
				key := strings.TrimPrefix(string(event.Kv.Key), l.namespace+"/")

				switch event.Type {
				case clientv3.EventTypePut:
					var value interface{}
					json.Unmarshal(event.Kv.Value, &value)
					l.updateKey(key, value)

				case clientv3.EventTypeDelete:
					l.updateKey(key, nil)
				}
			}

			// 触发回调
			if l.onChange != nil {
				l.configMu.RLock()
				cfg := l.config
				l.configMu.RUnlock()
				l.onChange(cfg)
			}
		}

		cancel()
	}
}

// updateKey 更新单个配置项
func (l *EtcdConfigLoader) updateKey(key string, value interface{}) {
	l.configMu.Lock()
	defer l.configMu.Unlock()

	if l.config == nil {
		return
	}

	// 简单的直接替换（实际应该更新对应的字段）
	parts := strings.Split(key, "/")
	if len(parts) < 2 {
		return
	}

	switch parts[0] {
	case "jwt":
		if len(parts) == 2 && parts[1] == "secret" {
			if s, ok := value.(string); ok {
				l.config.JWT.Secret = s
			}
		}
	case "redis":
		if len(parts) == 2 {
			switch parts[1] {
			case "host":
				if s, ok := value.(string); ok {
					l.config.Redis.Host = s
				}
			case "port":
				if f, ok := value.(float64); ok {
					l.config.Redis.Port = int(f)
				}
			}
		}
	case "database":
		// 类似处理...
	}
}

// GetConfig 返回当前配置（线程安全）
func GetConfig() *Config {
	if etcdLoader == nil {
		return nil
	}
	etcdLoader.configMu.RLock()
	defer etcdLoader.configMu.RUnlock()
	return etcdLoader.config
}

// setNestedValue 设置嵌套的配置值
func setNestedValue(m map[string]interface{}, key string, value interface{}) {
	parts := strings.Split(key, "/")
	current := m

	for i, part := range parts {
		if i == len(parts)-1 {
			current[part] = value
			return
		}

		if _, exists := current[part]; !exists {
			current[part] = make(map[string]interface{})
		}
		current = current[part].(map[string]interface{})
	}
}
