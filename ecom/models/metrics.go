package models

type Metrics struct {
	TotalOrdersReceived   int64   `json:"total_orders_received"`
	AverageProcessingTime float64 `json:"average_processing_time"` // In seconds
	OrdersPending         int     `json:"orders_pending"`
	OrdersProcessing      int     `json:"orders_processing"`
	OrdersCompleted       int     `json:"orders_completed"`
}
