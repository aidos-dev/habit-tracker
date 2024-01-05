package postgres

import (
	"context"
	"fmt"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RewardPostgres struct {
	dbpool *pgxpool.Pool
}

func NewRewardPostgres(dbpool *pgxpool.Pool) repository.Reward {
	return &RewardPostgres{dbpool: dbpool}
}

func (r *RewardPostgres) GetPersonalRewardsByHabitId(userId, habitId int) ([]models.Reward, error) {
	const op = "repository.postgres.reward_postgres.GetPersonalRewardsByHabitId"

	var rewards []models.Reward

	query := `SELECT 
					tl.id, tl.title, tl.description 
				FROM reward tl INNER JOIN user_reward ul on tl.id = ul.reward_id
				WHERE ul.user_id = $1 AND ul.habit_id = $2`

	rowReward, err := r.dbpool.Query(context.Background(), query, userId, habitId)
	if err != nil {
		return rewards, fmt.Errorf("%s:%s: %w", op, queryErr, err)
	}

	defer rowReward.Close()

	rewards, err = pgx.CollectRows(rowReward, pgx.RowToStructByName[models.Reward])
	if err != nil {
		return rewards, fmt.Errorf("%s:%s: %w", op, collectErr, err)
	}

	return rewards, err
}

func (r *RewardPostgres) GetAllPersonalRewards(userId int) ([]models.Reward, error) {
	const op = "repository.postgres.reward_postgres.GetAllPersonalRewards"

	var rewards []models.Reward

	query := `SELECT 
					tl.id, tl.title, tl.description 
				FROM reward tl INNER JOIN user_reward ul on tl.id = ul.reward_id
				WHERE ul.user_id = $1`

	rowRewards, err := r.dbpool.Query(context.Background(), query, userId)
	if err != nil {
		return rewards, fmt.Errorf("%s:%s: %w", op, queryErr, err)
	}

	defer rowRewards.Close()

	rewards, err = pgx.CollectRows(rowRewards, pgx.RowToStructByName[models.Reward])
	if err != nil {
		return rewards, fmt.Errorf("%s:%s: %w", op, collectErr, err)
	}

	return rewards, err
}
