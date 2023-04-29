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
	Create(habitTrackerId int, reward habit.Reward) (int, error)
	GetAll(userId int) ([]habit.Reward, error)
	GetById(userId, habitId int) ([]habit.Reward, error)
	Delete(habitTrackerId, rewardId int) error
	Update(userId, habitId int, input habit.UpdateTrackerInput) error
}

type Repository struct {
	Authorization
	Habit
	HabitTracker
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Habit:         NewHabitPostgres(db),
		HabitTracker:  NewHabitTrackerPostgres(db),
	}
}
