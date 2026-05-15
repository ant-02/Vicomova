package router

import (
	"vicomova/internal/commerce/interfaces/http/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz, commerceHandler *handler.CommerceHandler) {
	commerce := h.Group("/commerce")

	// Points/Wallet
	commerce.GET("/wallet", commerceHandler.GetWallet)
	commerce.POST("/wallet/recharge", commerceHandler.RechargePoints)

	// Sign-in
	commerce.POST("/sign-in", commerceHandler.DailySignIn)
	commerce.GET("/sign-in/info", commerceHandler.GetSignInInfo)

	// Products
	commerce.GET("/products", commerceHandler.ListProducts)
	commerce.GET("/products/:id", commerceHandler.GetProduct)

	// Orders
	commerce.POST("/orders", commerceHandler.CreateOrder)
	commerce.POST("/orders/:id/pay", commerceHandler.PayOrder)
	commerce.GET("/orders", commerceHandler.ListOrders)

	// Membership
	commerce.GET("/membership", commerceHandler.GetMembership)
	commerce.POST("/membership/buy", commerceHandler.BuyMembership)

	// Flash Sale
	commerce.GET("/flash-sales", commerceHandler.ListFlashSales)
	commerce.POST("/flash-sales/:id/purchase", commerceHandler.FlashSalePurchase)

	// Lottery
	commerce.GET("/lotteries", commerceHandler.ListLotteries)
	commerce.POST("/lotteries/:id/participate", commerceHandler.ParticipateLottery)
}
