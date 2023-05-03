package repository

import (
	"github.com/aidos-dev/habit-tracker"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user habit.User) (int, error)
	GetUser(username, password string) (habit.User, error)
}

type Habit interface {
	Create(userId int, habit habit.Habit) (int, error)
	GetAll(userId int) ([]habit.Habit, error)
	GetById(userId, habitId int) (habit.Habit, error)
	Delete(userId, habitId int) error
	Update(userId, habitId int, input habit.UpdateHabitInput) error
}

type HabitTracker interface {
	Create(habitId int, tracker habit.HabitTracker) (int, error)
	GetAll(userId int) ([]habit.HabitTracker, error)
	GetById(userId, habitId int) (habit.HabitTracker, error)
	Delete(userId, habitId int) error
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
	UpdateReward(rewardId int) error
	UpdateUsersReward(userId, rewardId int) error
}

type Repository struct {
	Authorization
	Habit
	HabitTracker
	Reward
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Habit:         NewHabitPostgres(db),
		HabitTracker:  NewHabitTrackerPostgres(db),
		Reward:        NewRewardPostgres(db),
	}
}
