package repository

import (
	"github.com/aidos-dev/habit-tracker/backend/internal/models"

	_ "github.com/jackc/pgx/v5"
)

type AdminUser interface {
	GetAllUsers() ([]models.GetUser, error)
	GetUserById(userId int) (models.GetUser, error)
	User
}

type AdminRole interface {
	AssignRole(userId int, role models.UpdateRoleInput) (int, error)
}

type AdminReward interface {
	Create(reward models.Reward) (int, error)
	GetById(rewardId int) (models.Reward, error)
	GetAllRewards() ([]models.Reward, error)
	Delete(rewardId int) error
	UpdateReward(rewardId int, input models.UpdateRewardInput) error
}

type AdminUserReward interface {
	AssignReward(userId, habitId, rewardId int) (int, error)
	RemoveFromUser(userId, habitId, rewardId int) error
	UpdateUserReward(userId, habitId, rewardId int, input models.UpdateUserRewardInput) error
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
	DeleteUser(userId int) (int, error)
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
	GetPersonalRewardsByHabitId(userId, habitId int) ([]models.Reward, error)
	GetAllPersonalRewards(userId int) ([]models.Reward, error)
}

type Repository struct {
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
