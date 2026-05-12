package wire

import (
	"vicomova/internal/video/application/command"
	"vicomova/internal/video/application/query"
	domainService "vicomova/internal/video/domain/service"
	infraKafka "vicomova/internal/video/infrastructure/mq/kafka"
	infraMysql "vicomova/internal/video/infrastructure/persistence/mysql"
	redisCache "vicomova/internal/video/infrastructure/persistence/redis"
	infraService "vicomova/internal/video/infrastructure/service"
	"vicomova/internal/video/interfaces/grpc"
	"vicomova/pkg/config"
	"vicomova/pkg/infrastructure/kafka"
	"vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/infrastructure/oss"
	"vicomova/pkg/infrastructure/redis"
	"vicomova/pkg/log"
)

type Provider struct {
	MySQL        *mysql.Client
	Redis        *redis.Client
	VideoHandler *grpc.VideoHandler
	Consumer     *infraKafka.ViewCountConsumer
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

	videoCache := redisCache.NewVideoCache(redisClient)
	videoRepo := infraMysql.NewVideoRepository(mysqlClient)

	hotAlgo := infraService.NewWilsonHotAlgorithm(videoRepo)

	// Kafka 播放量组件
	var viewCountService domainService.ViewCountService
	var viewCountConsumer *infraKafka.ViewCountConsumer

	if len(cfg.Kafka.Brokers) > 0 && len(cfg.Kafka.Topics) > 0 {
		if err := kafka.Init(cfg.Kafka.Brokers, cfg.Kafka.SASL.Username, cfg.Kafka.SASL.Password); err != nil {
			return nil, err
		}

		// 从配置获取 video-view topic
		if topicCfg, ok := cfg.Kafka.Topics["video-view"]; ok {
			producer := infraKafka.NewViewCountProducer(topicCfg.Name)
			viewCountService = infraKafka.NewViewCountService(producer)
			viewCountConsumer = infraKafka.NewViewCountConsumer(topicCfg.Name, topicCfg.Group, videoRepo)
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
	querySvc := query.NewVideoQueryService(videoRepo, videoCache, hotAlgo, ossClient, viewCountService)

	videoHandler := grpc.NewVideoHandler(cmdSvc, querySvc)

	config.RegisterCallback(func(newCfg *config.Config) {
		if err := mysql.Reload(&newCfg.Database); err != nil {
			log.Error.Printf("failed to reload mysql: %v", err)
		}
		if err := redis.Reload(&newCfg.Redis); err != nil {
			log.Error.Printf("failed to reload redis: %v", err)
		}
	})

	return &Provider{
		MySQL:        mysqlClient,
		Redis:        redis.GetClient(),
		VideoHandler: videoHandler,
		Consumer:     viewCountConsumer,
	}, nil
}