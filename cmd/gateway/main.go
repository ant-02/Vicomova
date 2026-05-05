package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"vicomova/internal/gateway"
	"vicomova/internal/user/interfaces/http/handler"
	"vicomova/internal/user/interfaces/http/router"
	rpc "vicomova/internal/user/interfaces/grpc"
	videoRpc "vicomova/internal/video/interfaces/grpc"
	interactionRpc "vicomova/internal/interaction/interfaces/grpc"
	videoHandler "vicomova/internal/video/interfaces/http/handler"
	interactionHandler "vicomova/internal/interaction/interfaces/http/handler"
	"vicomova/internal/shared/pkg/log"
	"vicomova/pkg/config"
	hertz "vicomova/pkg/hertz"
	"vicomova/internal/user/domain/service"
)

var (
	configPath string
	port       int
)

func init() {
	flag.StringVar(&configPath, "config", "config/base.yaml", "config file path")
	flag.IntVar(&port, "port", 8080, "gateway port")
}

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Error.Fatalf("Failed to load config: %v", err)
	}

	// Bootstrap - 初始化 etcd 连接
	bs, err := gateway.NewBootstrap(cfg)
	if err != nil {
		log.Error.Fatalf("Failed to bootstrap: %v", err)
	}
	defer bs.Close()

	// 创建 Hertz 服务器
	h := hertz.NewServer(port)

	// 初始化 User RPC Client（通过 bootstrap 发现地址）
	userAddr := bs.GetServiceAddr("user")
	userClient, err := rpc.NewUserClient("user", userAddr)
	if err != nil {
		log.Error.Fatalf("Failed to create user client: %v", err)
	}

	// 创建 handler 并注册路由
	tokenSvc := service.NewTokenService(cfg.JWT.Secret)
	userHandler := handler.NewUserHandler(userClient)
	router.RegisterRoutes(h, userHandler, tokenSvc)

	// 初始化 Video RPC Client
	videoAddr := bs.GetServiceAddr("video")
	videoClient, err := videoRpc.NewVideoClient("video", videoAddr)
	if err != nil {
		log.Error.Fatalf("Failed to create video client: %v", err)
	}

	// 初始化 Interaction RPC Client
	interactionAddr := bs.GetServiceAddr("interaction")
	interactionClient, err := interactionRpc.NewInteractionClient("interaction", interactionAddr)
	if err != nil {
		log.Error.Fatalf("Failed to create interaction client: %v", err)
	}

	// 创建 video/interaction handlers 并注册路由
	vh := videoHandler.NewVideoHandler(videoClient)
	ih := interactionHandler.NewInteractionHandler(interactionClient)
	gateway.RegisterVideoRoutes(h, vh, tokenSvc)
	gateway.RegisterInteractionRoutes(h, ih)

	// 优雅关闭
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Info.Print("Shutting down gateway...")
		_ = h.Shutdown(context.Background())
	}()

	log.Info.Printf("Gateway starting on :%d", port)
	if err := hertz.Run(h); err != nil {
		log.Error.Fatalf("Gateway error: %v", err)
	}
}
