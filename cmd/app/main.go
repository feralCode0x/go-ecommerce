package main

import (
	"go-ecommerce/configs"
	"go-ecommerce/db"
	"go-ecommerce/pkgs/casbin"
	"go-ecommerce/pkgs/logger"
	"go-ecommerce/pkgs/mail"
	"go-ecommerce/pkgs/minio"
	"go-ecommerce/pkgs/redis"
	"go-ecommerce/pkgs/token"
	"go-ecommerce/pkgs/validation"
	"sync"

	cartEntity "go-ecommerce/internals/cart/entity"
	orderEntity "go-ecommerce/internals/order/entity"
	productEntity "go-ecommerce/internals/product/entity"
	httpServer "go-ecommerce/internals/server/http"
	userEntity "go-ecommerce/internals/user/entity"
)

var wg sync.WaitGroup

func main() {
	cfg := configs.LoadConfig()
	logger.Initialize(cfg.Environment)

	database, err := db.NewDatabase(cfg.DatabaseURI)
	if err != nil {
		logger.Fatal("Cannot connect to database", err)
	}

	enforcer, err := casbin.InitCasbinEnforcer(database)
	if err != nil {
		logger.Fatal(err)
	}

	if err := database.AutoMigrate(
		&userEntity.User{},
		&productEntity.Product{},
		&orderEntity.Order{},
		&orderEntity.OrderLine{},
		&cartEntity.Cart{},
		&cartEntity.CartLine{}); err != nil {
		logger.Fatal("Database migration fail", err)
	}

	validator := validation.New()

	//minio
	minioClient, err := minio.NewMinioClient(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		cfg.MinioBaseurl,
		cfg.MinioUseSSL,
	)
	if err != nil {
		logger.Fatalf("Failed to connect to MinIO: %s", err)
	}

	//mailer
	mailer := mail.NewMailer(
		cfg.MailHost,
		cfg.MailPort,
		cfg.MailUser,
		cfg.MailPassword,
		cfg.MailFrom,
	)

	//cache
	cache := redis.New(redis.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	})

	//token
	tokenMaker, err := token.NewJTWMarker()
	if err != nil {
		logger.Fatal(err)
	}

	httpSvr := httpServer.NewServer(validator, database, minioClient, cache, tokenMaker, mailer, enforcer)

	wg.Add(1)

	// Run HTTP server
	go func() {
		defer wg.Done()
		if err := httpSvr.Run(); err != nil {
			logger.Fatal("Running HTTP server error:", err)
		}
	}()

	wg.Wait()
}