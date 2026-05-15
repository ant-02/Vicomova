package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	searchhttp "vicomova/internal/search/interfaces/http"
	"vicomova/internal/search/wire"
	"vicomova/pkg/config"
	"vicomova/pkg/constants"
	"vicomova/pkg/etcd"
	"vicomova/pkg/log"

	hertz "vicomova/pkg/infrastructure/hertz"
)

func init() {
	config.Init(constants.ServiceSearch)
}

func main() {
	log.Info.Printf("Search service starting...")

	p, err := wire.NewProvider()
	if err != nil {
		log.Error.Fatalf("Failed to init provider: %v", err)
	}

	addr := config.Get().Service.Addr
	h := hertz.NewServer(addr)

	searchHandler := searchhttp.NewSearchHandler(p.SearchService)
	searchhttp.RegisterRoutes(h, searchHandler)

	go func() {
		log.Info.Printf("Search HTTP server starting on %s", addr)
		if err := hertz.Run(h); err != nil {
			log.Error.Fatalf("Server error: %v", err)
		}
	}()

	tcpAddr, _ := net.ResolveTCPAddr("tcp", addr)
	_ = tcpAddr

	cli := config.GetClient()
	registry := etcd.NewRegistry(cli, constants.ServiceSearch, addr)
	if err := registry.Register(); err != nil {
		log.Error.Fatalf("Failed to register service: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	_ = h.Shutdown(context.Background())
	config.Close()
	log.Info.Printf("Search service stopped")
}
