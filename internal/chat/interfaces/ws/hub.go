package ws

import (
	"sync"
)

// Hub maintains active WebSocket connections
type Hub struct {
	clients map[int64]*Client
	mu      sync.RWMutex
}

var (
	hub     *Hub
	hubOnce sync.Once
)

// GetHub returns the singleton Hub instance
func GetHub() *Hub {
	hubOnce.Do(func() {
		hub = &Hub{
			clients: make(map[int64]*Client),
		}
	})
	return hub
}

// AddClient registers a client connection
func (h *Hub) AddClient(userID int64, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = client
}

// RemoveClient unregisters a client connection
func (h *Hub) RemoveClient(userID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if c, ok := h.clients[userID]; ok {
		close(c.SendChan)
		delete(h.clients, userID)
	}
}

// GetClient retrieves a client by user ID
func (h *Hub) GetClient(userID int64) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients[userID]
}

// IsOnline checks if a user is connected
func (h *Hub) IsOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID int64, data []byte) bool {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		return false
	}

	client.Send(data)
	return true
}

// GetOnlineUsers returns list of online user IDs
func (h *Hub) GetOnlineUsers() []int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]int64, 0, len(h.clients))
	for userID := range h.clients {
		users = append(users, userID)
	}
	return users
}

// GetOnlineCount returns the number of connected clients
func (h *Hub) GetOnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
