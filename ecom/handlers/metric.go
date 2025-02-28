package handlers

import (
	"net/http"

	"ecom.com/services"

	"github.com/gin-gonic/gin"
)

type MetricHandler struct {
	Service *services.Metric
}

func NewMetricHandler(service *services.Metric) *MetricHandler {
	return &MetricHandler{Service: service}
}

// MetricsHandler handles GET /metrics requests.
// It calls the metrics service to retrieve metrics data and returns it as JSON.
func (h *MetricHandler) GetMetricsHandler(c *gin.Context) {
	metrics, err := h.Service.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch metrics"})
		return
	}
	c.JSON(http.StatusOK, metrics)
}
