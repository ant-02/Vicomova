package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"vicomova/internal/chat/interfaces/ws"
	"vicomova/internal/chat/wire"
	"vicomova/pkg/config"
	"vicomova/pkg/constants"
	"vicomova/pkg/etcd"
	"vicomova/pkg/log"
	"vicomova/pkg/utils"

	chatservice "vicomova/third_party/kitex_gen/chat/chatservice"

	"github.com/cloudwego/kitex/pkg/klog"
	server "github.com/cloudwego/kitex/server"
)

func init() {
	config.Init(constants.ServiceChat)
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
	svr := chatservice.NewServer(p.ChatHandler, server.WithServiceAddr(tcpAddr))

	log.Info.Printf("Chat service starting on %s", addr)

	go func() {
		if err := svr.Run(); err != nil {
			log.Error.Fatalf("Server error: %v", err)
		}
	}()

	// Start WebSocket server on separate port
	wsAddr := cfg.Services["chat"].WebSocketAddr
	if wsAddr == "" {
		wsAddr = "127.0.0.1:8892" // default WebSocket port
	}

	tokenFn := func(token string) (int64, error) {
		claims, err := utils.ParseToken(token, cfg.JWT.Secret)
		if err != nil {
			return 0, err
		}
		return utils.GetUserIDFromClaims(claims)
	}

	chatSender := ws.NewChatSenderAdapter(p.ChatHandler.SendMessageWS)
	chatGroupSender := ws.NewChatGroupSenderAdapter(p.ChatHandler.SendGroupMessageWS)

	wsServer := ws.NewWSServer(wsAddr, tokenFn, chatSender, chatGroupSender)
	wsServer.Start()

	cli := config.GetClient()
	registry := etcd.NewRegistry(cli, constants.ServiceChat, addr)
	if err := registry.Register(); err != nil {
		log.Error.Fatalf("Failed to register service: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	_ = svr.Stop()
	wsServer.Stop()
	config.Close()
}
