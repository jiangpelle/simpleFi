package routes

import (
	"net/http"

	"defi-backend/handlers"
	"defi-backend/middleware"
	"defi-backend/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Router struct {
	userHandler *handlers.UserHandler
	defiHandler *handlers.DefiHandler
	logger      *zap.Logger
}

func NewRouter(userService *services.UserService, defiService *services.DefiService, logger *zap.Logger) *Router {
	return &Router{
		userHandler: handlers.NewUserHandler(userService),
		defiHandler: handlers.NewDefiHandler(defiService),
		logger:      logger,
	}
}

func (r *Router) SetupRouter() *gin.Engine {
	router := gin.Default()

	// 使用日志中间件
	router.Use(middleware.LoggerMiddleware(r.logger))

	// API 路由组
	api := router.Group("/api")
	{
		// 健康检查
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "Service is running",
			})
		})

		// 用户相关路由
		user := api.Group("/user")
		{
			user.POST("/register", r.userHandler.Register)
			user.POST("/login", r.userHandler.Login)
			user.GET("/profile", middleware.AuthMiddleware(), r.userHandler.GetProfile)
		}

		// DeFi 相关路由
		defi := api.Group("/defi")
		{
			// DEX 路由
			dex := defi.Group("/dex")
			{
				dex.POST("/swap", middleware.AuthMiddleware(), r.defiHandler.SwapTokens)
				dex.GET("/pairs", r.defiHandler.GetTradingPairs)
				dex.GET("/price/:pair", r.defiHandler.GetTokenPrice)
			}

			// 借贷路由
			lending := defi.Group("/lending")
			{
				lending.POST("/deposit", middleware.AuthMiddleware(), r.defiHandler.Deposit)
				lending.POST("/borrow", middleware.AuthMiddleware(), r.defiHandler.Borrow)
				lending.GET("/positions", middleware.AuthMiddleware(), r.defiHandler.GetPositions)
			}

			// 挖矿路由
			farming := defi.Group("/farming")
			{
				farming.POST("/stake", middleware.AuthMiddleware(), r.defiHandler.StakeTokens)
				farming.POST("/unstake", middleware.AuthMiddleware(), r.defiHandler.UnstakeTokens)
				farming.GET("/rewards", middleware.AuthMiddleware(), r.defiHandler.GetRewards)
			}
		}
	}

	return router
}
