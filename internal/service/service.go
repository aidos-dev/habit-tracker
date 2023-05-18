package service

import (
	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/aidos-dev/habit-tracker/internal/repository"
)

type Admin interface {
	AssignRole(userId int, role string) (int, error)
}

type Authorization interface {
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (*tokenClaims, error)
}

type User interface {
	CreateUser(user models.User) (int, error)
}

type Habit interface {
	Create(userId int, habit models.Habit) (int, error)
	GetAll(userId int) ([]models.Habit, error)
	GetById(userId, habitId int) (models.Habit, error)
	Delete(userId, habitId int) error
	Update(userId, habitId int, input models.UpdateHabitInput) error
}

type HabitTracker interface {
	// Create(userHabitId int, tracker habit.HabitTracker) (int, error) // temporarily disabled
	GetAll(userId int) ([]models.HabitTracker, error)
	GetById(userId, habitId int) (models.HabitTracker, error)
	// Delete(userId, habitId int) error // temporarily disabled
	Update(userId, habitId int, input models.UpdateTrackerInput) error
}

type Reward interface {
	Create(reward models.Reward) (int, error)
	GetById(rewardId int) (models.Reward, error)
	GetAllRewards() ([]models.Reward, error)
	Delete(rewardId int) error
	UpdateReward(rewardId int, input models.UpdateRewardInput) error
}

type UserReward interface {
	AssignReward(userId int, rewardId int, habitId int) (int, error)
	GetByUserId(userId int) ([]models.Reward, error)
	RemoveFromUser(userId, rewardId int) error
	UpdateUserReward(userId, rewardId int, input models.UpdateUserRewardInput) error
}

type Service struct {
	Admin
	Authorization
	User
	Habit
	HabitTracker
	Reward
	UserReward
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Admin:         NewAdminService(repos.Admin),
		Authorization: NewAuthService(repos.User),
		User:          NewUserService(repos.User),
		Habit:         NewHabitService(repos.Habit),
		HabitTracker:  NewHabitTrackerService(repos.HabitTracker),
		Reward:        NewRewardService(repos.Reward),
		UserReward:    NewUserRewardPostgres(repos.UserReward),
	}
}
