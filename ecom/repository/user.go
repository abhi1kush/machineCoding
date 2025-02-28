package repository

import (
	"ecom.com/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	GetAllUsers() ([]models.User, error)
}
