package repository

import (
	"context"
	"fmt"

	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
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

	var userRewardId int

	assignRewardQuery := `INSERT INTO 
								user_reward (user_id, reward_id, habit_id) 
							VALUES ($1, $2, $3)
							RETURNING id`

	rowUserReward := tx.QueryRow(context.Background(), assignRewardQuery, userId, rewardId, habitId)
	if err := rowUserReward.Scan(&userRewardId); err != nil {
		tx.Rollback(context.Background())
		return 0, err
	}

	return userRewardId, tx.Commit(context.Background())
}

// Take away from user
func (r *AdminUserRewardPostgres) RemoveFromUser(userId, rewardId int) error {
	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return err
	}

	var checkUserRewardId int

	queryUserReward := `DELETE FROM 
							user_reward 
						WHERE user_id = $1 AND reward_id=$2
						RETURNING id`

	rowUserReward := tx.QueryRow(context.Background(), queryUserReward, userId, rewardId)

	err = rowUserReward.Scan(&checkUserRewardId)
	if err != nil {
		tx.Rollback(context.Background())
		fmt.Printf("err: repository: admin_user_reward_postgres.go: RemoveFromUser: rowUserReward.Scan: user_reward doesn't exist: %v\n", err)
		return err
	}

	return tx.Commit(context.Background())
}

func (r *AdminUserRewardPostgres) UpdateUserReward(userId, rewardId int, input models.UpdateUserRewardInput) error {
	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return err
	}

	var userRewardId int

	query := `UPDATE 
					user_reward
				SET 
					reward_id = COALESCE($3, reward_id), 
					habit_id = COALESCE($4, habit_id)
				WHERE user_id = $1 AND reward_id = $2
				RETURNING id`

	rowUserReward := r.dbpool.QueryRow(context.Background(), query, userId, rewardId, input.RewardId, input.HabitId)
	err = rowUserReward.Scan(&userRewardId)

	if err != nil {
		tx.Rollback(context.Background())
		fmt.Printf("err: repository: admin_user_reward_postgres.go: UpdateUserReward: rowUserReward.Scan: user_reward doesn't exist: %v\n", err)
		return err
	}

	return tx.Commit(context.Background())
}
