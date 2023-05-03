package service

import (
	"github.com/aidos-dev/habit-tracker"
	"github.com/aidos-dev/habit-tracker/pkg/repository"
)

type RewardService struct {
	repo repository.Reward
}

func NewRewardService(repo repository.Reward) *RewardService {
	return &RewardService{repo: repo}
}

func (s *RewardService) Create(reward habit.Reward) (int, error) {
	return s.repo.Create(reward)
}

func (s *RewardService) AssignReward(userId int, rewardId int, habitId int) (int, error) {
	return s.repo.AssignReward(userId, rewardId, habitId)
}

func (s *RewardService) GetAllRewards() ([]habit.Reward, error) {
	return s.repo.GetAllRewards()
}

func (s *RewardService) GetById(rewardId int) (habit.Reward, error) {
	return s.repo.GetById(rewardId)
}

func (s *RewardService) GetByUserId(userId int) ([]habit.Reward, error) {
	return s.repo.GetByUserId(userId)
}

func (s *RewardService) Delete(rewardId int) error {
	return s.repo.Delete(rewardId)
}

func (s *RewardService) RemoveFromUser(userId, rewardId int) error {
	return s.repo.RemoveFromUser(userId, rewardId)
}

func (s *RewardService) UpdateReward(rewardId int, input habit.UpdateRewardInput) error {
	return s.repo.UpdateReward(rewardId, input)
}

func (s *RewardService) UpdateUsersReward(userId, rewardId int, input habit.UpdateUserRewardInput) error {
	return s.repo.UpdateUsersReward(userId, rewardId, input)
}
