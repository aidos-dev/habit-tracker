package postgres

import (
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminPostgres struct {
	dbpool *pgxpool.Pool
	// repository.AdminUser
	repository.AdminRole
	repository.AdminReward
	repository.AdminUserReward
}

func NewAdminPostgres(dbpool *pgxpool.Pool) repository.Admin {
	return &AdminPostgres{dbpool: dbpool}
}
