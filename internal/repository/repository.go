package repository

import (
	"github.com/aidos-dev/habit-tracker/internal/models"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Admin interface {
	AssignRole(userId int, role string) (int, error)
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

type Repository struct {
	Admin
	User
	Habit
	HabitTracker
	Reward
	UserReward
}

func NewRepository(dbpool *pgxpool.Pool) *Repository {
	return &Repository{
		Admin:        NewAdminPostgres(dbpool),
		User:         NewUserPostgres(dbpool),
		Habit:        NewHabitPostgres(dbpool),
		HabitTracker: NewHabitTrackerPostgres(dbpool),
		Reward:       NewRewardPostgres(dbpool),
		UserReward:   NewUserRewardPostgres(dbpool),
	}
}
