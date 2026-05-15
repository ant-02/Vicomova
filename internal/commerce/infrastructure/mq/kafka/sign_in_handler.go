package kafka

import (
	"log"
)

type SignInHandler struct{}

func NewSignInHandler() *SignInHandler {
	return &SignInHandler{}
}

func (h *SignInHandler) Start() error {
	// TODO: 使用 kafka.Init 后的 consumer 订阅 user-sign-in topic
	log.Printf("SignInHandler.Start: started")
	return nil
}

func (h *SignInHandler) Stop() error {
	log.Printf("SignInHandler.Stop: stopped")
	return nil
}
