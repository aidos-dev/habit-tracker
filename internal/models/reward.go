package models

import "errors"

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
