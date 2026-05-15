package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP 请求总数
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTP 请求延迟
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// WebSocket 连接数
	WebSocketConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_connections_total",
			Help: "Current number of WebSocket connections",
		},
	)

	// WebSocket 消息数
	WebSocketMessagesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_total",
			Help: "Total number of WebSocket messages",
		},
		[]string{"type"}, // 1v1, group
	)

	// 聊天消息数
	ChatMessagesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chat_messages_total",
			Help: "Total number of chat messages",
		},
		[]string{"type"}, // 1v1, group
	)

	// Redis 连接数
	RedisConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "redis_connections_total",
			Help: "Current number of Redis connections",
		},
	)

	// MySQL 连接数
	MySQLConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "mysql_connections_total",
			Help: "Current number of MySQL connections",
		},
	)

	// 在线用户数
	OnlineUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "online_users_total",
			Help: "Current number of online users",
		},
	)

	// 关注/粉丝数
	FollowersTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "followers_total",
			Help: "Total number of followers",
		},
	)

	// 消息未读数
	UnreadMessagesTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "unread_messages_total",
			Help: "Total number of unread messages",
		},
	)
)
