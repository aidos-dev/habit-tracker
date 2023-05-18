package service

import (
	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/aidos-dev/habit-tracker/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) User {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}
