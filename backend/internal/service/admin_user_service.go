package service

import (
	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
)

type AdminUserService struct {
	repo repository.AdminUser
	User
}

func NewAdminUserService(repo repository.AdminUser) AdminUser {
	return &AdminUserService{repo: repo}
}

func (r *AdminUserService) GetAllUsers() ([]models.GetUser, error) {
	return r.repo.GetAllUsers()
}

func (r *AdminUserService) GetUserById(userId int) (models.GetUser, error) {
	return r.repo.GetUserById(userId)
}
