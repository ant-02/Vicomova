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
	"vicomova/internal/shared/pkg/log"
	"vicomova/pkg/config"
	hertz "vicomova/pkg/hertz"
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

	// Bootstrap - 初始化 etcd 连接（可选，无 etcd 时使用默认地址）
	bs, err := gateway.NewBootstrap(cfg)
	if err != nil {
		log.Error.Fatalf("Failed to bootstrap: %v", err)
	}
	defer bs.Close()

	// 监听配置变更
	bs.WatchConfig("jwt.secret", func(oldVal, newVal interface{}) {
		log.Info.Printf("JWT secret changed: %v -> %v", oldVal, newVal)
	})

	// 创建 Hertz 服务器
	h := hertz.NewServer(port)

	// 初始化 User RPC Client（通过 bootstrap 发现地址）
	userAddr := bs.GetServiceAddr("user")
	if userAddr == "" {
		log.Error.Fatalf("No address found for user service")
	}
	userClient, err := rpc.NewUserClient("user", userAddr)
	if err != nil {
		log.Error.Fatalf("Failed to create user client: %v", err)
	}

	// 创建 handler 并注册路由
	userHandler := handler.NewUserHandler(userClient)
	router.RegisterRoutes(h, userHandler)

	// 初始化 etcd 路由管理器（用于热更新路由配置）
	if bs.Client() != nil {
		routeMgr := gateway.NewDynamicRouter(bs.Client(), h)
		routeMgr.LoadRoutes(context.Background())
		routeMgr.WatchRoutes(context.Background())
	}

	// 优雅关闭
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Info.Print("Shutting down gateway...")
		h.Shutdown(context.Background())
	}()

	log.Info.Printf("Gateway starting on :%d", port)
	if err := hertz.Run(h); err != nil {
		log.Error.Fatalf("Gateway error: %v", err)
	}
}
