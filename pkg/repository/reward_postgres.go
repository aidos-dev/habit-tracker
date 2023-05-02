package repository

import (
	"fmt"
	"strings"

	"github.com/aidos-dev/habit-tracker"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type RewardPostgres struct {
	db *sqlx.DB
}

func NewRewardPostgres(db *sqlx.DB) *RewardPostgres {
	return &RewardPostgres{db: db}
}

func (r *RewardPostgres) Create(reward habit.Reward) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createRewardQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", rewardTable)
	row := tx.QueryRow(createRewardQuery, reward.Title, reward.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *RewardPostgres) AssignReward(userId int, rewardId int, habitId int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int

	assignRewardQuery := fmt.Sprintf("INSERT INTO %s (user_id, reward_id, habit_id) VALUES ($1, $2, $3)", userRewardTable)
	row := tx.QueryRow(assignRewardQuery, userId, rewardId, habitId)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *RewardPostgres) GetAllRewards() ([]habit.Reward, error) {
	var rewards []habit.Reward
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s",
		rewardTable)
	err := r.db.Select(&rewards, query)

	return rewards, err
}

func (r *RewardPostgres) GetById(rewardId int) (habit.Reward, error) {
	var reward habit.Reward

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s WHERE tl.id = $1",
		rewardTable)
	err := r.db.Get(&reward, query, rewardId)

	return reward, err
}

func (r *RewardPostgres) GetByUserId(userId int) ([]habit.Reward, error) {
	var rewards []habit.Reward
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description, ul.habit_id FROM %s tl INNER JOIN %s ul on tl.id = ul.reward_id WHERE ul.user_id = $1",
		rewardTable, userRewardTable)
	err := r.db.Select(&rewards, query, userId)

	return rewards, err
}

func (r *RewardPostgres) Delete(rewardId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.habit_id AND ul.user_id=$1 AND ul.habit_id=$2",
		habitsTable, usersHabitsTable)
	_, err := r.db.Exec(query, userId, habitId)

	return err
}

func (r *RewardPostgres) RemoveFromUser(userId, rewardId int) error {
	// Take away from user
}

// PENDING
func (r *RewardPostgres) UpdateReward(rewardId int) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.habit_id AND ul.habit_id=$%d AND ul.user_id=$%d",
		habitsTable, setQuery, usersHabitsTable, argId, argId+1)

	args = append(args, habitId, userId)

	logrus.Debugf("updateQuerry: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}

// PENDING
func (r *RewardPostgres) UpdateUsersReward(userId, rewardId int) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.habit_id AND ul.habit_id=$%d AND ul.user_id=$%d",
		habitsTable, setQuery, usersHabitsTable, argId, argId+1)

	args = append(args, habitId, userId)

	logrus.Debugf("updateQuerry: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
