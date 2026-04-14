package gateway

import (
	"context"

	"vicomova/internal/shared/pkg/log"
	"vicomova/pkg/config"
	"vicomova/pkg/etcd"
)

type Bootstrap struct {
	Config       *config.Config
	EtcdClient   *etcd.Client
	Discovery    *etcd.Discovery
	ConfigWatcher *config.ConfigWatcher
}

func NewBootstrap(cfg *config.Config) (*Bootstrap, error) {
	b := &Bootstrap{
		Config: cfg,
	}

	// 尝试连接 etcd（可选，失败不影响启动）
	if len(cfg.Etcd.Endpoints) > 0 {
		cli, err := etcd.NewClient(&cfg.Etcd)
		if err != nil {
			log.Warn.Printf("Failed to connect to etcd: %v, continuing without etcd", err)
		} else {
			b.EtcdClient = cli
			b.Discovery = etcd.NewDiscovery(cli)
			b.ConfigWatcher = config.NewConfigWatcher(cli, cfg)

			// 初始化默认配置到 etcd
			if err := config.InitDefaultConfig(context.Background(), cli, cfg); err != nil {
				log.Warn.Printf("Failed to init default config: %v", err)
			}

			// 初始化默认路由到 etcd
			if err := InitRouteConfig(context.Background(), cli); err != nil {
				log.Warn.Printf("Failed to init default routes: %v", err)
			}

			log.Info.Printf("Connected to etcd at %v", cfg.Etcd.Endpoints)
		}
	}

	if b.EtcdClient == nil {
		log.Warn.Println("Running without etcd - using config file addresses")
	}

	return b, nil
}

// Close closes etcd connections
func (b *Bootstrap) Close() {
	if b.EtcdClient != nil {
		if b.Discovery != nil {
			b.Discovery.StopWatch()
		}
		b.EtcdClient.Close()
	}
}

// WatchConfig registers a config change callback
func (b *Bootstrap) WatchConfig(key string, callback func(oldVal, newVal interface{})) {
	if b.ConfigWatcher != nil {
		b.ConfigWatcher.Watch(key, callback)
		b.ConfigWatcher.Start(context.Background())
	}
}

// GetServiceAddr discovers a service address from etcd, falls back to config file
func (b *Bootstrap) GetServiceAddr(serviceName string) string {
	// 优先从 etcd 发现
	if b.Discovery != nil {
		addr, err := b.Discovery.GetOneInstance(context.Background(), serviceName)
		if err == nil {
			log.Info.Printf("Discovered service %s at %s from etcd", serviceName, addr)
			return addr
		}
		log.Warn.Printf("Service %s not found in etcd: %v", serviceName, err)
	}

	// fallback 到配置文件
	if addr, ok := b.Config.Services[serviceName]; ok && addr.Addr != "" {
		log.Info.Printf("Using config file address for service %s: %s", serviceName, addr.Addr)
		return addr.Addr
	}

	log.Error.Printf("No address found for service %s", serviceName)
	return ""
}

// Client returns the etcd client (nil if not connected)
func (b *Bootstrap) Client() *etcd.Client {
	return b.EtcdClient
}
