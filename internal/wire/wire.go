package wire

import (
	"vicomova/internal/application/user"
	infraMysql "vicomova/internal/infrastructure/persistence/mysql"
	infraRedis "vicomova/internal/infrastructure/persistence/redis"
	domainUser "vicomova/internal/domain/user"
	"vicomova/internal/interface/rpc"
	"vicomova/internal/data/mysql"
	redisClient "vicomova/internal/data/redis"
	"vicomova/pkg/config"
)

type Provider struct {
	UserHandler *rpc.UserHandler
}

func NewProvider(cfg *config.Config) (*Provider, error) {
	// 初始化数据库
	if err := mysql.Init(&cfg.Database); err != nil {
		return nil, err
	}

	// 初始化 Redis
	if err := redisClient.Init(&cfg.Redis); err != nil {
		return nil, err
	}

	// 初始化 Repository
	userRepo := infraMysql.NewUserRepository()
	refreshTokenRepo := infraRedis.NewRefreshTokenRepository()

	// 初始化 Domain Service
	tokenSvc := domainUser.NewTokenService(cfg.JWT.Secret)

	// 初始化 Application Service
	cmdSvc := user.NewUserCommandService(userRepo, refreshTokenRepo, tokenSvc)
	querySvc := user.NewUserQueryService(userRepo)

	// 初始化 Handler
	userHandler := rpc.NewUserHandler(cmdSvc, querySvc)

	return &Provider{
		UserHandler: userHandler,
	}, nil
}
