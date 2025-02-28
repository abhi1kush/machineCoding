package repository

import (
	"ecom.com/models"
)

type OrderRepositoryI interface {
	CreateOrder(order *models.Order) error
	UpdateOrderStatus(orderId string, status string) error
	GetOrderByID(id string) (*models.Order, error)
}
