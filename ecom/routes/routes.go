package routes

import (
	"ecom.com/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/healthcheck", handlers.HealthChecksHandler)
	r.POST("/orders", handlers.CreateOrderHandler)
	r.GET("/order/:order_id", handlers.GetOrderStatusHandler)
	r.GET("/metrics", handlers.MetricsHandler)
}
