package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type AdminUserRewardPostgres struct {
	dbpool *pgxpool.Pool
	Reward
}

func NewAdminUserRewardPostgres(dbpool *pgxpool.Pool) AdminUserReward {
	return &AdminUserRewardPostgres{dbpool: dbpool}
}

func (r *AdminUserRewardPostgres) AssignReward(userId, rewardId, habitId int) (int, error) {
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

// Take away from user
func (r *AdminUserRewardPostgres) RemoveFromUser(userId, rewardId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl WHERE tl.user_id = $1 AND tl.reward_id=$2",
		userRewardTable)
	_, err := r.dbpool.Exec(context.Background(), query, userId, rewardId)

	return err
}

func (r *AdminUserRewardPostgres) UpdateUserReward(userId, rewardId int, input models.UpdateUserRewardInput) error {
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
