package gateway

import (
	"context"

	interactionRouter "vicomova/internal/interaction/interfaces/http/router"
	userRouter "vicomova/internal/user/interfaces/http/router"
	videoRouter "vicomova/internal/video/interfaces/http/router"
	"vicomova/pkg/config"
	"vicomova/pkg/etcd"
	"vicomova/pkg/log"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz) error {
	discovery := etcd.NewDiscovery(config.GetClient())

	uh, err := NewUserClient(discovery)
	if err != nil {
		return err
	}
	userRouter.RegisterRoutes(h, uh)

	vh, err := NewVideoClient(discovery)
	if err != nil {
		return err
	}
	videoRouter.RegisterRoutes(h, vh)

	ih, err := NewInteractionClient(discovery)
	if err != nil {
		return err
	}
	interactionRouter.RegisterRoutes(h, ih)

	return nil
}

func GetServiceAddr(discovery *etcd.Discovery, serviceName string) string {
	addr, err := discovery.GetOneInstance(context.Background(), serviceName)
	if err != nil {
		log.Error.Fatalf("Failed to discover service %s: %v", serviceName, err)
	}
	log.Info.Printf("Discovered service %s at %s", serviceName, addr)
	return addr
}
