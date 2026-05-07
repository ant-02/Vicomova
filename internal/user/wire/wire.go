package wire

import (
	sharedMysql "vicomova/internal/shared/infrastructure/data/mysql"
	sharedRedis "vicomova/internal/shared/infrastructure/data/redis"
	sharedHasher "vicomova/internal/shared/pkg/hasher"
	appCommand "vicomova/internal/user/application/command"
	appQuery "vicomova/internal/user/application/query"
	svc "vicomova/internal/user/domain/service"
	infraEmail "vicomova/internal/user/infrastructure/external/email"
	infraMysql "vicomova/internal/user/infrastructure/persistence/mysql"
	infraRedis "vicomova/internal/user/infrastructure/persistence/redis"
	rpc "vicomova/internal/user/interfaces/grpc"
	"vicomova/pkg/config"
)

type Provider struct {
	MySQL       *sharedMysql.Client
	Redis       *sharedRedis.Client
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

	// 获取 MySQL/Redis 客户端
	mysqlClient := sharedMysql.GetClient()
	redisClient := sharedRedis.GetClient()

	// 自动迁移
	if err := sharedMysql.GetDB().AutoMigrate(&infraMysql.UserPO{}); err != nil {
		return nil, err
	}

	// 初始化 Repository
	userRepo := infraMysql.NewUserRepository(mysqlClient)
	refreshTokenRepo := infraRedis.NewRefreshTokenRepository(redisClient)
	emailCodeRepo := infraRedis.NewEmailCodeRepository(redisClient)

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

	// 初始化 Password Hasher
	hasher := &sharedHasher.SHA256Hasher{}

	// 初始化 Application Service
	cmdSvc := appCommand.NewUserCommandService(userRepo, refreshTokenRepo, emailCodeRepo, emailService, tokenSvc, hasher)
	querySvc := appQuery.NewUserQueryService(userRepo)

	// 初始化 Handler
	userHandler := rpc.NewUserHandler(cmdSvc, querySvc)

	return &Provider{
		MySQL:       mysqlClient,
		Redis:       redisClient,
		UserHandler: userHandler,
	}, nil
}
