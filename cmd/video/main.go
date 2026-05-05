package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"vicomova/internal/shared/pkg/log"
	"vicomova/internal/video/wire"
	"vicomova/pkg/config"
	"vicomova/pkg/etcd"

	videoservice "vicomova/third_party/kitex_gen/video/videoservice"

	"github.com/cloudwego/kitex/pkg/klog"
)

var (
	configPath string
	port       int
)

func init() {
	flag.StringVar(&configPath, "config", "config/base.yaml", "config file path")
	flag.IntVar(&port, "port", 8889, "video service port")
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
	svr := videoservice.NewServer(p.VideoHandler)

	// 注册到 etcd
	var registry *etcd.Registry
	if len(cfg.Etcd.Endpoints) > 0 {
		etcdClient, err := etcd.NewClient(&cfg.Etcd)
		if err != nil {
			log.Warn.Printf("Failed to create etcd client: %v", err)
		} else {
			defer func() { _ = etcdClient.Close() }()
			registry = etcd.NewRegistry(etcdClient, "video", fmt.Sprintf("127.0.0.1:%d", port))
			if err := registry.Register(); err != nil {
				log.Warn.Printf("Failed to register to etcd: %v", err)
			}
			defer func() { _ = registry.Unregister() }()
		}
	}

	// 优雅关闭
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		klog.Info("Shutting down video service...")
		if registry != nil {
			_ = registry.Unregister()
		}
		_ = svr.Stop()
	}()

	log.Info.Printf("Video service starting on 127.0.0.1:%d", port)
	if err := svr.Run(); err != nil {
		log.Error.Fatalf("Server error: %v", err)
	}
}