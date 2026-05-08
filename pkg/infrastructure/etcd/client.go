package etcd

import (
	"context"
	"fmt"

	"vicomova/pkg/config"
	"vicomova/pkg/constants"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	cli       *clientv3.Client
	namespace string
}

func NewClient(cfg *config.Etcd) (*Client, error) {
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

func (c *Client) Key(key string) string {
	if c.namespace == "" {
		return key
	}
	return fmt.Sprintf("/%s/%s", c.namespace, key)
}

func (c *Client) ServicesPath() string {
	return c.Key("services")
}

func (c *Client) ServicePath(serviceName string) string {
	return fmt.Sprintf("%s/%s", c.ServicesPath(), serviceName)
}

func (c *Client) RoutesPath() string {
	return c.Key("routes")
}

func (c *Client) Get(ctx context.Context, key string) (*clientv3.GetResponse, error) {
	return c.cli.Get(ctx, c.Key(key))
}

func (c *Client) Put(ctx context.Context, key, value string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	return c.cli.Put(ctx, c.Key(key), value, opts...)
}

func (c *Client) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	return c.cli.Delete(ctx, c.Key(key), opts...)
}

func (c *Client) Watch(ctx context.Context, prefix string) clientv3.WatchChan {
	watcher := clientv3.NewWatcher(c.cli)
	return watcher.Watch(ctx, c.Key(prefix))
}

func (c *Client) Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	return c.cli.Grant(ctx, ttl)
}

func (c *Client) KeepAlive(ctx context.Context, leaseID clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	return c.cli.KeepAlive(ctx, leaseID)
}

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
