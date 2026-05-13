package etcd

import (
	"context"
	"fmt"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	cli       *clientv3.Client
	namespace string
	watcher   clientv3.Watcher
	watcherMu sync.RWMutex
}

func NewClient(endpoints []string, username, password, namespace string) (*Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		Username:    username,
		Password:    password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return &Client{
		cli:       cli,
		namespace: namespace,
	}, nil
}

func (c *Client) Close() error {
	c.watcherMu.Lock()
	if c.watcher != nil {
		_ = c.watcher.Close()
	}
	c.watcherMu.Unlock()
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

func (c *Client) Get(ctx context.Context, key string) (*clientv3.GetResponse, error) {
	return c.cli.Get(ctx, c.Key(key))
}

func (c *Client) Put(ctx context.Context, key, value string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	return c.cli.Put(ctx, c.Key(key), value, opts...)
}

func (c *Client) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	return c.cli.Delete(ctx, c.Key(key), opts...)
}

func (c *Client) Watch(ctx context.Context, prefix string, opts ...clientv3.OpOption) clientv3.WatchChan {
	c.watcherMu.Lock()
	defer c.watcherMu.Unlock()

	if c.watcher != nil {
		_ = c.watcher.Close()
	}
	c.watcher = clientv3.NewWatcher(c.cli)

	return c.watcher.Watch(ctx, c.Key(prefix), opts...)
}

func (c *Client) Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	return c.cli.Grant(ctx, ttl)
}

func (c *Client) KeepAlive(ctx context.Context, leaseID clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	return c.cli.KeepAlive(ctx, leaseID)
}

func (c *Client) GetInstances(ctx context.Context, prefix string) ([]string, error) {
	resp, err := c.cli.Get(ctx, c.Key(prefix), clientv3.WithPrefix())
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
