package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"vicomova/internal/data/mysql"
	domainUser "vicomova/internal/domain/user"
	"vicomova/internal/pkg/log"
	"vicomova/internal/wire"
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

	// 初始化依赖
	p, err := wire.NewProvider(cfg)
	if err != nil {
		log.Error.Fatalf("Failed to init provider: %v", err)
	}

//

	// 自动迁移
	mysql.GetDB().AutoMigrate(&domainUser.User{})

	// 创建 Kitex Server
	svr := userService.NewServer(p.UserHandler)

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
