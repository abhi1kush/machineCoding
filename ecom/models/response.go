package models

type OrderResponse struct {
	Message string `json:"message"`
	OrderId string `json:"order_id"`
}

type OrderStatusResponse struct {
	OrderId string `json:"order_id"`
	Status  string `json:"status"`
}
