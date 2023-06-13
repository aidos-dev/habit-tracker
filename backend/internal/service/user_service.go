package service

import (
	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
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

func (s *UserService) GetUser(username, password string) (models.User, error) {
	return s.repo.GetUser(username, password)
}

func (s *UserService) DeleteUser(userId int) (int, error) {
	return s.repo.DeleteUser(userId)
}
