package service

import (
	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/aidos-dev/habit-tracker/internal/repository"
	"github.com/golang-jwt/jwt"
)

type Authorization interface {
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (jwt.MapClaims, error)
}

type AdminUser interface {
	GetAllUsers() ([]models.GetUser, error)
	User
}

type AdminRole interface {
	AssignRole(userId int, role string) (int, error)
}

type AdminReward interface {
	Create(reward models.Reward) (int, error)
	GetById(rewardId int) (models.Reward, error)
	GetAllRewards() ([]models.Reward, error)
	Delete(rewardId int) error
	UpdateReward(rewardId int, input models.UpdateRewardInput) error
}

type AdminUserReward interface {
	AssignReward(userId, rewardId, habitId int) (int, error)
	RemoveFromUser(userId, rewardId int) error
	UpdateUserReward(userId, rewardId int, input models.UpdateUserRewardInput) error
	Reward
}

type Admin interface {
	AdminUser
	AdminRole
	AdminReward
	AdminUserReward
}

type User interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
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
	GetPersonalRewardById(userId, rewardId int) (models.Reward, error)
	GetAllPersonalRewards(userId int) ([]models.Reward, error)
}

type Service struct {
	Authorization
	AdminUser
	AdminRole
	AdminReward
	AdminUserReward
	Admin
	User
	Habit
	HabitTracker
	Reward
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization:   NewAuthService(repos.User),
		AdminUser:       NewAdminUserService(repos.AdminUser),
		AdminRole:       NewAdminRoleService(repos.AdminRole),
		AdminReward:     NewAdminRewardService(repos.AdminReward),
		AdminUserReward: NewAdminUserRewardService(repos.AdminUserReward),
		Admin:           NewAdminService(repos.Admin),
		User:            NewUserService(repos.User),
		Habit:           NewHabitService(repos.Habit),
		HabitTracker:    NewHabitTrackerService(repos.HabitTracker),
		Reward:          NewRewardService(repos.Reward),
	}
}
