package wire

import (
	sharedMysql "vicomova/internal/shared/infrastructure/data/mysql"
	sharedRedis "vicomova/internal/shared/infrastructure/data/redis"
	videoAppCmd "vicomova/internal/video/application/command"
	videoAppQuery "vicomova/internal/video/application/query"
	infraMysql "vicomova/internal/video/infrastructure/persistence/mysql"
	infraService "vicomova/internal/video/infrastructure/service"
	infraStorage "vicomova/internal/video/infrastructure/storage"
	videoGrpc "vicomova/internal/video/interfaces/grpc"
	"vicomova/pkg/config"
)

type Provider struct {
	MySQL        *sharedMysql.Client
	Redis        *sharedRedis.Client
	VideoHandler *videoGrpc.VideoHandler
}

func NewProvider(cfg *config.Config) (*Provider, error) {
	// 初始化数据库
	if err := sharedMysql.Init(&cfg.Database); err != nil {
		return nil, err
	}

	// 初始化 Redis
	if err := sharedRedis.Init(&cfg.Redis); err != nil {
		return nil, err
	}

	// 获取 MySQL/Redis 客户端
	mysqlClient := sharedMysql.GetClient()

	// 自动迁移
	if err := sharedMysql.GetDB().AutoMigrate(&infraMysql.VideoPO{}); err != nil {
		return nil, err
	}

	// 初始化 Repository
	videoRepo := infraMysql.NewVideoRepository(mysqlClient)

	// 初始化 Storage
	var storage infraStorage.VideoStorage
	switch cfg.Video.Storage.Type {
	case "local":
		storage = infraStorage.NewLocalDiskStorage(&infraStorage.LocalDiskConfig{
			BasePath: cfg.Video.Storage.Local.BasePath,
			BaseURL:  cfg.Video.Storage.Local.BaseURL,
		})
	default:
		storage = infraStorage.NewLocalDiskStorage(&infraStorage.LocalDiskConfig{
			BasePath: cfg.Video.Storage.Local.BasePath,
			BaseURL:  cfg.Video.Storage.Local.BaseURL,
		})
	}

	// 初始化 Domain Service
	hotAlgo := infraService.NewWilsonHotAlgorithm(videoRepo)

	// 初始化 Application Service
	cmdSvc := videoAppCmd.NewVideoCommandService(videoRepo, storage)
	querySvc := videoAppQuery.NewVideoQueryService(videoRepo, hotAlgo)

	// 初始化 Handler
	videoHandler := videoGrpc.NewVideoHandler(cmdSvc, querySvc)

	return &Provider{
		MySQL:        mysqlClient,
		Redis:        sharedRedis.GetClient(),
		VideoHandler: videoHandler,
	}, nil
}
