package handlers

import (
	"net/http"

	"ecom.com/common"
	"ecom.com/services"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Service *services.Order
}

func NewOrderHandler(service *services.Order) *OrderHandler {
	return &OrderHandler{Service: service}
}

func (h *OrderHandler) CreateOrderHandler(c *gin.Context) {
	req := common.OrderRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID, err := h.Service.CreateOrder(req.UserID, req.ItemIDs, req.TotalAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	resp := &common.OrderAckResponse{Message: "Order created", OrderID: orderID}
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetOrderHandler(c *gin.Context) {
	orderID := c.Param("order_id")
	order, err := h.Service.GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order status"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrderStatusHandler(c *gin.Context) {
	orderID := c.Param("order_id")
	status, err := h.Service.GetOrderStatus(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order status"})
		return
	}

	if status == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	resp := common.OrderStatusResponse{OrderId: orderID, Status: status}
	c.JSON(http.StatusOK, resp)
}
