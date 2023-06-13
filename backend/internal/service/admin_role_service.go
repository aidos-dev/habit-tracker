package service

import (
	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
)

type AdminRoleService struct {
	repo repository.AdminRole
}

func NewAdminRoleService(repo repository.AdminRole) AdminRole {
	return &AdminRoleService{repo: repo}
}

func (r *AdminRoleService) AssignRole(userId int, role models.UpdateRoleInput) (int, error) {
	return r.repo.AssignRole(userId, role)
}
