package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"vicomova/internal/data/mysql"
	redisClient "vicomova/internal/data/redis"
	"vicomova/internal/pkg/log"
	"vicomova/internal/repository"
	"vicomova/internal/rpc/user"
	"vicomova/internal/service"
	"vicomova/pkg/config"

	userService "vicomova/third_party/kitex_gen/user/userservice"

	"github.com/cloudwego/kitex/pkg/klog"
)

var (
	configPath string
	port       int
)

func init() {
	flag.StringVar(&configPath, "config", "config/base.yaml", "config file path")
	flag.IntVar(&port, "port", 8888, "server port")
}

func main() {
	flag.Parse()

	// 初始化 Kitex 日志
	klog.SetLevel(klog.LevelInfo)

	// 加载配置
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Error.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := mysql.Init(&cfg.Database); err != nil {
		log.Error.Fatalf("Failed to init mysql: %v", err)
	}
	defer mysql.Close()

	// 初始化 Redis
	if err := redisClient.Init(&cfg.Redis); err != nil {
		log.Error.Fatalf("Failed to init redis: %v", err)
	}
	defer redisClient.Close()

	// 初始化 Repository
	userRepo := repository.NewUserRepository()

	// 初始化 Service
	userSvc := service.NewUserService(userRepo, "your-secret-key")

	// 初始化 Handler
	userHandler := user.NewUserHandler(userSvc)

	// 创建 Kitex Server
	svr := userService.NewServer(userHandler)

	// 优雅关闭
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		klog.Info("Shutting down server...")
		svr.Stop()
	}()

	addr := config.DefaultAddr("user", port)
	log.Info.Printf("User RPC server starting on %s", addr)
	if err := svr.Run(); err != nil {
		log.Error.Fatalf("Server error: %v", err)
	}
}
