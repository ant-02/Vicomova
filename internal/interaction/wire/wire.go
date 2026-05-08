package wire

import (
	"vicomova/internal/interaction/application/command"
	"vicomova/internal/interaction/application/query"
	infraMysql "vicomova/internal/interaction/infrastructure/persistence/mysql"
	"vicomova/internal/interaction/interfaces/grpc"
	"vicomova/pkg/config"
	"vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/infrastructure/redis"
	"vicomova/pkg/log"
)

type Provider struct {
	MySQL              *mysql.Client
	Redis              *redis.Client
	InteractionHandler *grpc.InteractionHandler
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

	if err := mysql.GetDB().AutoMigrate(&infraMysql.LikePO{}, &infraMysql.CommentPO{}, &infraMysql.FavoritePO{}); err != nil {
		return nil, err
	}

	likeRepo := infraMysql.NewLikeRepository(mysqlClient)
	commentRepo := infraMysql.NewCommentRepository(mysqlClient)
	favoriteRepo := infraMysql.NewFavoriteRepository(mysqlClient)

	cmdSvc := command.NewInteractionCommandService(likeRepo, commentRepo, favoriteRepo)
	querySvc := query.NewInteractionQueryService(likeRepo, commentRepo, favoriteRepo)

	handler := grpc.NewInteractionHandler(cmdSvc, querySvc)

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
		MySQL:              mysqlClient,
		Redis:              redis.GetClient(),
		InteractionHandler: handler,
	}, nil
}
