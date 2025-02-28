package common

type OrderRequest struct {
	UserID      string  `json:"user_id"`
	ItemIDs     string  `json:"item_ids"`
	TotalAmount float64 `json:"total_amount"`
}

type MetricRequest struct {
	OrderId        string
	ProcessingTime int
}
