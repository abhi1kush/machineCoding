package handlers

import (
	"net/http"

	"ecom.com/common"
	"ecom.com/models"
	"ecom.com/services"

	"github.com/gin-gonic/gin"
)

type MetricHandler struct {
	Service *services.Metric
}

func NewMetricHandler(service *services.Metric) *MetricHandler {
	return &MetricHandler{Service: service}
}

func (h *MetricHandler) CreateMetricsHandler(c *gin.Context) {
	req := common.MetricRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metric := &models.Metric{
		OrderId:        req.OrderId,
		ProcessingTime: float64(req.ProcessingTime),
	}
	err := h.Service.Repo.CreateMetric(metric)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	c.Status(http.StatusOK)
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
