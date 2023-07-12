package postgres

import (
	"context"
	"fmt"

	"github.com/aidos-dev/habit-tracker/backend/internal/config"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	nonUniqueValueCode = "23505"
	queryErr           = "queryRow failed"
	collectErr         = "collectRow failed"
	scanErr            = "row scan failed"
)

const (
	habitTable     = "habit-table"
	trackerTable   = "habit-tracker-table"
	userHabitTable = "user-habit-table"
)

func NewPostgresDB(cfg *config.Config) (*pgxpool.Pool, error) {
	const op = "repository.postgres.NewPostgresDB"

	dbpool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dbpool, err
}

func NewPostgresRepository(dbpool *pgxpool.Pool) *repository.Repository {
	return &repository.Repository{
		AdminRole:       NewAdminRolePostgres(dbpool),
		AdminReward:     NewAdminRewardPostgres(dbpool),
		AdminUserReward: NewAdminUserRewardPostgres(dbpool),
		Admin:           NewAdminPostgres(dbpool),
		User:            NewUserPostgres(dbpool),
		Habit:           NewHabitPostgres(dbpool),
		HabitTracker:    NewHabitTrackerPostgres(dbpool),
		Reward:          NewRewardPostgres(dbpool),
	}
}
