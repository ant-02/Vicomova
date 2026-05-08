package gateway

import (
	"context"
	"fmt"
	"os"
	"strings"

	"vicomova/pkg/config"
	"vicomova/pkg/infrastructure/etcd"
	"vicomova/pkg/log"
)

type Bootstrap struct {
	EtcdClient *etcd.Client
	Discovery  *etcd.Discovery
}

func NewBootstrap() (*Bootstrap, error) {
	b := &Bootstrap{}

	etcdAddr := os.Getenv("ETCD_ADDR")
	if etcdAddr == "" {
		return nil, fmt.Errorf("ETCD_ADDR is required")
	}

	cli, err := etcd.NewClient(&config.Etcd{
		Endpoints: strings.Split(etcdAddr, ","),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to etcd: %w", err)
	}
	b.EtcdClient = cli
	b.Discovery = etcd.NewDiscovery(cli)

	log.Info.Printf("Gateway connected to etcd at %v", etcdAddr)
	return b, nil
}

func (b *Bootstrap) Close() {
	if b.EtcdClient != nil {
		_ = b.EtcdClient.Close()
	}
}

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

func (b *Bootstrap) Client() *etcd.Client {
	return b.EtcdClient
}
