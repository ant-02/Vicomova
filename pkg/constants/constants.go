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

// Kafka Topic 名称
const (
	// KafkaTopicVideoView 视频播放量 topic
	KafkaTopicVideoView = "video-view"
)

// 服务配置 key
const (
	// ServiceKeyUser 用户服务配置 key
	ServiceKeyUser = "user"
)

// 服务名称
const (
	ServiceUser        = "user"
	ServiceVideo       = "video"
	ServiceInteraction = "interaction"
	ServiceGateway     = "gateway"
)

// Etcd config key
const (
	// EtcdConfigKey Etcd 配置 key（相对路径）
	EtcdConfigKey = "config"
)

// Etcd services key
const (
	// EtcdServicesKeyPrefix 服务注册 key 前缀
	EtcdServicesKeyPrefix = "services"
)

// HTTP Client 相关
const (
	// HTTPClientTimeout HTTP 客户端超时
	HTTPClientTimeout = 10 * time.Second
)

// Context Key
const (
	// ContextKeyUserID 用户ID上下文key
	ContextKeyUserID = "user_id"
)

// OSS 相关
const (
	// UploadTokenExpire 上传凭证过期时间（秒）
	UploadTokenExpire = 3600
)

// 上传类型，用于区分视频和封面在 OSS 存储的路径
const (
	UploadTokenTypeVideo = 1
	UploadTokenTypeCover = 2
)

// 数据库表名
const (
	TableVideos    = "videos"
	TableUsers     = "users"
	TableLikes     = "likes"
	TableComments  = "comments"
	TableFavorites = "favorites"
)

// 热门视频相关
const (
	DefaultHotVideoLimit = 20
	MaxHotVideoLimit     = 100
	WilsonZ              = 1.96 // Wilson 算法置信度参数
)
