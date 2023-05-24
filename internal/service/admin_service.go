package service

import "github.com/aidos-dev/habit-tracker/internal/repository"

type AdminService struct {
	repo repository.Admin
	AdminUser
	AdminRole
	AdminReward
	AdminUserReward
}

func NewAdminService(repo repository.Admin) Admin {
	return &AdminService{repo: repo}
}
