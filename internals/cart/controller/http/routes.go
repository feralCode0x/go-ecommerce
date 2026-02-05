package http

import (
	"go-ecommerce/db"
	"go-ecommerce/internals/cart/usecase"
	"go-ecommerce/pkgs/middlewares"
	"go-ecommerce/pkgs/redis"
	"go-ecommerce/pkgs/token"
	"go-ecommerce/pkgs/validation"

	"github.com/gin-gonic/gin"

	cartRepo "go-ecommerce/internals/cart/repository"
	productRepo "go-ecommerce/internals/product/repository"
)

func Routes(
	r *gin.RouterGroup,
	sqlDB db.IDatabase,
	validator validation.Validation,
	cache redis.IRedis,
	token token.IMarker,
) {

	cartRepository := cartRepo.NewCartRepository(sqlDB)
	productRepository := productRepo.NewProductRepository(sqlDB)
	cartUseCase := usecase.NewCartUseCase(validator, cartRepository, productRepository)
	cartHandler := NewCartHandler(cartUseCase)

	authMiddleware := middlewares.NewAuthMiddleware(token, cache).TokenAuth()

	cartRoute := r.Group("/carts", authMiddleware)
	{
		cartRoute.GET("/:userID", cartHandler.GetCart)
		cartRoute.POST("/:userID", cartHandler.AddProductToCart)
		cartRoute.PUT("/cart-line/:userID", cartHandler.UpdateCartLine)
		cartRoute.DELETE("/:userID", cartHandler.RemoveProductToCart)
	}
}