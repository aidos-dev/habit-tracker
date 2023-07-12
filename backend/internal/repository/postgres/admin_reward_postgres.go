package postgres

import (
	"context"
	"fmt"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRewardPostgres struct {
	dbpool *pgxpool.Pool
}

func NewAdminRewardPostgres(dbpool *pgxpool.Pool) repository.AdminReward {
	return &AdminRewardPostgres{dbpool: dbpool}
}

func (r *AdminRewardPostgres) Create(reward models.Reward) (int, error) {
	const op = "repository.postgres.admin_reward_postgres.Create"

	var rewardId int
	query := `INSERT INTO 
						reward (title, description) 
						VALUES ($1, $2) 
					RETURNING id`

	row := r.dbpool.QueryRow(context.Background(), query, reward.Title, reward.Description)
	if err := row.Scan(&rewardId); err != nil {
		return 0, fmt.Errorf("%s:%s: %w", op, scanErr, err)
	}

	return rewardId, nil
}

func (r *AdminRewardPostgres) GetById(rewardId int) (models.Reward, error) {
	const op = "repository.postgres.admin_reward_postgres.GetById"

	var reward models.Reward

	query := `SELECT 
					id, 
					title, 
					description 
				FROM 
					reward
				WHERE id = $1`

	rowReward, err := r.dbpool.Query(context.Background(), query, rewardId)
	if err != nil {
		return reward, fmt.Errorf("%s:%s: %w", op, queryErr, err)
	}

	defer rowReward.Close()

	reward, err = pgx.CollectOneRow(rowReward, pgx.RowToStructByName[models.Reward])
	if err != nil {
		return reward, fmt.Errorf("%s:%s: %w", op, collectErr, err)
	}

	return reward, nil
}

func (r *AdminRewardPostgres) GetAllRewards() ([]models.Reward, error) {
	const op = "repository.postgres.admin_reward_postgres.GetAllRewards"

	var rewards []models.Reward

	query := `SELECT 
					id, 
					title, 
					description 
				FROM 
					reward`

	rowsRewards, err := r.dbpool.Query(context.Background(), query)
	if err != nil {
		return rewards, fmt.Errorf("%s:%s: %w", op, queryErr, err)
	}

	defer rowsRewards.Close()

	rewards, err = pgx.CollectRows(rowsRewards, pgx.RowToStructByName[models.Reward])
	if err != nil {
		return rewards, fmt.Errorf("%s:%s: %w", op, collectErr, err)
	}

	return rewards, nil
}

func (r *AdminRewardPostgres) Delete(rewardId int) error {
	const op = "repository.postgres.admin_reward_postgres.Delete"

	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	queryReward := `DELETE FROM 
							reward 
						WHERE id = $1
						RETURNING id`

	var checkrewardId int

	rowReward := tx.QueryRow(context.Background(), queryReward, rewardId)
	err = rowReward.Scan(&checkrewardId)
	if err != nil {
		tx.Rollback(context.Background())
		return fmt.Errorf("%s:%s: %w", op, scanErr, err)
	}

	return tx.Commit(context.Background())
}

func (r *AdminRewardPostgres) UpdateReward(rewardId int, input models.UpdateRewardInput) error {
	const op = "repository.postgres.admin_reward_postgres.UpdateReward"

	query := `UPDATE 
					reward 
				SET 
					title=COALESCE($2, title), 
					description=COALESCE($3, description)
				WHERE id = $1
				RETURNING id`

	var checkRewardId int

	rowRewardd := r.dbpool.QueryRow(context.Background(), query, rewardId, input.Title, input.Description)
	err := rowRewardd.Scan(&checkRewardId)
	if err != nil {
		return fmt.Errorf("%s:%s: %w", op, scanErr, err)
	}

	return nil
}
