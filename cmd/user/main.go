package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"vicomova/internal/shared/infrastructure/data/mysql"
	userEntity "vicomova/internal/user/domain/entity"
	"vicomova/internal/shared/pkg/log"
	"vicomova/internal/user/wire"
	"vicomova/pkg/config"
	"vicomova/pkg/etcd"

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

	// 初始化依赖
	p, err := wire.NewProvider(cfg)
	if err != nil {
		log.Error.Fatalf("Failed to init provider: %v", err)
	}

	// 自动迁移
	if err := mysql.GetDB().AutoMigrate(&userEntity.User{}); err != nil {
		log.Error.Fatalf("Failed to auto migrate: %v", err)
	}

	// 创建 Kitex Server
	svr := userService.NewServer(p.UserHandler)

	// 启动 RPC server 获取实际地址
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	log.Info.Printf("User RPC server starting on %s", addr)

	// 注册到 etcd
	var registry *etcd.Registry
	var configWatcher *config.ConfigWatcher
	if len(cfg.Etcd.Endpoints) > 0 {
		etcdClient, err := etcd.NewClient(&cfg.Etcd)
		if err != nil {
			log.Error.Fatalf("Failed to create etcd client: %v", err)
		}
		defer func() { _ = etcdClient.Close() }()

		registry = etcd.NewRegistry(etcdClient, "user", addr)
		if err := registry.Register(); err != nil {
			log.Error.Fatalf("Failed to register to etcd: %v", err)
		}
		defer func() { _ = registry.Unregister() }()
		defer func() { _ = svr.Stop() }()

		// 初始化配置热更新
		configWatcher = config.NewConfigWatcher(etcdClient, cfg)
		configWatcher.Watch("jwt.secret", func(oldVal, newVal interface{}) {
			log.Info.Printf("JWT secret changed: %v -> %v", oldVal, newVal)
		})
		configWatcher.Start(context.Background())
		if err := config.InitDefaultConfig(context.Background(), etcdClient, cfg); err != nil {
			log.Error.Fatalf("Failed to init default config: %v", err)
		}
	}

	// 优雅关闭
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		klog.Info("Shutting down server...")
		if registry != nil {
			_ = registry.Unregister()
		}
		_ = svr.Stop()
	}()

	if err := svr.Run(); err != nil {
		log.Error.Fatalf("Server error: %v", err)
	}
}