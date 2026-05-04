package wire

import (
	interactionAppCmd "vicomova/internal/interaction/application/command"
	interactionAppQuery "vicomova/internal/interaction/application/query"
	infraMysql "vicomova/internal/interaction/infrastructure/persistence/mysql"
	interactionGrpc "vicomova/internal/interaction/interfaces/grpc"
	sharedMysql "vicomova/internal/shared/infrastructure/data/mysql"
	sharedRedis "vicomova/internal/shared/infrastructure/data/redis"
	"vicomova/pkg/config"
)

type Provider struct {
	MySQL             *sharedMysql.Client
	Redis             *sharedRedis.Client
	InteractionHandler *interactionGrpc.InteractionHandler
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

	// 初始化 Repositories
	likeRepo := infraMysql.NewLikeRepository(mysqlClient)
	commentRepo := infraMysql.NewCommentRepository(mysqlClient)
	favoriteRepo := infraMysql.NewFavoriteRepository(mysqlClient)

	// 初始化 Application Service
	cmdSvc := interactionAppCmd.NewInteractionCommandService(likeRepo, commentRepo, favoriteRepo)
	querySvc := interactionAppQuery.NewInteractionQueryService(likeRepo, commentRepo, favoriteRepo)

	// 初始化 Handler
	handler := interactionGrpc.NewInteractionHandler(cmdSvc, querySvc)

	return &Provider{
		MySQL:             mysqlClient,
		Redis:             sharedRedis.GetClient(),
		InteractionHandler: handler,
	}, nil
}