package main

import (
	"os"
	"os/signal"
	"syscall"

	"vicomova/internal/user/wire"
	"vicomova/pkg/config"
	"vicomova/pkg/constants"
	"vicomova/pkg/log"

	userService "vicomova/third_party/kitex_gen/user/userservice"

	"github.com/cloudwego/kitex/pkg/klog"
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

	svr := userService.NewServer(p.UserHandler)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		klog.Info("Shutting down server...")
		_ = svr.Stop()
	}()

	log.Info.Printf("User RPC server starting on %s", addr)
	if err := svr.Run(); err != nil {
		log.Error.Fatalf("Server error: %v", err)
	}
}