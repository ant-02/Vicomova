package kitex

import (
	"fmt"
	"vicomova/internal/pkg/log"
)

type Server interface {
	Run() error
	Stop()
}

type KitexServer struct {
	svr interface {
		Run() error
		Stop()
	}
	addr string
}

func NewServer(serviceName string, addr string) *KitexServer {
	return &KitexServer{
		addr: addr,
	}
}

func (s *KitexServer) Run() error {
	log.Info.Printf("Kitex server starting on %s", s.addr)
	return nil
}

func (s *KitexServer) Stop() {
	log.Info.Println("Kitex server stopped")
}

func DefaultAddr(serviceName string, port int) string {
	return fmt.Sprintf("%s:%d", serviceName, port)
}

func GetFreePort() int {
	// Simple port allocation, actual implementation would check availability
	return 8888
}
