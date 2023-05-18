package service

import "github.com/aidos-dev/habit-tracker/internal/repository"

type AdminService struct {
	repo repository.Admin
}

func NewAdminService(repo repository.Admin) Admin {
	return &AdminService{repo: repo}
}

func (s *AdminService) AssignRole(userId int, role string) (int, error) {
	return s.repo.AssignRole(userId, role)
}
