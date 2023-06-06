package service

import (
	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/aidos-dev/habit-tracker/internal/repository"
)

type RewardService struct {
	repo repository.Reward
}

func NewRewardService(repo repository.Reward) Reward {
	return &RewardService{repo: repo}
}

func (r *RewardService) GetPersonalRewardsByHabitId(userId, habitId int) ([]models.Reward, error) {
	return r.repo.GetPersonalRewardsByHabitId(userId, habitId)
}

func (r *RewardService) GetAllPersonalRewards(userId int) ([]models.Reward, error) {
	return r.repo.GetAllPersonalRewards(userId)
}
