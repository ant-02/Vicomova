package etcd

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vicomova/pkg/constants"
	"vicomova/pkg/log"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Registry struct {
	client      *Client
	instanceID  string
	serviceName string
	addr        string
	ttl         int64
	leaseID     clientv3.LeaseID
	ctx         context.Context
	cancel      context.CancelFunc
	mu          sync.RWMutex
	stopCh      chan struct{}
}

type RegistryOption func(*Registry)

func WithTTL(ttl int64) RegistryOption {
	return func(r *Registry) {
		r.ttl = ttl
	}
}

func WithInstanceID(id string) RegistryOption {
	return func(r *Registry) {
		r.instanceID = id
	}
}

func NewRegistry(client *Client, serviceName, addr string, opts ...RegistryOption) *Registry {
	r := &Registry{
		client:      client,
		serviceName: serviceName,
		addr:        addr,
		instanceID:  fmt.Sprintf("%d", time.Now().UnixNano()),
		ttl:         int64(constants.EtcdLeaseTTL.Seconds()),
		stopCh:      make(chan struct{}),
	}

	for _, opt := range opts {
		opt(r)
	}

	r.ctx, r.cancel = context.WithCancel(context.Background())
	return r
}

// Key returns the etcd key for this service instance
func (r *Registry) Key() string {
	return fmt.Sprintf("%s/%s/%s", constants.EtcdServicesKeyPrefix, r.serviceName, r.instanceID)
}

func (r *Registry) Register() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Debug.Printf("[Registry] Starting registration for %s at %s", r.serviceName, r.addr)

	// Create lease
	leaseResp, err := r.client.Grant(r.ctx, r.ttl)
	if err != nil {
		return fmt.Errorf("failed to grant lease: %w", err)
	}
	r.leaseID = leaseResp.ID
	log.Debug.Printf("[Registry] Lease granted: %d, TTL: %d", r.leaseID, r.ttl)

	// Register service with lease
	key := r.Key()
	log.Debug.Printf("[Registry] Putting key=%s value=%s with lease=%d", key, r.addr, r.leaseID)

	putResp, err := r.client.Put(r.ctx, key, r.addr, clientv3.WithLease(r.leaseID))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	log.Debug.Printf("[Registry] Put succeeded, revision: %d", putResp.Header.Revision)
	log.Info.Printf("Registered service %s/%s at %s with lease %d", r.serviceName, r.instanceID, r.addr, r.leaseID)

	// Start keepalive goroutine
	go r.keepAlive()

	return nil
}

func (r *Registry) keepAlive() {
	log.Debug.Printf("[Registry] keepAlive: started")
	ticker := time.NewTicker(constants.EtcdLeaseTTL / 2)
	defer ticker.Stop()
	defer log.Debug.Printf("[Registry] keepAlive: exited")

	for {
		select {
		case <-ticker.C:
			r.mu.RLock()
			leaseID := r.leaseID
			r.mu.RUnlock()

			if leaseID == 0 {
				log.Debug.Printf("[Registry] keepAlive: leaseID is 0, skipping")
				continue
			}

			log.Debug.Printf("[Registry] keepAlive: sending KeepAlive for lease %d", leaseID)

			ctx, cancel := context.WithTimeout(r.ctx, constants.EtcdKeepAliveTimeout)
			ch, err := r.client.KeepAlive(ctx, leaseID)
			if err != nil {
				cancel()
				log.Error.Printf("Failed to keep alive lease %d: %v", leaseID, err)
				continue
			}

			// KeepAlive returns ONE response per request, wait for it then close
			resp, ok := <-ch
			if !ok {
				log.Error.Printf("Keep alive channel closed for lease %d", leaseID)
				cancel()
				return
			}
			log.Debug.Printf("[Registry] KeepAlive response: lease=%d, TTL=%d", resp.ID, resp.TTL)
			cancel()

		case <-r.stopCh:
			log.Debug.Printf("[Registry] keepAlive: received stop signal")
			return
		}
	}
}

func (r *Registry) Unregister() error {
	r.cancel()
	close(r.stopCh)

	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.client.Delete(r.ctx, r.Key())
	if err != nil {
		return fmt.Errorf("failed to unregister service: %w", err)
	}

	log.Info.Printf("Unregistered service %s/%s", r.serviceName, r.instanceID)
	return nil
}

func (r *Registry) ServiceName() string {
	return r.serviceName
}

func (r *Registry) InstanceID() string {
	return r.instanceID
}
