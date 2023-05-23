package repository

import (
	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRewardPostgres struct {
	dbpool *pgxpool.Pool
}

func NewAdminRewardPostgres(dbpool *pgxpool.Pool) AdminReward {
	return &AdminPostgres{dbpool: dbpool}
}

func (r *AdminRewardPostgres) Create(reward models.Reward) (int, error) {
	return 0, nil
}

func (r *AdminRewardPostgres) Delete(rewardId int) error {
	return nil
}

func (r *AdminRewardPostgres) UpdateReward(rewardId int, input models.UpdateRewardInput) error {
	return nil
}
