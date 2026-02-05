package http

import (
	"go-ecommerce/db"
	"go-ecommerce/internals/order/repository"
	"go-ecommerce/internals/order/usecase"
	productRepo "go-ecommerce/internals/product/repository"
	"go-ecommerce/pkgs/middlewares"
	"go-ecommerce/pkgs/redis"
	"go-ecommerce/pkgs/token"
	"go-ecommerce/pkgs/validation"

	"github.com/gin-gonic/gin"
)

func Routes(
	r *gin.RouterGroup,
	sqlDB db.IDatabase,
	validator validation.Validation,
	cache redis.IRedis,
	token token.IMarker,
) {
	productRepository := productRepo.NewProductRepository(sqlDB)
	orderRepository := repository.NewOrderRepository(sqlDB)
	orderUsecase := usecase.NewOrderUseCase(validator, orderRepository, productRepository)
	orderHandler := NewOrderHandler(orderUsecase)

	authMiddleware := middlewares.NewAuthMiddleware(token, cache).TokenAuth()

	orderRoute := r.Group("/orders", authMiddleware)
	{
		orderRoute.POST("", orderHandler.PlaceOrder)
		orderRoute.GET("", orderHandler.GetOrders)
		orderRoute.GET("/:id", orderHandler.GetOrderByID)
		orderRoute.PUT("/:id/:status", orderHandler.UpdateOrder)
	}
}