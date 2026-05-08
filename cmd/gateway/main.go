package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"vicomova/internal/gateway"
	interactionRpc "vicomova/internal/interaction/interfaces/grpc"
	interactionHandler "vicomova/internal/interaction/interfaces/http/handler"
	"vicomova/internal/user/domain/service"
	rpc "vicomova/internal/user/interfaces/grpc"
	"vicomova/internal/user/interfaces/http/handler"
	"vicomova/internal/user/interfaces/http/router"
	videoRpc "vicomova/internal/video/interfaces/grpc"
	videoHandler "vicomova/internal/video/interfaces/http/handler"
	"vicomova/pkg/config"
	"vicomova/pkg/constants"
	hertz "vicomova/pkg/infrastructure/hertz"
	"vicomova/pkg/log"
)

func init() {
	config.Init(constants.ServiceGateway)
}

func main() {
	cfg := config.Get()
	addr := cfg.Service.Addr

	bs, err := gateway.NewBootstrap()
	if err != nil {
		log.Error.Fatalf("Failed to bootstrap: %v", err)
	}
	defer bs.Close()

	h := hertz.NewServer(addr)

	userAddr := bs.GetServiceAddr(constants.ServiceUser)
	userClient, err := rpc.NewUserClient(constants.ServiceUser, userAddr)
	if err != nil {
		log.Error.Fatalf("Failed to create user client: %v", err)
	}

	tokenSvc := service.NewTokenService(cfg.JWT.Secret)
	userHandler := handler.NewUserHandler(userClient)
	router.RegisterRoutes(h, userHandler, tokenSvc)

	videoAddr := bs.GetServiceAddr(constants.ServiceVideo)
	videoClient, err := videoRpc.NewVideoClient(constants.ServiceVideo, videoAddr)
	if err != nil {
		log.Error.Fatalf("Failed to create video client: %v", err)
	}

	interactionAddr := bs.GetServiceAddr(constants.ServiceInteraction)
	interactionClient, err := interactionRpc.NewInteractionClient(constants.ServiceInteraction, interactionAddr)
	if err != nil {
		log.Error.Fatalf("Failed to create interaction client: %v", err)
	}

	vh := videoHandler.NewVideoHandler(videoClient)
	ih := interactionHandler.NewInteractionHandler(interactionClient)
	gateway.RegisterVideoRoutes(h, vh, tokenSvc)
	gateway.RegisterInteractionRoutes(h, ih)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Info.Print("Shutting down gateway...")
		_ = h.Shutdown(context.Background())
	}()

	log.Info.Printf("Gateway starting on %s", addr)
	if err := hertz.Run(h); err != nil {
		log.Error.Fatalf("Gateway error: %v", err)
	}
}