package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/aidos-dev/habit-tracker/internal/models"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type RewardPostgres struct {
	dbpool *pgxpool.Pool
}

func NewRewardPostgres(dbpool *pgxpool.Pool) *RewardPostgres {
	return &RewardPostgres{dbpool: dbpool}
}

func (r *RewardPostgres) Create(reward models.Reward) (int, error) {
	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return 0, err
	}

	var id int
	createRewardQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", rewardTable)
	row := tx.QueryRow(context.Background(), createRewardQuery, reward.Title, reward.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback(context.Background())
		return 0, err
	}

	return id, tx.Commit(context.Background())
}

func (r *RewardPostgres) AssignReward(userId int, rewardId int, habitId int) (int, error) {
	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return 0, err
	}

	var id int

	assignRewardQuery := fmt.Sprintf("INSERT INTO %s (user_id, reward_id, habit_id) VALUES ($1, $2, $3)", userRewardTable)
	row := tx.QueryRow(context.Background(), assignRewardQuery, userId, rewardId, habitId)
	if err := row.Scan(&id); err != nil {
		tx.Rollback(context.Background())
		return 0, err
	}

	return id, tx.Commit(context.Background())
}

func (r *RewardPostgres) GetAllRewards() ([]models.Reward, error) {
	var rewards []models.Reward
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s",
		rewardTable)
	err := r.dbpool.QueryRow(context.Background(), query).Scan(&rewards)

	return rewards, err
}

func (r *RewardPostgres) GetById(rewardId int) (models.Reward, error) {
	var reward models.Reward

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s WHERE tl.id = $1",
		rewardTable)
	err := r.dbpool.QueryRow(context.Background(), query, rewardId).Scan(&reward)

	return reward, err
}

func (r *RewardPostgres) GetByUserId(userId int) ([]models.Reward, error) {
	var rewards []models.Reward
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description, ul.habit_id FROM %s tl INNER JOIN %s ul on tl.id = ul.reward_id WHERE ul.user_id = $1",
		rewardTable, userRewardTable)
	err := r.dbpool.QueryRow(context.Background(), query, userId).Scan(&rewards)

	return rewards, err
}

func (r *RewardPostgres) Delete(rewardId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl WHERE tl.id = $1",
		rewardTable)
	_, err := r.dbpool.Exec(context.Background(), query, rewardId)

	return err
}

// Take away from user
func (r *RewardPostgres) RemoveFromUser(userId, rewardId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl WHERE tl.user_id = $1 AND tl.reward_id=$2",
		userRewardTable)
	_, err := r.dbpool.Exec(context.Background(), query, userId, rewardId)

	return err
}

func (r *RewardPostgres) UpdateReward(rewardId int, input models.UpdateRewardInput) error {
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

	query := fmt.Sprintf("UPDATE %s tl SET %s WHERE tl.id = %d",
		rewardTable, setQuery, rewardId)

	logrus.Debugf("updateQuerry: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.dbpool.Exec(context.Background(), query, args...)
	return err
}

func (r *RewardPostgres) UpdateUserReward(userId, rewardId int, input models.UpdateUserRewardInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.RewardId != nil {
		setValues = append(setValues, fmt.Sprintf("reward_id=$%d", argId))
		args = append(args, *input.RewardId)
		argId++
	}

	if input.HabitId != nil {
		setValues = append(setValues, fmt.Sprintf("habit_id=$%d", argId))
		args = append(args, *input.HabitId)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s WHERE tl.user_id = $%d AND tl.reward_id=$%d",
		userRewardTable, setQuery, userId, rewardId)

	logrus.Debugf("updateQuerry: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.dbpool.Exec(context.Background(), query, args...)
	return err
}
