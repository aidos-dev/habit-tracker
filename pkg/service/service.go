package service

import (
	"github.com/aidos-dev/habit-tracker"
	"github.com/aidos-dev/habit-tracker/pkg/repository"
)

type Authorization interface {
	CreateUser(user habit.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Habit interface {
	Create(userId int, habit habit.Habit) (int, error)
	GetAll(userId int) ([]habit.Habit, error)
	GetById(userId, habitId int) (habit.Habit, error)
	Delete(userId, habitId int) error
	Update(userId, habitId int, input habit.UpdateHabitInput) error
}

type HabitTracker interface {
	// Create(userHabitId int, tracker habit.HabitTracker) (int, error) // temporarily disabled
	GetAll(userId int) ([]habit.HabitTracker, error)
	GetById(userId, habitId int) (habit.HabitTracker, error)
	// Delete(userId, habitId int) error // temporarily disabled
	Update(userId, habitId int, input habit.UpdateTrackerInput) error
}

type Reward interface {
	Create(reward habit.Reward) (int, error)
	AssignReward(userId int, rewardId int, habitId int) (int, error)
	GetAllRewards() ([]habit.Reward, error)
	GetById(rewardId int) (habit.Reward, error)
	GetByUserId(userId int) ([]habit.Reward, error)
	Delete(rewardId int) error
	RemoveFromUser(userId, rewardId int) error
	UpdateReward(rewardId int, input habit.UpdateRewardInput) error
	UpdateUserReward(userId, rewardId int, input habit.UpdateUserRewardInput) error
}
type Service struct {
	Authorization
	Habit
	HabitTracker
	Reward
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Habit:         NewHabitService(repos.Habit),
		HabitTracker:  NewHabitTrackerService(repos.HabitTracker),
		Reward:        NewRewardService(repos.Reward),
	}
}
