package services

import (
	"ecom.com/models"
	"ecom.com/repository"
)

type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.Repo.CreateUser(user)
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAllUsers()
}
