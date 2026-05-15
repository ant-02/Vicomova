package wire

import (
	"vicomova/internal/chat/application/command"
	"vicomova/internal/chat/application/query"
	infraMysql "vicomova/internal/chat/infrastructure/persistence/mysql"
	"vicomova/internal/chat/interfaces/grpc"
	userRpc "vicomova/internal/user/interfaces/grpc"
	"vicomova/pkg/config"
	"vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/infrastructure/redis"
	"vicomova/pkg/log"
)

type Provider struct {
	MySQL       *mysql.Client
	Redis       *redis.Client
	ChatHandler *grpc.ChatHandler
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

	if err := mysql.GetDB().AutoMigrate(
		&infraMysql.FollowPO{},
		&infraMysql.ConversationPO{},
		&infraMysql.GroupMemberPO{},
		&infraMysql.MessagePO{},
		&infraMysql.OfflineMessagePO{},
	); err != nil {
		return nil, err
	}

	followRepo := infraMysql.NewFollowRepository(mysqlClient)
	conversationRepo := infraMysql.NewConversationRepository(mysqlClient)
	groupMemberRepo := infraMysql.NewGroupMemberRepository(mysqlClient)
	messageRepo := infraMysql.NewMessageRepository(mysqlClient)
	offlineMessageRepo := infraMysql.NewOfflineMessageRepository(mysqlClient)

	cmdSvc := command.NewChatCommandService(followRepo, conversationRepo, groupMemberRepo, messageRepo, offlineMessageRepo)
	querySvc := query.NewChatQueryService(followRepo, groupMemberRepo, messageRepo, offlineMessageRepo)

	userCli, err := userRpc.NewUserClient("user", cfg.Services["user"].Addr)
	if err != nil {
		log.Error.Printf("failed to create user client: %v", err)
		return nil, err
	}

	handler := grpc.NewChatHandler(cmdSvc, querySvc, userCli)

	config.RegisterCallback(func(newCfg *config.Config) {
		if err := mysql.Reload(&newCfg.Database); err != nil {
			log.Error.Printf("failed to reload mysql: %v", err)
		}
		if err := redis.Reload(&newCfg.Redis); err != nil {
			log.Error.Printf("failed to reload redis: %v", err)
		}
	})

	return &Provider{
		MySQL:       mysqlClient,
		Redis:       redis.GetClient(),
		ChatHandler: handler,
	}, nil
}
