package service

import (
	"github.com/aidos-dev/habit-tracker/internal/repository"
)

type AdminRoleService struct {
	repo repository.AdminRole
}

func NewAdminRoleService(repo repository.AdminRole) AdminRole {
	return &AdminRoleService{repo: repo}
}

func (r *AdminRoleService) AssignRole(userId int, role string) (int, error) {
	return r.repo.AssignRole(userId, role)
}
