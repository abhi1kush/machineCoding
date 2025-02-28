package routes

import (
	"ecom.com/handlers"
	"ecom.com/middleware"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	UserHandler   *handlers.UserHandler
	OrderHandler  *handlers.OrderHandler
	MetricHandler *handlers.MetricHandler
}

// RegisterUserRoutes initializes user-related routes
func RegisterUserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", middleware.LoggerMiddleware(), userHandler.CreateUser) // POST /api/v1/users
		userRoutes.GET("/:id", userHandler.GetUserByID)                            // GET /api/v1/users/:id
	}
}

func RegisterOrderRoutes(router *gin.RouterGroup, orderHandler *handlers.OrderHandler) {
	orderRoutes := router.Group("/orders")
	{
		orderRoutes.POST("", middleware.LoggerMiddleware(), orderHandler.CreateOrderHandler)
		orderRoutes.GET("/:id", orderHandler.GetOrderHandler)
		orderRoutes.GET("/status/:id", orderHandler.GetOrderStatusHandler)
	}
}

func RegisterMetricRoutes(router *gin.RouterGroup, metricsHandler *handlers.MetricHandler) {
	metricRoutes := router.Group("/metrics")
	{
		metricRoutes.POST("", middleware.LoggerMiddleware(), metricsHandler.CreateMetricsHandler)
		metricRoutes.GET("/:id", metricsHandler.GetMetricsHandler)
	}
}

// RegisterRoutes initializes all API routes with middleware and versioning
func RegisterRoutes(router *gin.Engine, cfg *RouterConfig) {
	router.Use(middleware.LoggerMiddleware()) // Apply logging middleware globally

	apiV1 := router.Group("/api/v1") // Version 1 API group
	{
		RegisterUserRoutes(apiV1, cfg.UserHandler)
		RegisterOrderRoutes(apiV1, cfg.OrderHandler)
		RegisterMetricRoutes(apiV1, cfg.MetricHandler)
	}
}
