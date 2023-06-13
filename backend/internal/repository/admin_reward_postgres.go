package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type AdminRewardPostgres struct {
	dbpool *pgxpool.Pool
}

func NewAdminRewardPostgres(dbpool *pgxpool.Pool) AdminReward {
	return &AdminRewardPostgres{dbpool: dbpool}
}

func (r *AdminRewardPostgres) Create(reward models.Reward) (int, error) {
	var rewardId int
	query := `INSERT INTO 
						reward (title, description) 
						VALUES ($1, $2) 
					RETURNING id`

	row := r.dbpool.QueryRow(context.Background(), query, reward.Title, reward.Description)
	if err := row.Scan(&rewardId); err != nil {
		return 0, err
	}

	return rewardId, nil
}

func (r *AdminRewardPostgres) GetById(rewardId int) (models.Reward, error) {
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
		fmt.Fprintf(os.Stderr, "error from GetById: QueryRow failed: %v\n", err)
		return reward, err
	}

	defer rowReward.Close()

	reward, err = pgx.CollectOneRow(rowReward, pgx.RowToStructByName[models.Reward])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error from GetById: Collect One Row failed: %v\n", err)
		return reward, err
	}

	return reward, nil
}

func (r *AdminRewardPostgres) GetAllRewards() ([]models.Reward, error) {
	var rewards []models.Reward

	query := `SELECT 
					id, 
					title, 
					description 
				FROM 
					reward`

	rowsRewards, err := r.dbpool.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return rewards, err
	}

	defer rowsRewards.Close()

	rewards, err = pgx.CollectRows(rowsRewards, pgx.RowToStructByName[models.Reward])
	if err != nil {
		fmt.Fprintf(os.Stderr, "rowsRewards CollectRows failed: %v\n", err)
		return rewards, err
	}

	return rewards, nil
}

func (r *AdminRewardPostgres) Delete(rewardId int) error {
	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return err
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
		fmt.Printf("err: repository: admin_reward_postgres.go: Delete: rowReward.Scan: reward doesn't exist: %v\n", err)
		return err
	}

	return tx.Commit(context.Background())
}

func (r *AdminRewardPostgres) UpdateReward(rewardId int, input models.UpdateRewardInput) error {
	query := `UPDATE 
					reward 
				SET 
					title=COALESCE($2, title), 
					description=COALESCE($3, description)
				WHERE id = $1
				RETURNING id`

	logrus.Debugf("updateQuerry: %s", query)

	var checkRewardId int

	rowRewardd := r.dbpool.QueryRow(context.Background(), query, rewardId, input.Title, input.Description)
	err := rowRewardd.Scan(&checkRewardId)
	if err != nil {

		fmt.Printf("err: repository: admin_reward_postgres.go: UpdateReward: rowReward.Scan: reward doesn't exist: %v\n", err)
		return err
	}

	return err
}
