package gateway

import (
	"fmt"

	interactionRpc "vicomova/internal/interaction/interfaces/grpc"
	interactionHandler "vicomova/internal/interaction/interfaces/http/handler"
	userRpc "vicomova/internal/user/interfaces/grpc"
	userHandler "vicomova/internal/user/interfaces/http/handler"
	videoRpc "vicomova/internal/video/interfaces/grpc"
	videoHandler "vicomova/internal/video/interfaces/http/handler"
	"vicomova/pkg/constants"
	"vicomova/pkg/etcd"
)

func NewUserClient(discovery *etcd.Discovery) (*userHandler.UserHandler, error) {
	addr := GetServiceAddr(discovery, constants.ServiceUser)
	client, err := userRpc.NewUserClient(constants.ServiceUser, addr)
	if err != nil {
		return nil, fmt.Errorf("create user client: %w", err)
	}
	return userHandler.NewUserHandler(client), nil
}

func NewVideoClient(discovery *etcd.Discovery) (*videoHandler.VideoHandler, error) {
	addr := GetServiceAddr(discovery, constants.ServiceVideo)
	client, err := videoRpc.NewVideoClient(constants.ServiceVideo, addr)
	if err != nil {
		return nil, fmt.Errorf("create video client: %w", err)
	}
	return videoHandler.NewVideoHandler(client), nil
}

func NewInteractionClient(discovery *etcd.Discovery) (*interactionHandler.InteractionHandler, error) {
	addr := GetServiceAddr(discovery, constants.ServiceInteraction)
	client, err := interactionRpc.NewInteractionClient(constants.ServiceInteraction, addr)
	if err != nil {
		return nil, fmt.Errorf("create interaction client: %w", err)
	}
	return interactionHandler.NewInteractionHandler(client), nil
}
