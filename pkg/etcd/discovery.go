package etcd

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"vicomova/internal/shared/pkg/log"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Discovery struct {
	client      *Client
	instances   []string
	mu          sync.RWMutex
	index       uint32
	watchCancel context.CancelFunc
}

func NewDiscovery(client *Client) *Discovery {
	return &Discovery{
		client: client,
	}
}

func (d *Discovery) GetInstances(ctx context.Context, serviceName string) ([]string, error) {
	return d.client.GetInstances(ctx, serviceName)
}

func (d *Discovery) GetOneInstance(ctx context.Context, serviceName string) (string, error) {
	instances, err := d.client.GetInstances(ctx, serviceName)
	if err != nil {
		return "", err
	}
	if len(instances) == 0 {
		return "", fmt.Errorf("no instances found for service %s", serviceName)
	}
	// Round-robin select
	idx := atomic.AddUint32(&d.index, 1) - 1
	return instances[idx%uint32(len(instances))], nil
}

func (d *Discovery) WatchInstances(ctx context.Context, serviceName string, callback func(instances []string)) error {
	d.mu.Lock()
	ctx, d.watchCancel = context.WithCancel(ctx)
	d.mu.Unlock()

	// Initial fetch
	instances, err := d.client.GetInstances(ctx, serviceName)
	if err != nil {
		return err
	}
	d.mu.Lock()
	d.instances = instances
	d.mu.Unlock()
	callback(instances)

	// Watch for changes
	watcher := clientv3.NewWatcher(d.client.CLI())
	ch := watcher.Watch(ctx, d.client.ServicePath(serviceName), clientv3.WithPrefix(), clientv3.WithPrevKV())

	go func() {
		for {
			select {
			case resp, ok := <-ch:
				if !ok {
					log.Info.Printf("Watch channel closed for service %s", serviceName)
					return
				}
				if resp.Err() != nil {
					log.Error.Printf("Watch error for service %s: %v", serviceName, resp.Err())
					continue
				}

				// Re-fetch all instances on any change
				newInstances, err := d.client.GetInstances(ctx, serviceName)
				if err != nil {
					log.Error.Printf("Failed to get instances for %s: %v", serviceName, err)
					continue
				}

				d.mu.Lock()
				d.instances = newInstances
				d.mu.Unlock()

				log.Info.Printf("Service %s instances updated: %v", serviceName, newInstances)
				callback(newInstances)

			case <-ctx.Done():
				_ = watcher.Close()
				return
			}
		}
	}()

	return nil
}

func (d *Discovery) StopWatch() {
	d.mu.Lock()
	if d.watchCancel != nil {
		d.watchCancel()
	}
	d.mu.Unlock()
}

// LoadBalancer provides round-robin load balancing over discovered instances
type LoadBalancer struct {
	discovery   *Discovery
	serviceName string
	index       uint32
	mu          sync.Mutex
}

func NewLoadBalancer(discovery *Discovery, serviceName string) *LoadBalancer {
	return &LoadBalancer{
		discovery:   discovery,
		serviceName: serviceName,
	}
}

func (lb *LoadBalancer) GetNext(ctx context.Context) (string, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	instances, err := lb.discovery.GetInstances(ctx, lb.serviceName)
	if err != nil {
		return "", err
	}
	if len(instances) == 0 {
		return "", fmt.Errorf("no instances found for service %s", lb.serviceName)
	}

	idx := lb.index % uint32(len(instances))
	lb.index++
	return instances[idx], nil
}
