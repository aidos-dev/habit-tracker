package repository

import (
	"context"
	"fmt"

	"github.com/aidos-dev/habit-tracker/internal/models"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RewardPostgres struct {
	dbpool *pgxpool.Pool
}

func NewRewardPostgres(dbpool *pgxpool.Pool) Reward {
	return &RewardPostgres{dbpool: dbpool}
}

func (r *RewardPostgres) GetRewardById(rewardId int) (models.Reward, error) {
	var reward models.Reward

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s WHERE tl.id = $1",
		rewardTable)
	err := r.dbpool.QueryRow(context.Background(), query, rewardId).Scan(&reward)

	return reward, err
}

func (r *RewardPostgres) GetAllPersonalRewards(userId int) ([]models.Reward, error) {
	var rewards []models.Reward
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s",
		rewardTable)
	err := r.dbpool.QueryRow(context.Background(), query).Scan(&rewards)

	return rewards, err
}
