package service

import (
	"fmt"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
)

type AdminRewardService struct {
	repo repository.AdminReward
}

func NewAdminRewardService(repo repository.AdminReward) AdminReward {
	return &AdminRewardService{repo: repo}
}

func (s *AdminRewardService) Create(reward models.Reward) (int, error) {
	return s.repo.Create(reward)
}

func (s *AdminRewardService) GetById(rewardId int) (models.Reward, error) {
	return s.repo.GetById(rewardId)
}

func (s *AdminRewardService) GetAllRewards() ([]models.Reward, error) {
	return s.repo.GetAllRewards()
}

func (s *AdminRewardService) Delete(rewardId int) error {
	return s.repo.Delete(rewardId)
}

func (s *AdminRewardService) UpdateReward(rewardId int, input models.UpdateRewardInput) error {
	const op = "service.admin_reward_service.UpdateReward"

	if err := input.Validate(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return s.repo.UpdateReward(rewardId, input)
}
