package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"vicomova/pkg/log"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSServer WebSocket server using standard net/http
type WSServer struct {
	server          *http.Server
	hub             *Hub
	addr            string
	tokenFn         func(token string) (int64, error)
	chatSender      ChatSenderInterface
	chatGroupSender ChatGroupSenderInterface
}

var wsServerInstance *WSServer

// NewWSServer creates a new WebSocket server
// @Summary 创建 WebSocket 服务器
// @Description 启动 WebSocket 服务器用于实时聊天（消息收发）
// @Tags websocket
// @Param addr query string true "服务器地址"
// @Param tokenFn query string false "token验证函数"
// @Param chatSender query string false "1v1消息发送器"
// @Param chatGroupSender query string false "群聊消息发送器"
// @Success 200 {object} WSServer
// @Router /ws [post]
func NewWSServer(addr string, tokenFn func(token string) (int64, error), chatSender ChatSenderInterface, chatGroupSender ChatGroupSenderInterface) *WSServer {
	return &WSServer{
		addr:            addr,
		tokenFn:         tokenFn,
		hub:             GetHub(),
		chatSender:      chatSender,
		chatGroupSender: chatGroupSender,
	}
}

// Start begins listening for WebSocket connections
func (s *WSServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/chat", s.wsHandler)

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	wsServerInstance = s

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error.Fatalf("WebSocket server failed: %v", err)
		}
	}()

	log.Info.Printf("WebSocket server started at %s", s.addr)
}

// Stop gracefully shuts down the WebSocket server
func (s *WSServer) Stop() error {
	if s.server != nil {
		return s.server.Shutdown(context.Background())
	}
	return nil
}

func (s *WSServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token required", http.StatusUnauthorized)
		return
	}

	userID, err := s.tokenFn(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := NewClient(userID, conn)
	client.Start()

	log.Info.Printf("WebSocket client connected: userID=%d", userID)
}

// IsUserOnline checks if a user is connected
func IsUserOnline(userID int64) bool {
	return GetHub().IsOnline(userID)
}

// PushMessageToUser sends a message to a specific online user
func PushMessageToUser(userID int64, fromUserID int64, content string, msgID int64) {
	hub := GetHub()
	msg := WSMessage{
		Type:       MsgTypeChat,
		FromUserID: fromUserID,
		Content:    content,
		Timestamp:  msgID,
	}

	if client := hub.GetClient(userID); client != nil {
		data, _ := json.Marshal(msg)
		client.Send(data)
	}
}

// PushGroupMessage broadcasts a message to all group members
func PushGroupMessage(groupID int64, fromUserID int64, content string, msgID int64, memberIDs []int64) {
	hub := GetHub()
	msg := WSMessage{
		Type:       MsgTypeChat,
		FromUserID: fromUserID,
		ToGroupID:  groupID,
		Content:    content,
		Timestamp:  msgID,
	}

	data, _ := json.Marshal(msg)
	for _, memberID := range memberIDs {
		if memberID != fromUserID {
			if client := hub.GetClient(memberID); client != nil {
				client.Send(data)
			}
		}
	}
}

// WSMessage WebSocket 消息结构
// @Summary WebSocket 消息格式
// @Description 客户端与服务端交换的 JSON 消息格式
// @Tags websocket
// @Example {"type":3,"from_user_id":123,"to_user_id":456,"content":"你好","client_msg_id":"uuid-xxx"}
type WSMessage struct {
	// @Description 消息类型：1=Ping（心跳）, 2=Pong（心跳响应）, 3=1v1聊天, 4=确认（Ack）, 5=系统消息, 6=群聊
	// @Enum 1,2,3,4,5,6
	Type int32 `json:"type"`
	// @Description 发送者用户ID（1v1和群聊必填）
	FromUserID int64 `json:"from_user_id,omitempty"`
	// @Description 接收者用户ID（1v1聊天必填）
	ToUserID int64 `json:"to_user_id,omitempty"`
	// @Description 目标群ID（群聊必填）
	ToGroupID int64 `json:"to_group_id,omitempty"`
	// @Description 消息内容
	Content string `json:"content,omitempty"`
	// @Description 服务端返回的时间戳（毫秒）
	Timestamp int64 `json:"timestamp,omitempty"`
	// @Description 客户端消息ID（用于确认收到）
	ClientMsgID string `json:"client_msg_id,omitempty"`
}

// Message types
const (
	MsgTypePing      = 1
	MsgTypePong      = 2
	MsgTypeChat      = 3
	MsgTypeAck       = 4
	MsgTypeSystem    = 5
	MsgTypeGroupChat = 6 // 群聊消息
)

// Client represents a WebSocket client
type Client struct {
	UserID   int64
	Conn     *websocket.Conn
	SendChan chan []byte
	hub      *Hub
}

func NewClient(userID int64, conn *websocket.Conn) *Client {
	return &Client{
		UserID:   userID,
		Conn:     conn,
		SendChan: make(chan []byte, 256),
		hub:      GetHub(),
	}
}

func (c *Client) Start() {
	c.hub.AddClient(c.UserID, c)
	defer c.hub.RemoveClient(c.UserID)

	go c.writePump()
	go c.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.Conn.Close()
		close(c.SendChan)
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				break
			}
			break
		}

		c.handleMessage(message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.SendChan:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(data []byte) {
	var msg WSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	switch msg.Type {
	case MsgTypePing:
		c.SendPong()
	case MsgTypeChat:
		// 1v1 聊天消息
		if wsServerInstance != nil && wsServerInstance.chatSender != nil {
			if err := wsServerInstance.chatSender.SendMessage(context.Background(), msg.FromUserID, msg.ToUserID, msg.Content); err != nil {
				log.Error.Printf("Failed to send message via WebSocket: %v", err)
				return
			}
		}
		c.SendAck(msg.ClientMsgID)
	case MsgTypeGroupChat:
		// 群聊消息
		if wsServerInstance != nil && wsServerInstance.chatGroupSender != nil {
			if err := wsServerInstance.chatGroupSender.SendGroupMessage(context.Background(), msg.FromUserID, msg.ToGroupID, msg.Content); err != nil {
				log.Error.Printf("Failed to send group message via WebSocket: %v", err)
				return
			}
		}
		c.SendAck(msg.ClientMsgID)
	}
}

func (c *Client) Send(data []byte) {
	select {
	case c.SendChan <- data:
	default:
		c.hub.RemoveClient(c.UserID)
	}
}

func (c *Client) SendPong() {
	c.Send(marshalWSMessage(WSMessage{Type: MsgTypePong}))
}

func (c *Client) SendAck(clientMsgID string) {
	c.Send(marshalWSMessage(WSMessage{
		Type:        MsgTypeAck,
		ClientMsgID: clientMsgID,
	}))
}

func marshalWSMessage(msg WSMessage) []byte {
	data, _ := json.Marshal(msg)
	return data
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024
)
