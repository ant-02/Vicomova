package wire

import (
	commerceAppCmd "vicomova/internal/commerce/application/command"
	commerceAppQuery "vicomova/internal/commerce/application/query"
	commerceSvc "vicomova/internal/commerce/domain/service"
	infraMysql "vicomova/internal/commerce/infrastructure/persistence/mysql"
	redisCache "vicomova/internal/commerce/infrastructure/persistence/redis"
	commerceGrpc "vicomova/internal/commerce/interfaces/grpc"
	"vicomova/pkg/config"
	"vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/infrastructure/redis"
	"vicomova/pkg/log"
)

type Provider struct {
	MySQL           *mysql.Client
	Redis           *redis.Client
	CommerceHandler *commerceGrpc.CommerceHandler
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

	// AutoMigrate Commerce tables
	if err := mysql.GetDB().AutoMigrate(
		&infraMysql.PointsWalletPO{},
		&infraMysql.SignInRecordPO{},
		&infraMysql.TransactionPO{},
		&infraMysql.ProductPO{},
		&infraMysql.OrderPO{},
		&infraMysql.MembershipPO{},
		&infraMysql.VideoPermissionPO{},
		&infraMysql.FlashSalePO{},
		&infraMysql.FlashSaleStockPO{},
		&infraMysql.LotteryDrawPO{},
		&infraMysql.LotteryRecordPO{},
	); err != nil {
		return nil, err
	}

	// Repositories
	walletRepo := infraMysql.NewPointsWalletRepository(mysqlClient)
	signInRepo := infraMysql.NewSignInRecordRepository(mysqlClient)
	txnRepo := infraMysql.NewTransactionRepository(mysqlClient)
	productRepo := infraMysql.NewProductRepository(mysqlClient)
	orderRepo := infraMysql.NewOrderRepository(mysqlClient)
	membershipRepo := infraMysql.NewMembershipRepository(mysqlClient)
	videoPermRepo := infraMysql.NewVideoPermissionRepository(mysqlClient)
	flashSaleRepo := infraMysql.NewFlashSaleRepository(mysqlClient)
	stockRepo := infraMysql.NewFlashSaleStockRepository(mysqlClient)
	lotteryRepo := infraMysql.NewLotteryRepository(mysqlClient)
	lotteryRecordRepo := infraMysql.NewLotteryRecordRepository(mysqlClient)

	// Redis Cache
	walletCache := redisCache.NewWalletCache(redisClient)
	flashSaleCache := redisCache.NewFlashSaleCache(redisClient)
	lotteryCache := redisCache.NewLotteryCache(redisClient)
	_ = walletCache // suppress unused warning (used in command services)

	// Domain Services
	bonusCalc := commerceSvc.NewSignInBonusCalculator()

	// Command Services
	walletCmdSvc := commerceAppCmd.NewPointsWalletCommandService(walletRepo, txnRepo)
	signInCmdSvc := commerceAppCmd.NewSignInCommandService(signInRepo, walletRepo, txnRepo, bonusCalc)
	orderCmdSvc := commerceAppCmd.NewOrderCommandService(orderRepo, productRepo, walletRepo, txnRepo, membershipRepo, videoPermRepo)
	membershipCmdSvc := commerceAppCmd.NewMembershipCommandService(membershipRepo)
	videoAccessCmdSvc := commerceAppCmd.NewVideoAccessCommandService(videoPermRepo)
	flashSaleCmdSvc := commerceAppCmd.NewFlashSaleCommandService(flashSaleRepo, stockRepo, productRepo, orderRepo, flashSaleCache)
	lotteryCmdSvc := commerceAppCmd.NewLotteryCommandService(lotteryRepo, lotteryRecordRepo, lotteryCache)

	// Query Services
	walletQrySvc := commerceAppQuery.NewPointsWalletQueryService(walletRepo)
	signInQrySvc := commerceAppQuery.NewSignInQueryService(signInRepo)
	_ = commerceAppQuery.NewProductQueryService(productRepo) // productQrySvc (used in handler)
	orderQrySvc := commerceAppQuery.NewOrderQueryService(orderRepo)
	membershipQrySvc := commerceAppQuery.NewMembershipQueryService(membershipRepo)
	flashSaleQrySvc := commerceAppQuery.NewFlashSaleQueryService(flashSaleRepo)
	lotteryQrySvc := commerceAppQuery.NewLotteryQueryService(lotteryRepo)

	// gRPC Handler
	commerceHandler := commerceGrpc.NewCommerceHandler(
		walletCmdSvc, walletQrySvc,
		signInCmdSvc, signInQrySvc,
		orderCmdSvc, orderQrySvc,
		membershipCmdSvc, membershipQrySvc,
		videoAccessCmdSvc,
		flashSaleCmdSvc, flashSaleQrySvc,
		lotteryCmdSvc, lotteryQrySvc,
	)

	config.RegisterCallback(func(newCfg *config.Config) {
		if err := mysql.Reload(&newCfg.Database); err != nil {
			log.Error.Printf("failed to reload mysql: %v", err)
		}
		if err := redis.Reload(&newCfg.Redis); err != nil {
			log.Error.Printf("failed to reload redis: %v", err)
		}
	})

	return &Provider{
		MySQL:           mysqlClient,
		Redis:           redis.GetClient(),
		CommerceHandler: commerceHandler,
	}, nil
}
