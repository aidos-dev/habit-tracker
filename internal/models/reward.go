package models

import "errors"

type Reward struct {
	Id          int    `json:"rewardId" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UserReward struct {
	Id       int `json:"userRewardId" db:"id"`
	UserId   int `json:"title" db:"title" binding:"required"`
	HabitId  int `json:"habitId" db:"habitId"`
	RewardId int `json:"rewardId" db:"rewardId"`
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
	HabitId  *int `json:"habitId"`
	RewardId *int `json:"rewardId"`
}

func (i UpdateUserRewardInput) Validate() error {
	if i.RewardId == nil && i.HabitId == nil {
		return errors.New("user reward update structure has no values")
	}
	return nil
}
