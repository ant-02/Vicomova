package grpc

import (
	"context"

	commerceAppCmd "vicomova/internal/commerce/application/command"
	commerceAppQuery "vicomova/internal/commerce/application/query"
	commerce "vicomova/third_party/kitex_gen/commerce"
)

type CommerceHandler struct {
	// Command Services
	walletCmdSvc      *commerceAppCmd.PointsWalletCommandService
	signInCmdSvc      *commerceAppCmd.SignInCommandService
	orderCmdSvc       *commerceAppCmd.OrderCommandService
	membershipCmdSvc  *commerceAppCmd.MembershipCommandService
	videoAccessCmdSvc *commerceAppCmd.VideoAccessCommandService
	flashSaleCmdSvc   *commerceAppCmd.FlashSaleCommandService
	lotteryCmdSvc     *commerceAppCmd.LotteryCommandService

	// Query Services
	walletQrySvc     *commerceAppQuery.PointsWalletQueryService
	signInQrySvc     *commerceAppQuery.SignInQueryService
	productQrySvc    *commerceAppQuery.ProductQueryService
	orderQrySvc      *commerceAppQuery.OrderQueryService
	membershipQrySvc *commerceAppQuery.MembershipQueryService
	flashSaleQrySvc  *commerceAppQuery.FlashSaleQueryService
	lotteryQrySvc    *commerceAppQuery.LotteryQueryService
}

func NewCommerceHandler(
	walletCmdSvc *commerceAppCmd.PointsWalletCommandService,
	walletQrySvc *commerceAppQuery.PointsWalletQueryService,
	signInCmdSvc *commerceAppCmd.SignInCommandService,
	signInQrySvc *commerceAppQuery.SignInQueryService,
	orderCmdSvc *commerceAppCmd.OrderCommandService,
	orderQrySvc *commerceAppQuery.OrderQueryService,
	membershipCmdSvc *commerceAppCmd.MembershipCommandService,
	membershipQrySvc *commerceAppQuery.MembershipQueryService,
	videoAccessCmdSvc *commerceAppCmd.VideoAccessCommandService,
	flashSaleCmdSvc *commerceAppCmd.FlashSaleCommandService,
	flashSaleQrySvc *commerceAppQuery.FlashSaleQueryService,
	lotteryCmdSvc *commerceAppCmd.LotteryCommandService,
	lotteryQrySvc *commerceAppQuery.LotteryQueryService,
) *CommerceHandler {
	return &CommerceHandler{
		walletCmdSvc:      walletCmdSvc,
		walletQrySvc:      walletQrySvc,
		signInCmdSvc:      signInCmdSvc,
		signInQrySvc:      signInQrySvc,
		orderCmdSvc:       orderCmdSvc,
		orderQrySvc:       orderQrySvc,
		membershipCmdSvc:  membershipCmdSvc,
		membershipQrySvc:  membershipQrySvc,
		videoAccessCmdSvc: videoAccessCmdSvc,
		flashSaleCmdSvc:   flashSaleCmdSvc,
		flashSaleQrySvc:   flashSaleQrySvc,
		lotteryCmdSvc:     lotteryCmdSvc,
		lotteryQrySvc:     lotteryQrySvc,
	}
}

func (h *CommerceHandler) GetWallet(ctx context.Context, req *commerce.GetWalletRequest) (*commerce.GetWalletResponse, error) {
	result, err := h.walletQrySvc.GetWallet(ctx, &commerceAppQuery.GetWalletQuery{UserID: req.UserId})
	if err != nil {
		return nil, err
	}
	if result.Wallet == nil {
		return &commerce.GetWalletResponse{Wallet: &commerce.PointsWallet{UserId: req.UserId, Balance: 0}}, nil
	}
	return &commerce.GetWalletResponse{
		Wallet: &commerce.PointsWallet{
			UserId:  result.Wallet.UserID,
			Balance: result.Wallet.Points,
			Version: result.Wallet.Version,
		},
	}, nil
}

func (h *CommerceHandler) RechargePoints(ctx context.Context, req *commerce.RechargePointsRequest) (*commerce.RechargePointsResponse, error) {
	result, err := h.walletCmdSvc.RechargePoints(ctx, &commerceAppCmd.RechargePointsCommand{
		UserID:    req.UserId,
		Amount:    req.Amount,
		PaymentID: req.PaymentId,
	})
	if err != nil {
		return nil, err
	}
	return &commerce.RechargePointsResponse{
		Wallet: &commerce.PointsWallet{
			UserId:  result.Wallet.UserID,
			Balance: result.NewPoints,
			Version: result.Wallet.Version,
		},
	}, nil
}

func (h *CommerceHandler) DailySignIn(ctx context.Context, req *commerce.DailySignInRequest) (*commerce.DailySignInResponse, error) {
	result, err := h.signInCmdSvc.ProcessSignIn(ctx, &commerceAppCmd.ProcessSignInCommand{
		UserID:          req.UserId,
		SignDate:        "", // 由 Kafka 事件触发，这里不需要
		ConsecutiveDays: 1,  // 默认1天，由事件携带真实连续天数
	})
	if err != nil {
		return nil, err
	}
	return &commerce.DailySignInResponse{
		Success:         result.Success,
		PointsEarned:    result.PointsEarned,
		ConsecutiveDays: int32(result.ConsecutiveDays),
		Message:         result.Message,
	}, nil
}

func (h *CommerceHandler) GetSignInInfo(ctx context.Context, req *commerce.GetSignInInfoRequest) (*commerce.GetSignInInfoResponse, error) {
	result, err := h.signInQrySvc.GetSignInInfo(ctx, &commerceAppQuery.GetSignInInfoQuery{UserID: req.UserId})
	if err != nil {
		return nil, err
	}
	info := &commerce.SignInInfo{
		UserId:        req.UserId,
		SignedInToday: false,
	}
	if result.Info != nil {
		info.ConsecutiveDays = int32(result.Info.ConsecutiveDays)
		info.LastSignInAt = result.Info.LastSignInAt
		info.TodayBonus = result.Info.TodayBonus
		info.SignedInToday = result.Info.SignedInToday
	}
	return &commerce.GetSignInInfoResponse{Info: info}, nil
}

func (h *CommerceHandler) ListProducts(ctx context.Context, req *commerce.ListProductsRequest) (*commerce.ListProductsResponse, error) {
	result, err := h.productQrySvc.ListProducts(ctx, &commerceAppQuery.ListProductsQuery{
		ProductType: int(req.ProductType),
		Cursor:      req.Cursor,
		Limit:       int(req.Limit),
	})
	if err != nil {
		return nil, err
	}
	products := make([]*commerce.Product, len(result.Products))
	for i, p := range result.Products {
		products[i] = &commerce.Product{
			Id:              p.ID,
			Name:            p.Name,
			Description:     p.Description,
			ProductType:     int32(p.ProductType),
			Price:           p.Price,
			Stock:           int64(p.Stock),
			VideoId:         p.VideoID,
			MembershipLevel: int32(p.MembershipLevel),
			IsActive:        p.IsActive,
		}
	}
	return &commerce.ListProductsResponse{
		Products:   products,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

func (h *CommerceHandler) ListOrders(ctx context.Context, req *commerce.ListOrdersRequest) (*commerce.ListOrdersResponse, error) {
	result, err := h.orderQrySvc.ListOrders(ctx, &commerceAppQuery.ListOrdersQuery{
		UserID: req.UserId,
		Cursor: req.Cursor,
		Limit:  int(req.Limit),
	})
	if err != nil {
		return nil, err
	}
	orders := make([]*commerce.Order, len(result.Orders))
	for i, o := range result.Orders {
		orders[i] = &commerce.Order{
			Id:        o.ID,
			UserId:    o.UserID,
			ProductId: o.ProductID,
			PricePaid: o.PricePaid,
			PayType:   int32(o.PayType),
			Status:    int32(o.Status),
		}
	}
	return &commerce.ListOrdersResponse{
		Orders:     orders,
		HasMore:    result.HasMore,
		NextCursor: result.NextCursor,
	}, nil
}

func (h *CommerceHandler) GetProduct(ctx context.Context, req *commerce.GetProductRequest) (*commerce.GetProductResponse, error) {
	// TODO: 实现
	return nil, nil
}

func (h *CommerceHandler) CreateOrder(ctx context.Context, req *commerce.CreateOrderRequest) (*commerce.CreateOrderResponse, error) {
	result, err := h.orderCmdSvc.CreateOrder(ctx, &commerceAppCmd.CreateOrderCommand{
		UserID:    req.UserId,
		ProductID: req.ProductId,
		PayType:   int(req.PayType),
	})
	if err != nil {
		return nil, err
	}
	return &commerce.CreateOrderResponse{
		Order: &commerce.Order{
			Id:        result.Order.ID,
			UserId:    result.Order.UserID,
			ProductId: result.Order.ProductID,
			PayType:   int32(result.Order.PayType),
			PricePaid: result.Order.PricePaid,
		},
	}, nil
}

func (h *CommerceHandler) PayOrder(ctx context.Context, req *commerce.PayOrderRequest) (*commerce.PayOrderResponse, error) {
	err := h.orderCmdSvc.PayOrder(ctx, &commerceAppCmd.PayOrderCommand{
		OrderID: req.OrderId,
		UserID:  req.UserId,
	})
	if err != nil {
		return nil, err
	}
	return &commerce.PayOrderResponse{Success: true}, nil
}

func (h *CommerceHandler) BuyMembership(ctx context.Context, req *commerce.BuyMembershipRequest) (*commerce.BuyMembershipResponse, error) {
	result, err := h.membershipCmdSvc.BuyMembership(ctx, &commerceAppCmd.BuyMembershipCommand{
		UserID:  req.UserId,
		Level:   int(req.Level),
		PayType: int(req.PayType),
	})
	if err != nil {
		return nil, err
	}
	return &commerce.BuyMembershipResponse{
		Membership: &commerce.Membership{
			Id:       result.Membership.ID,
			UserId:   result.Membership.UserID,
			Level:    int32(result.Membership.Level),
			StartAt:  result.Membership.StartAt,
			ExpireAt: result.Membership.ExpireAt,
			IsActive: result.Membership.IsActive,
		},
	}, nil
}

func (h *CommerceHandler) GetMembership(ctx context.Context, req *commerce.GetMembershipRequest) (*commerce.GetMembershipResponse, error) {
	result, err := h.membershipQrySvc.GetMembership(ctx, &commerceAppQuery.GetMembershipQuery{UserID: req.UserId})
	if err != nil {
		return nil, err
	}
	if result.Membership == nil {
		return &commerce.GetMembershipResponse{Membership: nil}, nil
	}
	return &commerce.GetMembershipResponse{
		Membership: &commerce.Membership{
			Id:       result.Membership.ID,
			UserId:   result.Membership.UserID,
			Level:    int32(result.Membership.Level),
			StartAt:  result.Membership.StartAt,
			ExpireAt: result.Membership.ExpireAt,
			IsActive: result.Membership.IsActive,
		},
	}, nil
}

func (h *CommerceHandler) BuyVideoAccess(ctx context.Context, req *commerce.BuyVideoAccessRequest) (*commerce.BuyVideoAccessResponse, error) {
	result, err := h.videoAccessCmdSvc.BuyVideoAccess(ctx, &commerceAppCmd.BuyVideoAccessCommand{
		UserID:  req.UserId,
		VideoID: req.VideoId,
		PayType: int(req.PayType),
	})
	if err != nil {
		return nil, err
	}
	return &commerce.BuyVideoAccessResponse{Success: result.Success}, nil
}

func (h *CommerceHandler) CheckVideoAccess(ctx context.Context, req *commerce.CheckVideoAccessRequest) (*commerce.CheckVideoAccessResponse, error) {
	// TODO: 实现视频权限检查
	return &commerce.CheckVideoAccessResponse{HasAccess: false, AccessLevel: 0}, nil
}

func (h *CommerceHandler) ListFlashSales(ctx context.Context, req *commerce.ListFlashSalesRequest) (*commerce.ListFlashSalesResponse, error) {
	result, err := h.flashSaleQrySvc.ListFlashSales(ctx, &commerceAppQuery.ListFlashSalesQuery{})
	if err != nil {
		return nil, err
	}
	sales := make([]*commerce.FlashSale, len(result.FlashSales))
	for i, s := range result.FlashSales {
		sales[i] = &commerce.FlashSale{
			Id:      s.ID,
			Name:    s.Name,
			StartAt: s.StartTime,
			EndAt:   s.EndTime,
		}
	}
	return &commerce.ListFlashSalesResponse{FlashSales: sales}, nil
}

func (h *CommerceHandler) FlashSalePurchase(ctx context.Context, req *commerce.FlashSalePurchaseRequest) (*commerce.FlashSalePurchaseResponse, error) {
	result, err := h.flashSaleCmdSvc.Purchase(ctx, &commerceAppCmd.FlashSalePurchaseCommand{
		UserID:      req.UserId,
		FlashSaleID: req.FlashSaleId,
		ProductID:   req.ProductId,
		PayType:     int(req.PayType),
	})
	if err != nil {
		return nil, err
	}
	return &commerce.FlashSalePurchaseResponse{
		Success: result.Success,
		OrderId: result.OrderID,
		Message: result.Message,
	}, nil
}

func (h *CommerceHandler) ListLotteries(ctx context.Context, req *commerce.ListLotteriesRequest) (*commerce.ListLotteriesResponse, error) {
	result, err := h.lotteryQrySvc.ListLotteries(ctx, &commerceAppQuery.ListLotteriesQuery{})
	if err != nil {
		return nil, err
	}
	lotteries := make([]*commerce.LotteryDraw, len(result.Lotteries))
	for i, l := range result.Lotteries {
		lotteries[i] = &commerce.LotteryDraw{
			Id:          l.ID,
			Name:        l.Name,
			Description: l.Description,
			EntryPoints: l.EntryPoints,
			StartAt:     l.StartTime,
			EndAt:       l.EndTime,
		}
	}
	return &commerce.ListLotteriesResponse{Lotteries: lotteries}, nil
}

func (h *CommerceHandler) ParticipateLottery(ctx context.Context, req *commerce.ParticipateLotteryRequest) (*commerce.ParticipateLotteryResponse, error) {
	result, err := h.lotteryCmdSvc.Participate(ctx, &commerceAppCmd.ParticipateLotteryCommand{
		UserID:    req.UserId,
		LotteryID: req.LotteryId,
	})
	if err != nil {
		return nil, err
	}
	return &commerce.ParticipateLotteryResponse{
		Success: result.Success,
		Won:     result.Won,
		Prize:   result.Prize,
		Message: result.Message,
	}, nil
}
