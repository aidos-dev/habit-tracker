package habit

import (
	"errors"
	"time"
)

type Habit struct {
	Id          int    `json:"habitId" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UsersHabits struct {
	Id             int
	UserId         int
	HabitId        int
	HabitTrackerId int
}

type HabitTracker struct {
	Id            int       `json:"trackerId" db:"id"`
	HabitId       int       `json:"habitId" db:"habit_id"`
	UnitOfMessure string    `json:"unit_of_messure" db:"unit_of_messure" binding:"required"`
	Goal          string    `json:"goal" db:"goal" binding:"required"`
	Frequency     string    `json:"frequency" db:"frequency" binding:"required"`
	StartDate     time.Time `json:"start_date" db:"start_date"`
	EndDate       time.Time `json:"end_date" db:"end_date"`
	Counter       int       `json:"counter" db:"counter"`
	Done          bool      `json:"done" db:"done"`
}

type Reward struct {
	Id          int    `json:"rewardId" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	HabitId     int    `json:"habitId" db:"habitId"`
}

type UserReward struct {
	Id       int
	UserId   int
	RewardId int
	HabitId  int
}

type UpdateHabitInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdateHabitInput) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("habit update structure has no values")
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
		return errors.New("habit tracker update structure has no values")
	}

	return nil
}

type UpdateRewardInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdateRewardInput) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("reward update structure has no values")
	}
	return nil
}

type UpdateUserRewardInput struct {
	RewardId *string `json:"rewardId"`
	HabitId  *string `json:"habitId"`
}

func (i UpdateUserRewardInput) Validate() error {
	if i.RewardId == nil && i.HabitId == nil {
		return errors.New("user reward update structure has no values")
	}
	return nil
}
