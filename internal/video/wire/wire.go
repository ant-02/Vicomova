package wire

import (
	"context"

	usergrpc "vicomova/internal/user/interfaces/grpc"
	"vicomova/internal/video/application/command"
	"vicomova/internal/video/application/query"
	"vicomova/internal/video/domain/repository"
	infraKafka "vicomova/internal/video/infrastructure/mq/kafka"
	infraMysql "vicomova/internal/video/infrastructure/persistence/mysql"
	redisCache "vicomova/internal/video/infrastructure/persistence/redis"
	videogrpc "vicomova/internal/video/interfaces/grpc"
	"vicomova/pkg/config"
	"vicomova/pkg/constants"
	"vicomova/pkg/infrastructure/kafka"
	"vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/infrastructure/oss"
	"vicomova/pkg/infrastructure/redis"
	"vicomova/pkg/log"
)

type Provider struct {
	MySQL           *mysql.Client
	Redis           *redis.Client
	VideoHandler    *videogrpc.VideoHandler
	ConsumerManager *infraKafka.ConsumerManager
}

func NewProvider() (*Provider, error) {
	cfg := config.Get()

	if err := mysql.Init(&cfg.Database); err != nil {
		return nil, err
	}

	if err := redis.Init(&cfg.Redis); err != nil {
		return nil, err
	}

	mysqlClient := mysql.GetClient()
	redisClient := redis.GetClient()

	if err := mysql.GetDB().AutoMigrate(&infraMysql.VideoPO{}); err != nil {
		return nil, err
	}
	if err := mysql.GetDB().AutoMigrate(&infraMysql.CategoryPO{}); err != nil {
		return nil, err
	}

	videoCache := redisCache.NewVideoCache(redisClient)
	hotVideoCache := redisCache.NewHotVideoCache(redisClient)
	categoryVideoCache := redisCache.NewCategoryVideoCache(redisClient)
	videoRepo := infraMysql.NewVideoRepository(mysqlClient)

	// User Service 客户端
	var userClient *usergrpc.UserClient
	if userCfg, ok := cfg.Services[constants.ServiceKeyUser]; ok && userCfg.Addr != "" {
		var err error
		userClient, err = usergrpc.NewUserClient(userCfg.Name, userCfg.Addr)
		if err != nil {
			log.Warn.Printf("failed to create user client: %v", err)
		}
	}

	// Kafka 播放量组件
	var viewCountProducer repository.ViewCountProducer
	var videoIndexProducer *infraKafka.VideoIndexProducer
	var consumerManager *infraKafka.ConsumerManager

	if len(cfg.Kafka.Brokers) > 0 && len(cfg.Kafka.Topics) > 0 {
		if err := kafka.Init(cfg.Kafka.Brokers, cfg.Kafka.SASL.Username, cfg.Kafka.SASL.Password); err != nil {
			return nil, err
		}

		var consumers []infraKafka.Consumer

		// 从配置获取 video-view topic
		if topicCfg, ok := cfg.Kafka.Topics[constants.KafkaTopicVideoView]; ok {
			sender := kafka.GetSender()
			producer := infraKafka.NewViewCountProducer(sender, topicCfg.Name)
			viewCountProducer = producer

			consumer := infraKafka.NewViewCountConsumer(topicCfg.Name, topicCfg.Group, videoRepo)
			consumers = append(consumers, consumer)
		}

		// video-index topic 用于 Search 服务索引视频
		if topicCfg, ok := cfg.Kafka.Topics[constants.KafkaTopicVideoIndex]; ok {
			sender := kafka.GetSender()
			videoIndexProducer = infraKafka.NewVideoIndexProducer(sender, topicCfg.Name)
		}

		consumerManager = infraKafka.NewConsumerManager(consumers...)
		if err := consumerManager.StartAll(); err != nil {
			return nil, err
		}
	}

	// 初始化 OSS（可选，未配置时跳过）
	var ossClient oss.OSS
	if cfg.OSS.Qiniu.AccessKey != "" {
		ossCfg := &oss.OSSConfig{
			AccessKey:  cfg.OSS.Qiniu.AccessKey,
			SecretKey:  cfg.OSS.Qiniu.SecretKey,
			Bucket:     cfg.OSS.Qiniu.Bucket,
			Domain:     cfg.OSS.Qiniu.Domain,
			UploadHost: cfg.OSS.Qiniu.UploadHost,
		}
		var err error
		ossClient, err = oss.NewOSS(oss.OSSTypeQiniu, ossCfg)
		if err != nil {
			return nil, err
		}
	}

	cmdSvc := command.NewVideoCommandService(videoRepo, videoCache, ossClient)
	querySvc := query.NewVideoQueryService(videoRepo, videoCache, ossClient, viewCountProducer, hotVideoCache, categoryVideoCache, userClient)

	// 设置视频索引 Producer
	if videoIndexProducer != nil {
		cmdSvc.SetIndexProducer(videoIndexProducer)
	}

	// 启动时预热热门视频缓存
	if err := querySvc.WarmUp(context.Background()); err != nil {
		log.Warn.Printf("hot cache warm up failed, will fallback on first request: %v", err)
	}

	videoHandler := videogrpc.NewVideoHandler(cmdSvc, querySvc, viewCountProducer)

	config.RegisterCallback(func(newCfg *config.Config) {
		if err := mysql.Reload(&newCfg.Database); err != nil {
			log.Error.Printf("failed to reload mysql: %v", err)
		}
		if err := redis.Reload(&newCfg.Redis); err != nil {
			log.Error.Printf("failed to reload redis: %v", err)
		}
	})

	return &Provider{
		MySQL:           mysqlClient,
		Redis:           redis.GetClient(),
		VideoHandler:    videoHandler,
		ConsumerManager: consumerManager,
	}, nil
}
