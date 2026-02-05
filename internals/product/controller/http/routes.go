package http

import (
	"go-ecommerce/db"
	"go-ecommerce/internals/product/repository"
	"go-ecommerce/internals/product/usecase"
	"go-ecommerce/pkgs/middlewares"
	"go-ecommerce/pkgs/minio"
	"go-ecommerce/pkgs/redis"
	"go-ecommerce/pkgs/token"
	"go-ecommerce/pkgs/validation"
	"github.com/gin-gonic/gin"
)

func Routes(
	r *gin.RouterGroup,
	sqlDB db.IDatabase,
	validator validation.Validation,
	minioClient minio.IUploadService,
	cache redis.IRedis,
	token token.IMarker,
) {
	productRepository := repository.NewProductRepository(sqlDB)
	productUseCase := usecase.NewProductUseCase(validator, productRepository, minioClient)
	productHandler := NewProductHandler(productUseCase, cache)

	authMiddleware := middlewares.NewAuthMiddleware(token, cache).TokenAuth()

	productRoute := r.Group("/products").Use(authMiddleware)
	{
	productRoute.GET("", productHandler.GetProducts)
	productRoute.GET("/:id", productHandler.GetProduct)
	productRoute.POST("", middlewares.AuthorizePolicy("products", "write"), productHandler.CreateProduct)
	productRoute.PUT("/:id", middlewares.AuthorizePolicy("products", "write"), productHandler.UpdateProduct)
	productRoute.DELETE("/:id", middlewares.AuthorizePolicy("products", "delete"), productHandler.DeleteProduct)
	}
}