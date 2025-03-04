package routes

import (
	"ecom.com/handlers"
	"ecom.com/middleware"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	OrderHandler  *handlers.OrderHandler
	MetricHandler *handlers.MetricHandler
}

func RegisterOrderRoutes(router *gin.RouterGroup, orderHandler *handlers.OrderHandler) {
	orderRoutes := router.Group("/orders")
	{
		orderRoutes.POST("", middleware.LoggerMiddleware(), orderHandler.CreateOrderHandler)
		orderRoutes.GET("/:id", orderHandler.GetOrderHandler)
		orderRoutes.GET("/status/:id", orderHandler.GetOrderStatusHandler)
	}
}

// RegisterRoutes initializes all API routes with middleware and versioning
func RegisterRoutes(router *gin.Engine, cfg *RouterConfig) {
	router.Use(middleware.LoggerMiddleware()) // Apply logging middleware globally
	router.GET("health", handlers.HealthChecksHandler)
	apiV1 := router.Group("/api/v1") // Version 1 API group
	{
		RegisterOrderRoutes(apiV1, cfg.OrderHandler)
		apiV1.GET("/metrics", cfg.MetricHandler.GetMetricsHandler)
	}
}
