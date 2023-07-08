package postgres

import (
	"context"
	"fmt"

	"github.com/aidos-dev/habit-tracker/backend/internal/config"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(cfg *config.Config) (*pgxpool.Pool, error) {
	// db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
	// 	cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	// if err != nil {
	// 	return nil, err
	// }

	dbpool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName))
	if err != nil {
		return nil, err
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return dbpool, nil
}

func NewPostgresRepository(dbpool *pgxpool.Pool) *repository.Repository {
	return &repository.Repository{
		AdminUser:       NewAdminUserPostgres(dbpool),
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
