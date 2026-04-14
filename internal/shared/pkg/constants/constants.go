package constants

import "time"

// Token 相关
const (
	// AccessTokenExpiry Access Token 过期时间
	AccessTokenExpiry = 30 * time.Minute
	// RefreshTokenExpiry Refresh Token 过期时间
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

// Email Code 相关
const (
	// EmailCodeTTL 邮箱验证码过期时间
	EmailCodeTTL = 10 * time.Minute
	// EmailCodeLength 邮箱验证码长度
	EmailCodeLength = 6
)

// Etcd 相关
const (
	// EtcdDialTimeout Etcd 拨号超时
	EtcdDialTimeout = 5 * time.Second
	// EtcdLeaseTTL 服务注册 Lease 过期时间
	EtcdLeaseTTL = 10 * time.Second
	// EtcdKeepAliveTimeout Lease 续约超时
	EtcdKeepAliveTimeout = 2 * time.Second
)

// Kafka 相关
const (
	// KafkaProducerRetryMax Producer 最大重试次数
	KafkaProducerRetryMax = 3
)

// HTTP Client 相关
const (
	// HTTPClientTimeout HTTP 客户端超时
	HTTPClientTimeout = 10 * time.Second
)
