package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"vicomova/internal/shared/pkg/log"
	"vicomova/pkg/etcd"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// RouteConfig represents route configuration stored in etcd
type RouteConfig struct {
	Services []ServiceRoute `json:"services"`
}

// ServiceRoute describes routes for a single service
type ServiceRoute struct {
	Name    string       `json:"name"`
	Prefix  string       `json:"prefix"`
	Addr    string       `json:"addr"`
	Routes  []RouteInfo `json:"routes"`
}

// RouteInfo describes a single route
type RouteInfo struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"handler"` // e.g., "Login", "GetUser"
}

// DynamicRouter manages routes from etcd
type DynamicRouter struct {
	etcdClient *etcd.Client
	h          *server.Hertz
	routes     RouteConfig
	mu         sync.RWMutex
}

// NewDynamicRouter creates a new DynamicRouter
func NewDynamicRouter(etcdClient *etcd.Client, h *server.Hertz) *DynamicRouter {
	return &DynamicRouter{
		etcdClient: etcdClient,
		h:          h,
	}
}

// LoadRoutes loads route config from etcd
func (dr *DynamicRouter) LoadRoutes(ctx context.Context) error {
	if dr.etcdClient == nil {
		log.Warn.Println("DynamicRouter: etcd not available, no routes loaded")
		return nil
	}

	resp, err := dr.etcdClient.Get(ctx, "routes")
	if err != nil {
		return fmt.Errorf("failed to get routes from etcd: %w", err)
	}

	if len(resp.Kvs) == 0 {
		log.Warn.Println("DynamicRouter: no routes configured in etcd")
		return nil
	}

	var routes RouteConfig
	if err := json.Unmarshal(resp.Kvs[0].Value, &routes); err != nil {
		return fmt.Errorf("failed to unmarshal routes: %w", err)
	}

	dr.mu.Lock()
	dr.routes = routes
	dr.mu.Unlock()

	log.Info.Printf("DynamicRouter: loaded %d services from etcd", len(routes.Services))
	return nil
}

// GetRoutes returns current route configuration
func (dr *DynamicRouter) GetRoutes() RouteConfig {
	dr.mu.RLock()
	defer dr.mu.RUnlock()
	return dr.routes
}

// InitRouteConfig initializes default route config in etcd if not exists
func InitRouteConfig(ctx context.Context, etcdClient *etcd.Client) error {
	if etcdClient == nil {
		return nil
	}

	resp, err := etcdClient.Get(ctx, "routes")
	if err != nil {
		return err
	}
	if len(resp.Kvs) > 0 {
		return nil // already exists
	}

	// Default user service routes
	defaultRoutes := RouteConfig{
		Services: []ServiceRoute{
			{
				Name:   "user",
				Prefix: "/user",
				Addr:   "127.0.0.1:8888",
				Routes: []RouteInfo{
					{Method: "POST", Path: "/register/send", Handler: "SendVerificationCode"},
					{Method: "POST", Path: "/register/verify", Handler: "VerifyAndRegister"},
					{Method: "POST", Path: "/login", Handler: "Login"},
					{Method: "POST", Path: "/refresh", Handler: "RefreshToken"},
					{Method: "POST", Path: "/logout", Handler: "Logout"},
					{Method: "GET", Path: "/:id", Handler: "GetUser"},
				},
			},
		},
	}

	data, err := json.Marshal(defaultRoutes)
	if err != nil {
		return err
	}

	_, err = etcdClient.Put(ctx, "routes", string(data))
	if err != nil {
		return fmt.Errorf("failed to put default routes: %w", err)
	}

	log.Info.Println("DynamicRouter: initialized default routes in etcd")
	return nil
}

// WatchRoutes watches for route configuration changes
func (dr *DynamicRouter) WatchRoutes(ctx context.Context) {
	if dr.etcdClient == nil {
		return
	}

	go func() {
		ch := dr.etcdClient.Watch(ctx, dr.etcdClient.RoutesPath())

		for {
			select {
			case resp, ok := <-ch:
				if !ok {
					log.Info.Println("DynamicRouter: route watch channel closed")
					return
				}
				if resp.Err() != nil {
					log.Error.Printf("DynamicRouter: watch error: %v", resp.Err())
					continue
				}

				log.Info.Println("DynamicRouter: route config changed, reloading...")
				if err := dr.LoadRoutes(ctx); err != nil {
					log.Error.Printf("DynamicRouter: failed to reload routes: %v", err)
					continue
				}

				// Note: Full hot reload of routes requires rebuilding handlers
				// This is a simplified version - routes take effect on restart
				log.Info.Println("DynamicRouter: routes reloaded (restart required for full effect)")

			case <-ctx.Done():
				return
			}
		}
	}()
}
