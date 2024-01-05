package service

import (
	"fmt"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
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

func (s *AdminUserRewardService) RemoveFromUser(userId, habitId, rewardId int) error {
	return s.repo.RemoveFromUser(userId, habitId, rewardId)
}

func (s *AdminUserRewardService) UpdateUserReward(userId, habitId, rewardId int, input models.UpdateUserRewardInput) error {
	const op = "service.admin_user_reward_service.UpdateUserReward"

	if err := input.Validate(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return s.repo.UpdateUserReward(userId, habitId, rewardId, input)
}
