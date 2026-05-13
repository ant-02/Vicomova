package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vicomova/internal/user/wire"
	"vicomova/pkg/config"
	"vicomova/pkg/constants"
	"vicomova/pkg/etcd"
	"vicomova/pkg/log"

	userService "vicomova/third_party/kitex_gen/user/userservice"

	"github.com/cloudwego/kitex/pkg/klog"
	server "github.com/cloudwego/kitex/server"
)

func init() {
	config.Init(constants.ServiceUser)
}

func main() {
	klog.SetLevel(klog.LevelInfo)

	cfg := config.Get()
	addr := cfg.Service.Addr

	p, err := wire.NewProvider()
	if err != nil {
		log.Error.Fatalf("Failed to init provider: %v", err)
	}

	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	svr := userService.NewServer(p.UserHandler, server.WithServiceAddr(tcpAddr))

	log.Info.Printf("User RPC server starting on %s", addr)

	go func() {
		if err := svr.Run(); err != nil {
			log.Error.Fatalf("Server error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	cli := config.GetClient()
	registry := etcd.NewRegistry(cli, constants.ServiceUser, addr)
	if err := registry.Register(); err != nil {
		log.Error.Fatalf("Failed to register service: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	_ = svr.Stop()
	config.Close()
}
