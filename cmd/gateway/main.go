package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"vicomova/internal/gateway"
	"vicomova/pkg/config"
	"vicomova/pkg/constants"
	hertz "vicomova/pkg/infrastructure/hertz"
	"vicomova/pkg/log"
)

func init() {
	config.Init(constants.ServiceGateway)

	config.RegisterCallback(func(cfg *config.Config) {
		if cfg.JWT.Secret != config.Get().JWT.Secret {
			log.Warn.Printf("JWT secret changed, please restart gateway to take effect")
		}
		if cfg.Service.Addr != config.Get().Service.Addr {
			log.Warn.Printf("Service addr changed, please restart gateway to take effect")
		}
	})
}

func main() {
	h := hertz.NewServer(config.Get().Service.Addr)

	if err := gateway.RegisterRoutes(h); err != nil {
		log.Error.Fatalf("Failed to register routes: %v", err)
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Info.Print("Shutting down gateway...")
		_ = h.Shutdown(context.Background())
	}()

	log.Info.Printf("Gateway starting on %s", config.Get().Service.Addr)
	if err := hertz.Run(h); err != nil {
		log.Error.Fatalf("Gateway error: %v", err)
	}

	config.Close()
}
