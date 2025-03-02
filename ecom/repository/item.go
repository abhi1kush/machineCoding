package repository

import (
	"ecom.com/models"
)

type ItemRepositoryI interface {
	CreateItem(item *models.Item) error
	GetItem(id string) (*models.Item, error)
	GetItemsByOrderId(id string) ([]models.Item, error)
	RemoveItem(itemId string, orderId string) error
}
