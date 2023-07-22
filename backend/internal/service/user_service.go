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
	// const op = "service.user_service.CreateUser"

	user.Password = generatePasswordHash(user.Password)

	// if err := user.Validate(); err != nil {
	// 	return 0, fmt.Errorf("%s: %w", op, err)
	// }

	return s.repo.CreateUser(user)
}

func (s *UserService) GetUser(username, password string) (models.User, error) {
	return s.repo.GetUser(username, password)
}

func (r *UserService) GetUserByTgUsername(TGusername string) (models.GetUser, error) {
	return r.repo.GetUserByTgUsername(TGusername)
}

func (r *UserService) GetUserById(userId int) (models.GetUser, error) {
	return r.repo.GetUserById(userId)
}

func (r *UserService) GetAllUsers() ([]models.GetUser, error) {
	return r.repo.GetAllUsers()
}

func (s *UserService) DeleteUser(userId int) (int, error) {
	return s.repo.DeleteUser(userId)
}
