package handlers

import (
	"net/http"

	"ecom.com/models"
	"ecom.com/services"

	"github.com/gin-gonic/gin"
)

func CreateOrderHandler(c *gin.Context) {
	req := models.OrderRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID, err := services.CreateOrder(req.UserID, req.ItemIDs, req.TotalAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	resp := models.OrderResponse{Message: "Order created", OrderId: orderID}
	c.JSON(http.StatusOK, resp)
}

func GetOrderStatusHandler(c *gin.Context) {
	orderID := c.Param("order_id")
	status, err := services.GetOrderStatus(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order status"})
		return
	}

	if status == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	resp := models.OrderStatusResponse{OrderId: orderID, Status: status}
	c.JSON(http.StatusOK, resp)
}
