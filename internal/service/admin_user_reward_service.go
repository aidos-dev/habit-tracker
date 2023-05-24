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

func (s *AdminUserRewardService) AssignReward(userId int, rewardId int, habitId int) (int, error) {
	return s.repo.AssignReward(userId, rewardId, habitId)
}

func (s *AdminUserRewardService) RemoveFromUser(userId, rewardId int) error {
	return s.repo.RemoveFromUser(userId, rewardId)
}

func (s *AdminUserRewardService) UpdateUserReward(userId, rewardId int, input models.UpdateUserRewardInput) error {
	return s.repo.UpdateUserReward(userId, rewardId, input)
}
