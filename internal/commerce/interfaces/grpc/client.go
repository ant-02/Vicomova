package grpc

import (
	"context"

	commerce "vicomova/third_party/kitex_gen/commerce"
	commerceservice "vicomova/third_party/kitex_gen/commerce/commerceservice"

	"github.com/cloudwego/kitex/client"
)

type CommerceClient struct {
	cli commerceservice.Client
}

// NewCommerceClient creates RPC client, directly connects to specified address
func NewCommerceClient(serviceName, addr string) (*CommerceClient, error) {
	cli, err := commerceservice.NewClient(serviceName,
		client.WithHostPorts(addr),
	)
	if err != nil {
		return nil, err
	}
	return &CommerceClient{cli: cli}, nil
}

func (c *CommerceClient) GetWallet(ctx context.Context, userID int64) (*commerce.GetWalletResponse, error) {
	return c.cli.GetWallet(ctx, &commerce.GetWalletRequest{UserId: userID})
}

func (c *CommerceClient) RechargePoints(ctx context.Context, userID int64, amount int64, paymentID string) (*commerce.RechargePointsResponse, error) {
	return c.cli.RechargePoints(ctx, &commerce.RechargePointsRequest{
		UserId:    userID,
		Amount:    amount,
		PaymentId: paymentID,
	})
}

func (c *CommerceClient) DailySignIn(ctx context.Context, userID int64) (*commerce.DailySignInResponse, error) {
	return c.cli.DailySignIn(ctx, &commerce.DailySignInRequest{UserId: userID})
}

func (c *CommerceClient) GetSignInInfo(ctx context.Context, userID int64) (*commerce.GetSignInInfoResponse, error) {
	return c.cli.GetSignInInfo(ctx, &commerce.GetSignInInfoRequest{UserId: userID})
}

func (c *CommerceClient) ListProducts(ctx context.Context, productType int32, limit int32, cursor string) (*commerce.ListProductsResponse, error) {
	return c.cli.ListProducts(ctx, &commerce.ListProductsRequest{
		ProductType: productType,
		Limit:       limit,
		Cursor:      cursor,
	})
}

func (c *CommerceClient) GetProduct(ctx context.Context, productID int64) (*commerce.GetProductResponse, error) {
	return c.cli.GetProduct(ctx, &commerce.GetProductRequest{ProductId: productID})
}

func (c *CommerceClient) CreateOrder(ctx context.Context, userID int64, productID int64, payType int32) (*commerce.CreateOrderResponse, error) {
	return c.cli.CreateOrder(ctx, &commerce.CreateOrderRequest{
		UserId:    userID,
		ProductId: productID,
		PayType:   payType,
	})
}

func (c *CommerceClient) PayOrder(ctx context.Context, orderID int64, userID int64) (*commerce.PayOrderResponse, error) {
	return c.cli.PayOrder(ctx, &commerce.PayOrderRequest{
		OrderId: orderID,
		UserId:  userID,
	})
}

func (c *CommerceClient) ListOrders(ctx context.Context, userID int64, limit int32, cursor string) (*commerce.ListOrdersResponse, error) {
	return c.cli.ListOrders(ctx, &commerce.ListOrdersRequest{
		UserId: userID,
		Limit:  limit,
		Cursor: cursor,
	})
}

func (c *CommerceClient) BuyMembership(ctx context.Context, userID int64, level int32, payType int32) (*commerce.BuyMembershipResponse, error) {
	return c.cli.BuyMembership(ctx, &commerce.BuyMembershipRequest{
		UserId:  userID,
		Level:   level,
		PayType: payType,
	})
}

func (c *CommerceClient) GetMembership(ctx context.Context, userID int64) (*commerce.GetMembershipResponse, error) {
	return c.cli.GetMembership(ctx, &commerce.GetMembershipRequest{UserId: userID})
}

func (c *CommerceClient) CheckVideoAccess(ctx context.Context, userID int64, videoID int64) (*commerce.CheckVideoAccessResponse, error) {
	return c.cli.CheckVideoAccess(ctx, &commerce.CheckVideoAccessRequest{
		UserId:  userID,
		VideoId: videoID,
	})
}

func (c *CommerceClient) ListFlashSales(ctx context.Context) (*commerce.ListFlashSalesResponse, error) {
	return c.cli.ListFlashSales(ctx, &commerce.ListFlashSalesRequest{})
}

func (c *CommerceClient) FlashSalePurchase(ctx context.Context, userID int64, flashSaleID int64, productID int64, payType int32) (*commerce.FlashSalePurchaseResponse, error) {
	return c.cli.FlashSalePurchase(ctx, &commerce.FlashSalePurchaseRequest{
		UserId:      userID,
		FlashSaleId: flashSaleID,
		ProductId:   productID,
		PayType:     payType,
	})
}

func (c *CommerceClient) ListLotteries(ctx context.Context) (*commerce.ListLotteriesResponse, error) {
	return c.cli.ListLotteries(ctx, &commerce.ListLotteriesRequest{})
}

func (c *CommerceClient) ParticipateLottery(ctx context.Context, userID int64, lotteryID int64) (*commerce.ParticipateLotteryResponse, error) {
	return c.cli.ParticipateLottery(ctx, &commerce.ParticipateLotteryRequest{
		UserId:    userID,
		LotteryId: lotteryID,
	})
}
