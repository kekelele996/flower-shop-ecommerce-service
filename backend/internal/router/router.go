package router

import (
	"log/slog"

	"github.com/flowershop/backend/internal/config"
	"github.com/flowershop/backend/internal/handler"
	"github.com/flowershop/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Router 聚合所有模块路由。
type Router struct {
	cfg    *config.Config
	db     *gorm.DB
	logger *slog.Logger
	engine *gin.Engine
}

// Deps 路由依赖装配。
type Deps struct {
	UserHandler     *handler.UserHandler
	CategoryHandler *handler.CategoryHandler
	ProductHandler  *handler.ProductHandler
	CartHandler     *handler.CartHandler
	OrderHandler    *handler.OrderHandler
	ReviewHandler   *handler.ReviewHandler
	CouponHandler   *handler.CouponHandler
	AuditHandler    *handler.AuditHandler
	UploadHandler   *handler.UploadHandler
	StatsHandler    *handler.StatsHandler
}

// New 创建并注册全部路由。
func New(cfg *config.Config, db *gorm.DB, logger *slog.Logger, deps *Deps) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestID())
	engine.Use(middleware.RequestLogger(logger))
	engine.Use(middleware.CORS(cfg.CORSOrigins))
	engine.Use(middleware.ErrorHandler())

	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "app": cfg.AppName})
	})

	v1 := engine.Group("/api/v1")

	// 公开路由
	auth := v1.Group("/auth")
	auth.POST("/register", deps.UserHandler.Register)
	auth.POST("/login", deps.UserHandler.Login)

	v1.GET("/categories", deps.CategoryHandler.List)
	v1.GET("/products", deps.ProductHandler.List)
	v1.GET("/products/recommendations", deps.ProductHandler.Recommendations)
	v1.GET("/products/:id", deps.ProductHandler.Detail)
	v1.POST("/products/:id/view", deps.ProductHandler.View)
	v1.GET("/reviews", deps.ReviewHandler.List)

	// 需要登录
	user := v1.Group("")
	user.Use(middleware.Auth(cfg.JWTSecret))
	user.GET("/auth/profile", deps.UserHandler.Profile)
	user.PUT("/auth/profile", deps.UserHandler.UpdateProfile)
	user.PUT("/auth/password", deps.UserHandler.ChangePassword)

	user.GET("/cart", deps.CartHandler.List)
	user.POST("/cart", deps.CartHandler.Add)
	user.PUT("/cart/:id", deps.CartHandler.Update)
	user.DELETE("/cart/:id", deps.CartHandler.Remove)

	user.POST("/orders", deps.OrderHandler.Checkout)
	user.GET("/orders", deps.OrderHandler.List)
	user.GET("/orders/:id", deps.OrderHandler.Detail)
	user.POST("/orders/:id/pay", deps.OrderHandler.Pay)
	user.POST("/orders/:id/cancel", deps.OrderHandler.Cancel)
	user.POST("/orders/:id/complete", deps.OrderHandler.Complete)
	user.GET("/orders/:id/logistics", deps.OrderHandler.Logistics)

	user.POST("/reviews", deps.ReviewHandler.Create)
	user.GET("/coupons", deps.CouponHandler.Mine)
	user.GET("/coupons/available", deps.CouponHandler.Available)
	user.POST("/coupons/:id/claim", deps.CouponHandler.Claim)

	user.POST("/uploads", deps.UploadHandler.Upload)

	// 管理员路由
	admin := v1.Group("/admin")
	admin.Use(middleware.Auth(cfg.JWTSecret), middleware.RBAC("ADMIN"))
	admin.GET("/products", deps.ProductHandler.AdminList)
	admin.POST("/categories", deps.CategoryHandler.Create)
	admin.PUT("/categories/:id", deps.CategoryHandler.Update)
	admin.DELETE("/categories/:id", deps.CategoryHandler.Delete)
	admin.POST("/products", deps.ProductHandler.Create)
	admin.PUT("/products/:id", deps.ProductHandler.Update)
	admin.PUT("/products/:id/status", deps.ProductHandler.ChangeStatus)
	admin.GET("/orders", deps.OrderHandler.List)
	admin.POST("/orders/:id/ship", deps.OrderHandler.Ship)
	admin.POST("/orders/:id/complete", deps.OrderHandler.CompleteAdmin)
	admin.POST("/orders/:id/pay", deps.OrderHandler.Pay)
	admin.GET("/orders/:id/logistics", deps.OrderHandler.Logistics)
	admin.GET("/reviews", deps.ReviewHandler.List)
	admin.POST("/reviews/:id/reply", deps.ReviewHandler.Reply)
	admin.POST("/coupons", deps.CouponHandler.CreateTemplate)
	admin.GET("/audit-logs", deps.AuditHandler.List)
	admin.GET("/stats", deps.StatsHandler.Stats)

	return engine
}
