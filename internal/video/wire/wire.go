package wire

import (
	"vicomova/internal/video/application/command"
	"vicomova/internal/video/application/query"
	infraMysql "vicomova/internal/video/infrastructure/persistence/mysql"
	"vicomova/internal/video/infrastructure/service"
	"vicomova/internal/video/interfaces/grpc"
	"vicomova/pkg/config"
	"vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/infrastructure/oss"
	"vicomova/pkg/infrastructure/redis"
	"vicomova/pkg/log"
)

type Provider struct {
	MySQL        *mysql.Client
	Redis        *redis.Client
	VideoHandler *grpc.VideoHandler
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

	if err := mysql.GetDB().AutoMigrate(&infraMysql.VideoPO{}); err != nil {
		return nil, err
	}

	videoRepo := infraMysql.NewVideoRepository(mysqlClient)

	// 初始化 OSS（可选，未配置时跳过）
	var ossClient oss.OSS
	if cfg.OSS.Qiniu.AccessKey != "" {
		ossCfg := &oss.OSSConfig{
			AccessKey: cfg.OSS.Qiniu.AccessKey,
			SecretKey: cfg.OSS.Qiniu.SecretKey,
			Bucket:    cfg.OSS.Qiniu.Bucket,
			Domain:    cfg.OSS.Qiniu.Domain,
			UploadHost: cfg.OSS.Qiniu.UploadHost,
		}
		var err error
		ossClient, err = oss.NewOSS(oss.OSSTypeQiniu, ossCfg)
		if err != nil {
			return nil, err
		}
	}

	hotAlgo := service.NewWilsonHotAlgorithm(videoRepo)

	cmdSvc := command.NewVideoCommandService(videoRepo, ossClient)
	querySvc := query.NewVideoQueryService(videoRepo, hotAlgo, ossClient)

	videoHandler := grpc.NewVideoHandler(cmdSvc, querySvc)

	config.RegisterCallback(func(newCfg *config.Config) {
		if err := mysql.Reload(&newCfg.Database); err != nil {
			log.Error.Printf("failed to reload mysql: %v", err)
		}
		if err := redis.Reload(&newCfg.Redis); err != nil {
			log.Error.Printf("failed to reload redis: %v", err)
		}
		if newCfg.Service.Addr != cfg.Service.Addr {
			log.Warn.Printf("service addr changed to %s, please restart service to take effect", newCfg.Service.Addr)
		}
	})

	return &Provider{
		MySQL:        mysqlClient,
		Redis:        redis.GetClient(),
		VideoHandler: videoHandler,
	}, nil
}
