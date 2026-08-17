package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flowershop/backend/internal/config"
	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/database"
	"github.com/flowershop/backend/internal/handler"
	"github.com/flowershop/backend/internal/repository"
	"github.com/flowershop/backend/internal/router"
	"github.com/flowershop/backend/internal/service"
	"github.com/flowershop/backend/internal/util"
)

func main() {
	cfg := config.Load()
	logger := util.InitLogger(cfg.AppEnv)
	logger.Info(constants.LogServerStarted, "app", cfg.AppName, "env", cfg.AppEnv, "port", cfg.ServerPort)

	db, err := database.Connect(cfg)
	if err != nil {
		logger.Error("connect database failed", "err", err)
		os.Exit(1)
	}
	database.LogDBStats(logger, cfg)
	tables, err := database.Migrate(db)
	if err != nil {
		logger.Error("migrate database failed", "err", err)
		os.Exit(1)
	}
	logger.Info(constants.LogDBMigrated, "tables", tables)

	rdb, err := database.ConnectRedis(cfg)
	if err != nil {
		logger.Error("connect redis failed", "err", err)
		os.Exit(1)
	}
	database.LogRedis(logger, cfg)

	minioClient, err := database.ConnectMinIO(cfg)
	if err != nil {
		logger.Error("connect minio failed", "err", err)
		os.Exit(1)
	}
	database.LogMinIO(logger, cfg)

	if err := database.Seed(db, minioClient, cfg.MinIOBucket, logger); err != nil {
		logger.Error("seed database failed", "err", err)
		os.Exit(1)
	}

	// 仓储
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	orderItemRepo := repository.NewOrderItemRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	logisticsRepo := repository.NewLogisticsRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	couponRepo := repository.NewCouponRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	// 服务
	auditSvc := service.NewAuditService(auditRepo)
	auditQuerySvc := service.NewAuditQueryService(auditRepo)
	jwtCfg := &util.JWTConfig{Secret: cfg.JWTSecret, AccessExpire: cfg.JWTAccessExpire}
	userSvc := service.NewUserService(userRepo, auditSvc, jwtCfg, logger)
	categorySvc := service.NewCategoryService(categoryRepo, auditSvc, logger)
	productSvc := service.NewProductService(productRepo, categoryRepo, auditSvc, rdb, logger)
	cartSvc := service.NewCartService(cartRepo, productRepo, auditSvc, logger)
	orderSvc := service.NewOrderService(db, orderRepo, orderItemRepo, paymentRepo, logisticsRepo, cartRepo, productRepo, couponRepo, auditSvc, logger)
	reviewSvc := service.NewReviewService(db, reviewRepo, orderItemRepo, orderRepo, productRepo, userRepo, auditSvc, logger)
	couponSvc := service.NewCouponService(db, couponRepo, auditSvc, logger)
	uploadSvc := service.NewUploadService(minioClient, cfg.MinIOBucket, logger)
	statsSvc := service.NewStatsService(orderRepo, productRepo, userRepo)

	// 处理器
	deps := &router.Deps{
		UserHandler:     handler.NewUserHandler(userSvc),
		CategoryHandler: handler.NewCategoryHandler(categorySvc),
		ProductHandler:  handler.NewProductHandler(productSvc),
		CartHandler:     handler.NewCartHandler(cartSvc),
		OrderHandler:    handler.NewOrderHandler(orderSvc),
		ReviewHandler:   handler.NewReviewHandler(reviewSvc),
		CouponHandler:   handler.NewCouponHandler(couponSvc),
		AuditHandler:    handler.NewAuditHandler(auditQuerySvc),
		UploadHandler:   handler.NewUploadHandler(uploadSvc),
		StatsHandler:    handler.NewStatsHandler(statsSvc),
	}

	engine := router.New(cfg, db, logger, deps)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server listen failed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(constants.LogServerStopped, "reason", "signal")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
