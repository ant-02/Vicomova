package ws

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// ChatSenderInterface for sending 1v1 messages via WebSocket
type ChatSenderInterface interface {
	SendMessage(ctx context.Context, senderID, receiverID int64, content string) error
}

// ChatGroupSenderInterface for sending group messages via WebSocket
type ChatGroupSenderInterface interface {
	SendGroupMessage(ctx context.Context, senderID, groupID int64, content string) error
}

// chatSenderAdapter implements ChatSenderInterface
type chatSenderAdapter struct {
	fn func(ctx context.Context, senderID, receiverID int64, content string) error
}

func (a *chatSenderAdapter) SendMessage(ctx context.Context, senderID, receiverID int64, content string) error {
	return a.fn(ctx, senderID, receiverID, content)
}

// NewChatSenderAdapter creates an adapter for ChatHandler methods
func NewChatSenderAdapter(
	sendMessageFn func(ctx context.Context, senderID, receiverID int64, content string) error,
) ChatSenderInterface {
	return &chatSenderAdapter{fn: sendMessageFn}
}

// chatGroupSenderAdapter implements ChatGroupSenderInterface
type chatGroupSenderAdapter struct {
	fn func(ctx context.Context, senderID, groupID int64, content string) error
}

func (a *chatGroupSenderAdapter) SendGroupMessage(ctx context.Context, senderID, groupID int64, content string) error {
	return a.fn(ctx, senderID, groupID, content)
}

// NewChatGroupSenderAdapter creates an adapter for ChatHandler methods
func NewChatGroupSenderAdapter(
	sendGroupMessageFn func(ctx context.Context, senderID, groupID int64, content string) error,
) ChatGroupSenderInterface {
	return &chatGroupSenderAdapter{fn: sendGroupMessageFn}
}

// WSToken WebSocket 连接 Token 结构
// @Summary WebSocket Token 结构
// @Description 临时 WebSocket 连接 token，包含用户和目标信息
// @Tags websocket
type WSToken struct {
	// @Description 用户ID
	UserID int64 `json:"user_id"`
	// @Description 目标ID（peer_id 或 group_id）
	TargetID int64 `json:"target_id"`
	// @Description 连接类型（1v1 或 group）
	Type string `json:"type"`
	// @Description 过期时间戳（秒）
	ExpiresAt int64 `json:"expires_at"`
	// @Description 签名
	Signature string `json:"signature"`
}

// GenerateWSToken 生成临时 WebSocket 连接 Token
// @Summary 生成 WebSocket Token
// @Description 生成一个临时 token 用于 WebSocket 连接验证（5分钟有效期）
// @Tags websocket
// @Param user_id query int64 true "用户ID"
// @Param target_id query int64 true "目标用户ID或群ID"
// @Param type query string true "连接类型：1v1 或 group"
// @Param secret query string true "签名密钥"
// @Success 200 {object} map[string]interface{} "返回 token 和过期时间"
// @Router /ws/token [post]
func GenerateWSToken(userID, targetID int64, tokenType, secret string) (string, int64, error) {
	expiresAt := time.Now().Add(5 * time.Minute).Unix() // 5分钟后过期

	// 构建签名字符串
	data := fmt.Sprintf("%d:%d:%s:%d", userID, targetID, tokenType, expiresAt)

	// 使用 SHA256 生成签名
	h := sha256.New()
	h.Write([]byte(data + secret))
	signature := hex.EncodeToString(h.Sum(nil))

	// 拼接 token
	token := fmt.Sprintf("%d:%d:%s:%d:%s", userID, targetID, tokenType, expiresAt, signature)

	return token, expiresAt, nil
}

// ValidateWSToken validates and parses a WebSocket token
func ValidateWSToken(token, secret string) (*WSToken, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 5 {
		return nil, fmt.Errorf("invalid token format")
	}

	var userID, targetID, expiresAt int64
	var tokenType string

	fmt.Sscanf(parts[0], "%d", &userID)
	fmt.Sscanf(parts[1], "%d", &targetID)
	tokenType = parts[2]
	fmt.Sscanf(parts[3], "%d", &expiresAt)
	signature := parts[4]

	// 检查是否过期
	if time.Now().Unix() > expiresAt {
		return nil, fmt.Errorf("token expired")
	}

	// 验证签名
	data := fmt.Sprintf("%d:%d:%s:%d", userID, targetID, tokenType, expiresAt)
	h := sha256.New()
	h.Write([]byte(data + secret))
	expectedSig := hex.EncodeToString(h.Sum(nil))

	if signature != expectedSig {
		return nil, fmt.Errorf("invalid signature")
	}

	return &WSToken{
		UserID:    userID,
		TargetID:  targetID,
		Type:      tokenType,
		ExpiresAt: expiresAt,
	}, nil
}
