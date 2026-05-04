package wire

import (
	appCommand "vicomova/internal/user/application/command"
	appQuery "vicomova/internal/user/application/query"
	infraEmail "vicomova/internal/user/infrastructure/external/email"
	infraMysql "vicomova/internal/user/infrastructure/persistence/mysql"
	infraRedis "vicomova/internal/user/infrastructure/persistence/redis"
	svc "vicomova/internal/user/domain/service"
	rpc "vicomova/internal/user/interfaces/grpc"
	sharedMysql "vicomova/internal/shared/infrastructure/data/mysql"
	sharedRedis "vicomova/internal/shared/infrastructure/data/redis"
	"vicomova/pkg/config"
)

type Provider struct {
	MySQL      *sharedMysql.Client
	Redis      *sharedRedis.Client
	UserHandler *rpc.UserHandler
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

	// 初始化 Repository
	userRepo := infraMysql.NewUserRepository()
	refreshTokenRepo := infraRedis.NewRefreshTokenRepository()
	emailCodeRepo := infraRedis.NewEmailCodeRepository()

	// 初始化 Email Service
	var emailService infraEmail.EmailService
	if cfg.Email.AccessKey != "" && cfg.Email.AccountName != "" && cfg.Email.Region != "" {
		var err error
		emailService, err = infraEmail.NewAliyunEmailService(&cfg.Email)
		if err != nil {
			return nil, err
		}
	}

	// 初始化 Domain Service
	tokenSvc := svc.NewTokenService(cfg.JWT.Secret)

	// 初始化 Application Service
	cmdSvc := appCommand.NewUserCommandService(userRepo, refreshTokenRepo, emailCodeRepo, emailService, tokenSvc)
	querySvc := appQuery.NewUserQueryService(userRepo)

	// 初始化 Handler
	userHandler := rpc.NewUserHandler(cmdSvc, querySvc)

	return &Provider{
		MySQL:      sharedMysql.GetClient(),
		Redis:      sharedRedis.GetClient(),
		UserHandler: userHandler,
	}, nil
}