package postgres

import (
	"context"
	"fmt"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminUserRewardPostgres struct {
	dbpool *pgxpool.Pool
	repository.Reward
}

func NewAdminUserRewardPostgres(dbpool *pgxpool.Pool) repository.AdminUserReward {
	return &AdminUserRewardPostgres{dbpool: dbpool}
}

func (r *AdminUserRewardPostgres) AssignReward(userId, habitId, rewardId int) (int, error) {
	const op = "repository.postgres.AssignReward"

	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var userRewardId int

	assignRewardQuery := `INSERT INTO
								user_reward (user_id, habit_id, reward_id)
							SELECT 
								ht.user_id, ht.habit_id, rt.id
							FROM user_habit AS ht, reward AS rt	
							WHERE
								ht.user_id = $1
							AND  
								ht.habit_id = $2
							AND  
								rt.id = $3

							RETURNING id`

	rowUserReward := tx.QueryRow(context.Background(), assignRewardQuery, userId, habitId, rewardId)
	if err := rowUserReward.Scan(&userRewardId); err != nil {
		tx.Rollback(context.Background())
		return 0, fmt.Errorf("%s:%s: %w", op, scanErr, err)
	}

	return userRewardId, tx.Commit(context.Background())
}

// Take away from user
func (r *AdminUserRewardPostgres) RemoveFromUser(userId, habitId, rewardId int) error {
	const op = "repository.postgres.RemoveFromUser"

	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var checkUserRewardId int

	queryUserReward := `DELETE FROM
								user_reward
							WHERE user_id = $1 AND habit_id = $2 AND reward_id=$3
							RETURNING id`

	rowUserReward := tx.QueryRow(context.Background(), queryUserReward, userId, habitId, rewardId)

	err = rowUserReward.Scan(&checkUserRewardId)
	if err != nil {
		tx.Rollback(context.Background())
		return fmt.Errorf("%s:%s: %w", op, scanErr, err)
	}

	// queryUserReward := `IF EXISTS (
	// 						SELECT 1
	// 						FROM
	// 							user_reward
	// 						WHERE
	// 							(user_id=$1) AND (reward_id=$2)
	// 					)
	// 					BEGIN
	// 					DELETE FROM
	// 						user_reward
	// 					WHERE
	// 						(user_id=$1) AND (reward_id=$2)
	// 					END
	// 					ELSE
	// 					BEGIN
	// 						RAISERROR ('No rows found', 16, 1)
	// 					END`

	// _, err = tx.Exec(context.Background(), queryUserReward, userId, rewardId)
	// if err != nil {
	// 	tx.Rollback(context.Background())
	// 	fmt.Printf("err: repository: admin_user_reward_postgres.go: RemoveFromUser: rowUserReward.Scan: user_reward doesn't exist: %v\n", err)
	// 	return err
	// }

	return tx.Commit(context.Background())
}

func (r *AdminUserRewardPostgres) UpdateUserReward(userId, habitId, rewardId int, input models.UpdateUserRewardInput) error {
	const op = "repository.postgres.UpdateUserReward"

	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var userRewardId int

	query := `UPDATE 
					user_reward
				SET 
					habit_id = COALESCE($4, habit_id),
					reward_id = COALESCE($5, reward_id) 
				WHERE user_id = $1 AND habit_id = $2 AND reward_id=$3
				RETURNING id`

	rowUserReward := r.dbpool.QueryRow(context.Background(), query, userId, habitId, rewardId, input.HabitId, input.RewardId)
	err = rowUserReward.Scan(&userRewardId)

	if err != nil {
		tx.Rollback(context.Background())
		return fmt.Errorf("%s:%s: %w", op, scanErr, err)
	}

	return tx.Commit(context.Background())
}
