package kafka

import (
	"context"
	"encoding/json"
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

func (h *SignInHandler) handleSignIn(ctx context.Context, msg []byte) error {
	var event SignInEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		log.Printf("SignInHandler.handleSignIn: unmarshal failed: %v", err)
		return err
	}

	log.Printf("SignInHandler.handleSignIn: received sign-in event for user=%d, consecutive_days=%d",
		event.UserID, event.ConsecutiveDays)

	// TODO: 调用 Application 层的 SignInCommandService 处理签到积分发放
	// 这部分在 application/command/sign_in.go 中实现后连接

	return nil
}
