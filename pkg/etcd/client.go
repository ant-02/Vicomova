package etcd

import (
	"context"
	"fmt"

	"vicomova/internal/shared/pkg/constants"
	"vicomova/internal/shared/pkg/log"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	cli       *clientv3.Client
	namespace string
}

func NewClient(cfg *EtcdConfig) (*Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		Username:    cfg.Username,
		Password:    cfg.Password,
		DialTimeout: constants.EtcdDialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return &Client{
		cli:       cli,
		namespace: cfg.Namespace,
	}, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) CLI() *clientv3.Client {
	return c.cli
}

// Key appends namespace prefix to key
func (c *Client) Key(key string) string {
	if c.namespace == "" {
		return key
	}
	return fmt.Sprintf("/%s/%s", c.namespace, key)
}

// ServicesPath returns the prefix path for service registration
func (c *Client) ServicesPath() string {
	return c.Key("services")
}

// ServicePath returns the path for a specific service
func (c *Client) ServicePath(serviceName string) string {
	return fmt.Sprintf("%s/%s", c.ServicesPath(), serviceName)
}

// RoutesPath returns the path for routes configuration
func (c *Client) RoutesPath() string {
	return c.Key("routes")
}

// Get retrieves a value from etcd
func (c *Client) Get(ctx context.Context, key string) (*clientv3.GetResponse, error) {
	return c.cli.Get(ctx, c.Key(key))
}

// Put puts a key-value pair to etcd
func (c *Client) Put(ctx context.Context, key, value string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	return c.cli.Put(ctx, c.Key(key), value, opts...)
}

// Delete deletes a key from etcd
func (c *Client) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	return c.cli.Delete(ctx, c.Key(key), opts...)
}

// Watch watches a key prefix in etcd and returns the watch channel
func (c *Client) Watch(ctx context.Context, prefix string) clientv3.WatchChan {
	watcher := clientv3.NewWatcher(c.cli)
	return watcher.Watch(ctx, c.Key(prefix))
}

// Grant creates a lease with the given TTL
func (c *Client) Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	return c.cli.Grant(ctx, ttl)
}

// KeepAlive keeps a lease alive
func (c *Client) KeepAlive(ctx context.Context, leaseID clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	return c.cli.KeepAlive(ctx, leaseID)
}

// GetInstances gets all instance addresses for a service
func (c *Client) GetInstances(ctx context.Context, serviceName string) ([]string, error) {
	resp, err := c.cli.Get(ctx, c.ServicePath(serviceName), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	addrs := make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		if len(kv.Value) > 0 {
			addrs = append(addrs, string(kv.Value))
		}
	}
	return addrs, nil
}

// RegisterService registers a service instance with lease
func (c *Client) RegisterService(ctx context.Context, serviceName, instanceID, addr string, ttl int64) error {
	// Create lease
	leaseResp, err := c.Grant(ctx, ttl)
	if err != nil {
		return fmt.Errorf("failed to grant lease: %w", err)
	}

	// Register with lease
	key := fmt.Sprintf("%s/%s", c.ServicePath(serviceName), instanceID)
	_, err = c.Put(ctx, key, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	log.Info.Printf("Registered service %s/%s at %s with lease %d", serviceName, instanceID, addr, leaseResp.ID)
	return nil
}

// UnregisterService unregisters a service instance
func (c *Client) UnregisterService(ctx context.Context, serviceName, instanceID string) error {
	key := fmt.Sprintf("%s/%s", c.ServicePath(serviceName), instanceID)
	_, err := c.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to unregister service: %w", err)
	}

	log.Info.Printf("Unregistered service %s/%s", serviceName, instanceID)
	return nil
}

// WatchServiceInstances watches for changes in service instances
func (c *Client) WatchServiceInstances(ctx context.Context, serviceName string) clientv3.WatchChan {
	watcher := clientv3.NewWatcher(c.cli)
	r := watcher.Watch(ctx, c.ServicePath(serviceName), clientv3.WithPrefix(), clientv3.WithPrevKV())
	go func() {
		for {
			_, ok := <-r
			if !ok {
				log.Info.Printf("Watch channel closed for service %s", serviceName)
				return
			}
		}
	}()
	return r
}
