package service

import (
	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/aidos-dev/habit-tracker/internal/repository"
)

type AdminUserRewardService struct {
	repo repository.AdminUserReward
	Reward
}

func NewAdminUserRewardService(repo repository.AdminUserReward) AdminUserReward {
	return &AdminUserRewardService{repo: repo}
}

func (s *AdminUserRewardService) AssignReward(userId, habitId, rewardId int) (int, error) {
	return s.repo.AssignReward(userId, habitId, rewardId)
}

func (s *AdminUserRewardService) RemoveFromUser(userId, rewardId int) error {
	return s.repo.RemoveFromUser(userId, rewardId)
}

func (s *AdminUserRewardService) UpdateUserReward(userId, rewardId int, input models.UpdateUserRewardInput) error {
	return s.repo.UpdateUserReward(userId, rewardId, input)
}
