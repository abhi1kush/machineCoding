package handlers

import (
	"net/http"

	"ecom.com/services"

	"github.com/gin-gonic/gin"
)

// MetricsHandler handles GET /metrics requests.
// It calls the metrics service to retrieve metrics data and returns it as JSON.
func MetricsHandler(c *gin.Context) {
	metrics, err := services.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch metrics"})
		return
	}
	c.JSON(http.StatusOK, metrics)
}
