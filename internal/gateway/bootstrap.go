package gateway

import (
	"context"

	"vicomova/internal/shared/pkg/log"
	"vicomova/pkg/config"
	"vicomova/pkg/etcd"
)

type Bootstrap struct {
	Config     *config.Config
	EtcdClient *etcd.Client
	Discovery  *etcd.Discovery
}

func NewBootstrap(cfg *config.Config) (*Bootstrap, error) {
	b := &Bootstrap{
		Config: cfg,
	}

	// 连接 etcd 进行服务发现
	if len(cfg.Etcd.Endpoints) > 0 {
		cli, err := etcd.NewClient(&cfg.Etcd)
		if err != nil {
			log.Error.Fatalf("Failed to connect to etcd: %v", err)
		}
		b.EtcdClient = cli
		b.Discovery = etcd.NewDiscovery(cli)

		log.Info.Printf("Connected to etcd at %v", cfg.Etcd.Endpoints)
	} else {
		log.Error.Fatalf("No etcd endpoints configured")
	}

	return b, nil
}

// Close closes etcd connections
func (b *Bootstrap) Close() {
	if b.EtcdClient != nil {
		b.EtcdClient.Close()
	}
}

// GetServiceAddr discovers a service address from etcd
func (b *Bootstrap) GetServiceAddr(serviceName string) string {
	if b.Discovery == nil {
		log.Error.Fatalf("etcd not connected, cannot discover service %s", serviceName)
	}

	addr, err := b.Discovery.GetOneInstance(context.Background(), serviceName)
	if err != nil {
		log.Error.Fatalf("Failed to discover service %s: %v", serviceName, err)
	}

	log.Info.Printf("Discovered service %s at %s", serviceName, addr)
	return addr
}

// Client returns the etcd client
func (b *Bootstrap) Client() *etcd.Client {
	return b.EtcdClient
}
