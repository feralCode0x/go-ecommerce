package http

import (
	//"go-ecommerce/docs"
	"go-ecommerce/pkgs/mail"
	"go-ecommerce/pkgs/middlewares"
	"go-ecommerce/pkgs/minio"
	"go-ecommerce/pkgs/token"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go-ecommerce/db"
	"go-ecommerce/pkgs/logger"
	"go-ecommerce/pkgs/validation"
	"net/http"

	"go-ecommerce/configs"
	"go-ecommerce/pkgs/redis"

	cartHttp "go-ecommerce/internals/cart/controller/http"
	orderHttp "go-ecommerce/internals/order/controller/http"
	productHttp "go-ecommerce/internals/product/controller/http"
	userHttp "go-ecommerce/internals/user/controller/http"
)

type Server struct {
	engine      *gin.Engine
	cfg         *configs.Config
	validator   validation.Validation
	db          db.IDatabase
	minioClient minio.IUploadService
	cache       redis.IRedis
	tokenMarker token.IMarker
	mailer      mail.IMailer
	enforcer    *casbin.Enforcer
}

func NewServer(
	validator validation.Validation,
	db db.IDatabase,
	minioClient minio.IUploadService,
	cache redis.IRedis,
	tokenMarker token.IMarker,
	mailer mail.IMailer,
	enforcer *casbin.Enforcer,
) *Server {
	return &Server{
		engine:      gin.Default(),
		cfg:         configs.GetConfig(),
		validator:   validator,
		db:          db,
		minioClient: minioClient,
		cache:       cache,
		tokenMarker: tokenMarker,
		mailer:      mailer,
		enforcer:    enforcer,
	}
}

func (s Server) Run() error {
	_ = s.engine.SetTrustedProxies(nil)
	if s.cfg.Environment == configs.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	s.engine.Use(func(c *gin.Context) {
		c.Set("enforcer", s.enforcer)
		c.Next()
	})

	s.engine.Use(middlewares.PrometheusMiddleware())
	s.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	s.engine.Use(middlewares.CorsMiddleware())

	if err := s.MapRoutes(); err != nil {
		logger.Fatalf("MapRoutes Error: %v", err)
	}

	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to Ecommerce Clean Architecture"})
	})

	//Start http server
	logger.Info("HTTP server is listening on PORT: ", s.cfg.HttpPort)
	//if err := s.engine.Run(fmt.Sprintf(":8080"); err != nil {
	if err := s.engine.Run(fmt.Sprintf("0.0.0.0:%d", s.cfg.HttpPort)); err != nil {
		logger.Fatalf("Running HTTP server: %v", err)
	}

	 if err := s.MapRoutes(); err != nil {
        logger.Fatalf("Error mapeando rutas: %v", err)
    }

    // 2. Imprimir confirmación (Si no ves esto al arrancar, MapRoutes no se ejecutó)
    fmt.Println("Rutas mapeadas con éxito")

    // 3. Iniciar servidor
    return s.engine.Run(fmt.Sprintf("0.0.0.0:%d", s.cfg.HttpPort))
}

func (s Server) GetEngine() *gin.Engine {
	return s.engine
}

//	@title			Ecommerce Clean Architecture Swagger API
//	@version		1.0
// @host        	localhost:8080
// @BasePath    	/api/v1
//	@description	Swagger API for Go Clean Architecture.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Tran Phuoc Anh Quoc
//	@contact.email	anhquoc18092003@gmail.com

//	@license.name	MIT
//	@license.url	https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func (s Server) MapRoutes() error {
	routesV1 := s.engine.Group("/api/v1")
	userHttp.Routes(routesV1, s.db, s.validator, s.minioClient, s.cache, s.mailer, s.tokenMarker)
	productHttp.Routes(routesV1, s.db, s.validator, s.minioClient, s.cache, s.tokenMarker)
	cartHttp.Routes(routesV1, s.db, s.validator, s.cache, s.tokenMarker)
	orderHttp.Routes(routesV1, s.db, s.validator, s.cache, s.tokenMarker)
	return nil
}