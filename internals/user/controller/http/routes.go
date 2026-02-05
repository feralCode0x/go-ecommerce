package http

import (
	"go-ecommerce/db"
	"go-ecommerce/internals/user/repository"
	"go-ecommerce/internals/user/usecase"
	"go-ecommerce/pkgs/mail"
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
	mailer mail.IMailer,
	token token.IMarker,
) {
	userRepository := repository.NewUserRepository(sqlDB)
	userUseCase := usecase.NewUserUseCase(validator, userRepository, minioClient, cache, mailer, token)
	userHandler := NewAuthHandler(userUseCase)

	authMiddleware := middlewares.NewAuthMiddleware(token, cache).TokenAuth()

	authRouter := r.Group("/auth")
	{
		authRouter.POST("/signup", userHandler.SignUp)
		authRouter.POST("/signin", userHandler.SignIn)
		authRouter.POST("/signout", authMiddleware, userHandler.SignOut)
	}
	userRouter := r.Group("/users").Use(authMiddleware)
	{
		userRouter.GET("", middlewares.AuthorizePolicy("users", "read"), userHandler.GetUsers)
		userRouter.GET("/:id", userHandler.GetUser)
		userRouter.DELETE("/:id", middlewares.AuthorizePolicy("users", "delete"), userHandler.DeleteUser)
	}
}
