package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"vicomova/internal/interaction/wire"
	"vicomova/internal/shared/pkg/log"
	"vicomova/pkg/config"
	"vicomova/pkg/etcd"

	interactionservice "vicomova/third_party/kitex_gen/interaction/interactionservice"

	"github.com/cloudwego/kitex/pkg/klog"
)

var (
	configPath string
	port       int
)

func init() {
	flag.StringVar(&configPath, "config", "config/base.yaml", "config file path")
	flag.IntVar(&port, "port", 8890, "interaction service port")
}

func main() {
	flag.Parse()

	klog.SetLevel(klog.LevelInfo)

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Error.Fatalf("Failed to load config: %v", err)
	}

	// 初始化依赖
	p, err := wire.NewProvider(cfg)
	if err != nil {
		log.Error.Fatalf("Failed to init provider: %v", err)
	}

	// 创建 Kitex server
	svr := interactionservice.NewServer(p.InteractionHandler)

	// 注册到 etcd
	var registry *etcd.Registry
	if len(cfg.Etcd.Endpoints) > 0 {
		etcdClient, err := etcd.NewClient(&cfg.Etcd)
		if err != nil {
			log.Info.Printf("Failed to create etcd client: %v", err)
		} else {
			defer func() { _ = etcdClient.Close() }()
			defer func() { _ = registry.Unregister() }()
			if err := registry.Register(); err != nil {
				log.Info.Printf("Failed to register to etcd: %v", err)
			}
			defer registry.Unregister()
		}
	}

	// 优雅关闭
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		klog.Info("Shutting down interaction service...")
		if registry != nil {
			_ = registry.Unregister()
		}
		_ = svr.Stop()
	}()

	log.Info.Printf("Interaction service starting on 127.0.0.1:%d", port)
	if err := svr.Run(); err != nil {
		log.Error.Fatalf("Server error: %v", err)
	}
}
