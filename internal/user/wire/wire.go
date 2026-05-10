package wire

import (
	"vicomova/internal/user/application/command"
	"vicomova/internal/user/application/query"
	"vicomova/internal/user/domain/service"
	email "vicomova/pkg/infrastructure/email"
	infraMysql "vicomova/internal/user/infrastructure/persistence/mysql"
	infraRedis "vicomova/internal/user/infrastructure/persistence/redis"
	usergrpc "vicomova/internal/user/interfaces/grpc"
	"vicomova/pkg/config"
	"vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/infrastructure/redis"
	"vicomova/pkg/log"
	"vicomova/pkg/utils"
)

type Provider struct {
	MySQL       *mysql.Client
	Redis       *redis.Client
	UserHandler *usergrpc.UserHandler
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

	if err := mysql.GetDB().AutoMigrate(&infraMysql.UserPO{}); err != nil {
		return nil, err
	}

	userRepo := infraMysql.NewUserRepository(mysqlClient)
	refreshTokenRepo := infraRedis.NewRefreshTokenRepository(redisClient)
	emailCodeRepo := infraRedis.NewEmailCodeRepository(redisClient)

	var emailService email.EmailService
	if cfg.Email.AccountName != "" && cfg.Email.Region != "" {
		var err error
		emailService, err = email.NewAliyunEmailService(&cfg.Email)
		if err != nil {
			return nil, err
		}
	}

	tokenSvc := service.NewTokenService(cfg.JWT.Secret)

	hasher := &utils.SHA256Hasher{}

	cmdSvc := command.NewUserCommandService(userRepo, refreshTokenRepo, emailCodeRepo, emailService, tokenSvc, hasher)
	querySvc := query.NewUserQueryService(userRepo)

	userHandler := usergrpc.NewUserHandler(cmdSvc, querySvc)

	// Register config change callbacks for hot-reload
	config.RegisterCallback(func(newCfg *config.Config) {
		if err := mysql.Reload(&newCfg.Database); err != nil {
			log.Error.Printf("failed to reload mysql: %v", err)
		}
		if err := redis.Reload(&newCfg.Redis); err != nil {
			log.Error.Printf("failed to reload redis: %v", err)
		}
		if newCfg.JWT.Secret != cfg.JWT.Secret {
			log.Warn.Printf("JWT secret changed, please restart service to take effect")
		}
		if newCfg.Service.Addr != cfg.Service.Addr {
			log.Warn.Printf("service addr changed to %s, please restart service to take effect", newCfg.Service.Addr)
		}
	})

	return &Provider{
		MySQL:       mysqlClient,
		Redis:       redisClient,
		UserHandler: userHandler,
	}, nil
}
