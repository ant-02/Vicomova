package handler

import (
	"context"
	"strconv"

	rpc "vicomova/internal/commerce/interfaces/grpc"
	"vicomova/pkg/constants"
	hertz "vicomova/pkg/infrastructure/hertz"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// CommerceHandler handles HTTP requests for commerce service
type CommerceHandler struct {
	commerceClient *rpc.CommerceClient
}

func NewCommerceHandler(commerceClient *rpc.CommerceClient) *CommerceHandler {
	return &CommerceHandler{commerceClient: commerceClient}
}

// GetWallet 获取用户积分钱包
func (h *CommerceHandler) GetWallet(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	resp, err := h.commerceClient.GetWallet(ctx, userID)
	if err != nil {
		hlog.Errorf("GetWallet: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get wallet"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"user_id": resp.Wallet.UserId,
		"balance": resp.Wallet.Balance,
		"version": resp.Wallet.Version,
	}))
}

// RechargePoints 充值积分
func (h *CommerceHandler) RechargePoints(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req struct {
		Amount    int64  `json:"amount"`
		PaymentID string `json:"payment_id"`
	}
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request"))
		return
	}

	resp, err := h.commerceClient.RechargePoints(ctx, userID, req.Amount, req.PaymentID)
	if err != nil {
		hlog.Errorf("RechargePoints: userID=%d amount=%d failed: %v", userID, req.Amount, err)
		c.JSON(500, hertz.Fail(500, "Failed to recharge"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"user_id": resp.Wallet.UserId,
		"balance": resp.Wallet.Balance,
	}))
}

// DailySignIn 每日签到
func (h *CommerceHandler) DailySignIn(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	resp, err := h.commerceClient.DailySignIn(ctx, userID)
	if err != nil {
		hlog.Errorf("DailySignIn: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to sign in"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"success":          resp.Success,
		"points_earned":    resp.PointsEarned,
		"consecutive_days": resp.ConsecutiveDays,
		"message":          resp.Message,
	}))
}

// GetSignInInfo 获取签到信息
func (h *CommerceHandler) GetSignInInfo(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	resp, err := h.commerceClient.GetSignInInfo(ctx, userID)
	if err != nil {
		hlog.Errorf("GetSignInInfo: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get sign in info"))
		return
	}

	if resp.Info == nil {
		c.JSON(200, hertz.Success(map[string]interface{}{
			"signed_in_today": false,
		}))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"user_id":          resp.Info.UserId,
		"consecutive_days": resp.Info.ConsecutiveDays,
		"last_sign_in_at":  resp.Info.LastSignInAt,
		"today_bonus":      resp.Info.TodayBonus,
		"signed_in_today":  resp.Info.SignedInToday,
	}))
}

// ListProducts 获取商品列表
func (h *CommerceHandler) ListProducts(ctx context.Context, c *app.RequestContext) {
	productType, _ := strconv.ParseInt(c.Query("type"), 10, 32)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)
	cursor := c.Query("cursor")

	resp, err := h.commerceClient.ListProducts(ctx, int32(productType), int32(limit), cursor)
	if err != nil {
		hlog.Errorf("ListProducts: failed: %v", err)
		c.JSON(500, hertz.Fail(500, "Failed to list products"))
		return
	}

	products := make([]map[string]interface{}, 0, len(resp.Products))
	for _, p := range resp.Products {
		products = append(products, map[string]interface{}{
			"id":               p.Id,
			"name":             p.Name,
			"description":      p.Description,
			"product_type":     p.ProductType,
			"price":            p.Price,
			"stock":            p.Stock,
			"video_id":         p.VideoId,
			"membership_level": p.MembershipLevel,
			"is_active":        p.IsActive,
		})
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"products":    products,
		"next_cursor": resp.NextCursor,
		"has_more":    resp.HasMore,
	}))
}

// GetProduct 获取单个商品
func (h *CommerceHandler) GetProduct(ctx context.Context, c *app.RequestContext) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid product_id"))
		return
	}

	resp, err := h.commerceClient.GetProduct(ctx, productID)
	if err != nil {
		hlog.Errorf("GetProduct: productID=%d failed: %v", productID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get product"))
		return
	}

	if resp.Product == nil {
		c.JSON(404, hertz.Fail(404, "Product not found"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"id":               resp.Product.Id,
		"name":             resp.Product.Name,
		"description":      resp.Product.Description,
		"product_type":     resp.Product.ProductType,
		"price":            resp.Product.Price,
		"stock":            resp.Product.Stock,
		"video_id":         resp.Product.VideoId,
		"membership_level": resp.Product.MembershipLevel,
		"is_active":        resp.Product.IsActive,
	}))
}

// CreateOrder 创建订单
func (h *CommerceHandler) CreateOrder(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req struct {
		ProductID int64 `json:"product_id"`
		PayType   int   `json:"pay_type"`
	}
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request"))
		return
	}

	resp, err := h.commerceClient.CreateOrder(ctx, userID, req.ProductID, int32(req.PayType))
	if err != nil {
		hlog.Errorf("CreateOrder: userID=%d productID=%d failed: %v", userID, req.ProductID, err)
		c.JSON(500, hertz.Fail(500, "Failed to create order"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"order_id": resp.Order.Id,
	}))
}

// PayOrder 支付订单
func (h *CommerceHandler) PayOrder(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid order_id"))
		return
	}

	resp, err := h.commerceClient.PayOrder(ctx, orderID, userID)
	if err != nil {
		hlog.Errorf("PayOrder: orderID=%d userID=%d failed: %v", orderID, userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to pay order"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"success": resp.Success,
	}))
}

// ListOrders 获取订单列表
func (h *CommerceHandler) ListOrders(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 32)
	cursor := c.Query("cursor")

	resp, err := h.commerceClient.ListOrders(ctx, userID, int32(limit), cursor)
	if err != nil {
		hlog.Errorf("ListOrders: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to list orders"))
		return
	}

	orders := make([]map[string]interface{}, 0, len(resp.Orders))
	for _, o := range resp.Orders {
		orders = append(orders, map[string]interface{}{
			"id":         o.Id,
			"user_id":    o.UserId,
			"product_id": o.ProductId,
			"price_paid": o.PricePaid,
			"pay_type":   o.PayType,
			"status":     o.Status,
		})
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"orders":      orders,
		"has_more":    resp.HasMore,
		"next_cursor": resp.NextCursor,
	}))
}

// GetMembership 获取会员信息
func (h *CommerceHandler) GetMembership(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	resp, err := h.commerceClient.GetMembership(ctx, userID)
	if err != nil {
		hlog.Errorf("GetMembership: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get membership"))
		return
	}

	if resp.Membership == nil {
		c.JSON(200, hertz.Success(map[string]interface{}{
			"is_active": false,
		}))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"id":        resp.Membership.Id,
		"user_id":   resp.Membership.UserId,
		"level":     resp.Membership.Level,
		"start_at":  resp.Membership.StartAt,
		"expire_at": resp.Membership.ExpireAt,
		"is_active": resp.Membership.IsActive,
	}))
}

// BuyMembership 购买会员
func (h *CommerceHandler) BuyMembership(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req struct {
		Level   int32 `json:"level"`
		PayType int32 `json:"pay_type"`
	}
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request"))
		return
	}

	resp, err := h.commerceClient.BuyMembership(ctx, userID, req.Level, req.PayType)
	if err != nil {
		hlog.Errorf("BuyMembership: userID=%d level=%d failed: %v", userID, req.Level, err)
		c.JSON(500, hertz.Fail(500, "Failed to buy membership"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"id":        resp.Membership.Id,
		"user_id":   resp.Membership.UserId,
		"level":     resp.Membership.Level,
		"start_at":  resp.Membership.StartAt,
		"expire_at": resp.Membership.ExpireAt,
		"is_active": resp.Membership.IsActive,
	}))
}

// ListFlashSales 获取秒杀列表
func (h *CommerceHandler) ListFlashSales(ctx context.Context, c *app.RequestContext) {
	resp, err := h.commerceClient.ListFlashSales(ctx)
	if err != nil {
		hlog.Errorf("ListFlashSales: failed: %v", err)
		c.JSON(500, hertz.Fail(500, "Failed to list flash sales"))
		return
	}

	sales := make([]map[string]interface{}, 0, len(resp.FlashSales))
	for _, s := range resp.FlashSales {
		sales = append(sales, map[string]interface{}{
			"id":           s.Id,
			"name":         s.Name,
			"product_id":   s.ProductId,
			"flash_price":  s.FlashPrice,
			"stock":        s.Stock,
			"remain_stock": s.RemainStock,
			"start_at":     s.StartAt,
			"end_at":       s.EndAt,
		})
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"flash_sales": sales,
	}))
}

// FlashSalePurchase 秒杀购买
func (h *CommerceHandler) FlashSalePurchase(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	flashSaleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid flash_sale_id"))
		return
	}

	var req struct {
		ProductID int64 `json:"product_id"`
		PayType   int   `json:"pay_type"`
	}
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request"))
		return
	}

	resp, err := h.commerceClient.FlashSalePurchase(ctx, userID, flashSaleID, req.ProductID, int32(req.PayType))
	if err != nil {
		hlog.Errorf("FlashSalePurchase: userID=%d flashSaleID=%d failed: %v", userID, flashSaleID, err)
		c.JSON(500, hertz.Fail(500, "Failed to purchase"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"success":  resp.Success,
		"order_id": resp.OrderId,
		"message":  resp.Message,
	}))
}

// ListLotteries 获取抽奖列表
func (h *CommerceHandler) ListLotteries(ctx context.Context, c *app.RequestContext) {
	resp, err := h.commerceClient.ListLotteries(ctx)
	if err != nil {
		hlog.Errorf("ListLotteries: failed: %v", err)
		c.JSON(500, hertz.Fail(500, "Failed to list lotteries"))
		return
	}

	lotteries := make([]map[string]interface{}, 0, len(resp.Lotteries))
	for _, l := range resp.Lotteries {
		lotteries = append(lotteries, map[string]interface{}{
			"id":             l.Id,
			"name":           l.Name,
			"description":    l.Description,
			"entry_points":   l.EntryPoints,
			"start_time":     l.StartAt,
			"end_time":       l.EndAt,
			"remain_tickets": l.RemainTickets,
		})
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"lotteries": lotteries,
	}))
}

// ParticipateLottery 参与抽奖
func (h *CommerceHandler) ParticipateLottery(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	lotteryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid lottery_id"))
		return
	}

	resp, err := h.commerceClient.ParticipateLottery(ctx, userID, lotteryID)
	if err != nil {
		hlog.Errorf("ParticipateLottery: userID=%d lotteryID=%d failed: %v", userID, lotteryID, err)
		c.JSON(500, hertz.Fail(500, "Failed to participate"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"success": resp.Success,
		"won":     resp.Won,
		"prize":   resp.Prize,
		"message": resp.Message,
	}))
}
