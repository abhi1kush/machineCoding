package common

type OrderAckResponse struct {
	Message string `json:"message"`
	OrderID string `json:"order_id"`
}

type OrderResponse struct {
	OrderID     string   `json:"order_id"`
	UserID      string   `json:"user_id"`
	ItemIDs     []string `json:"item_ids"`
	TotalAmount float64  `json:"total_amount"`
	Status      string   `json:"status"`
}

type OrderStatusResponse struct {
	OrderId string `json:"order_id"`
	Status  string `json:"status"`
}

type Metrics struct {
	TotalOrdersReceived   int64   `json:"total_orders_received"`
	AverageProcessingTime float64 `json:"average_processing_time"` // In seconds
	OrdersPending         int     `json:"orders_pending"`
	OrdersProcessing      int     `json:"orders_processing"`
	OrdersCompleted       int     `json:"orders_completed"`
}
