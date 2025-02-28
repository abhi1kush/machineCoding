package cache

type CacheI interface {
	SetOrderStatus(orderID, status string) error
	GetOrderStatus(orderID string) (string, error)
}
