package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

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
	userAddr   string
)

func init() {
	flag.StringVar(&configPath, "config", "config/base.yaml", "config file path")
	flag.IntVar(&port, "port", 8080, "gateway port")
	flag.StringVar(&userAddr, "user_addr", "127.0.0.1:8888", "user service address")
}

func main() {
	flag.Parse()

	// 加载配置
	_, err := config.Load(configPath)
	if err != nil {
		log.Error.Fatalf("Failed to load config: %v", err)
	}

	// 初始化 User RPC Client
	userClient, err := rpc.NewUserClient("user", userAddr)
	if err != nil {
		log.Error.Fatalf("Failed to create user client: %v", err)
	}

	// 创建 Hertz 服务器
	h := hertz.NewServer(port)

	// 创建 handler
	userHandler := handler.NewUserHandler(userClient)

	// 注册路由
	router.RegisterRoutes(h, userHandler)

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