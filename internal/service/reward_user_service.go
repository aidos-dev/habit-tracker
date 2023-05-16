package service

import (
	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/aidos-dev/habit-tracker/internal/repository"
)

type UserRewardService struct {
	repo repository.UserReward
}

func NewUserRewardPostgres(repo repository.UserReward) UserReward {
	return &UserRewardService{repo: repo}
}

func (s *UserRewardService) AssignReward(userId int, rewardId int, habitId int) (int, error) {
	return s.repo.AssignReward(userId, rewardId, habitId)
}

func (s *UserRewardService) GetByUserId(userId int) ([]models.Reward, error) {
	return s.repo.GetByUserId(userId)
}

func (s *UserRewardService) RemoveFromUser(userId, rewardId int) error {
	return s.repo.RemoveFromUser(userId, rewardId)
}

func (s *UserRewardService) UpdateUserReward(userId, rewardId int, input models.UpdateUserRewardInput) error {
	return s.repo.UpdateUserReward(userId, rewardId, input)
}
