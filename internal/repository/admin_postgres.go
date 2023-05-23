package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminPostgres struct {
	dbpool *pgxpool.Pool
	AdminRole
	AdminReward
}

func NewAdminPostgres(dbpool *pgxpool.Pool) Admin {
	return &AdminPostgres{dbpool: dbpool}
}
