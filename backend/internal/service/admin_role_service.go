package service

import (
	"fmt"

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
	const op = "service.admin_role_service.AssignRole"

	if err := role.Validate(); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return r.repo.AssignRole(userId, role)
}
