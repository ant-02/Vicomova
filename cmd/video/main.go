package main

import (
	"os"
	"os/signal"
	"syscall"

	"vicomova/internal/video/wire"
	"vicomova/pkg/config"
	"vicomova/pkg/constants"
	"vicomova/pkg/log"

	videoservice "vicomova/third_party/kitex_gen/video/videoservice"

	"github.com/cloudwego/kitex/pkg/klog"
)

func init() {
	config.Init(constants.ServiceVideo)
}

func main() {
	klog.SetLevel(klog.LevelInfo)

	cfg := config.Get()
	addr := cfg.Service.Addr

	p, err := wire.NewProvider()
	if err != nil {
		log.Error.Fatalf("Failed to init provider: %v", err)
	}

	svr := videoservice.NewServer(p.VideoHandler)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		klog.Info("Shutting down video service...")
		_ = svr.Stop()
	}()

	log.Info.Printf("Video service starting on %s", addr)
	if err := svr.Run(); err != nil {
		log.Error.Fatalf("Server error: %v", err)
	}
}