package models

type Order struct {
	OrderID     string  `json:"order_id"`
	UserID      string  `json:"user_id"`
	ItemIDs     string  `json:"item_ids"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}
