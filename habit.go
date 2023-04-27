package habit

import (
	"errors"
	"time"
)

type Habit struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UsersHabits struct {
	Id      int
	UserId  int
	HabitId int
}

type HabitTracker struct {
	Id            int       `json:"id" db:"id"`
	UserHabitId   int       `json:"user_habit_id" db:"user_habit_id" binding:"required"`
	UnitOfMessure string    `json:"unit_of_messure" db:"unit_of_messure" binding:"required"`
	Goal          string    `json:"goal" db:"goal" binding:"required"`
	Frequency     string    `json:"frequency" db:"frequency" binding:"required"`
	StartDate     time.Time `json:"start_date" db:"start_date"`
	EndDate       time.Time `json:"end_date" db:"end_date"`
	Counter       int       `json:"counter" db:"counter"`
	Done          bool      `json:"done" db:"done"`
}

type Reward struct {
	Id             int    `json:"id" db:"id"`
	HabitTrackerId int    `json:"habit_tracker_id" db:"habit_tracker_id" binding:"required"`
	Title          string `json:"title" db:"title" binding:"required"`
	Description    string `json:"description" db:"description"`
}

type UpdateHabitInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdateHabitInput) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

type UpdateTrackerInput struct {
	UnitOfMessure *string    `json:"unit_of_messure"`
	Goal          *string    `json:"goal"`
	Frequency     *string    `json:"frequency"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	Counter       *int       `json:"counter"`
	Done          *bool      `json:"done"`
}

func (i UpdateTrackerInput) Validate() error {
	if i.UnitOfMessure == nil && i.Goal == nil && i.Frequency == nil && i.StartDate == nil && i.EndDate == nil && i.Counter == nil && i.Done == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
