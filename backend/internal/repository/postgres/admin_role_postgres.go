package postgres

import (
	"context"
	"fmt"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRolePostgres struct {
	dbpool *pgxpool.Pool
}

func NewAdminRolePostgres(dbpool *pgxpool.Pool) repository.AdminRole {
	return &AdminRolePostgres{dbpool: dbpool}
}

func (r *AdminRolePostgres) AssignRole(userId int, role models.UpdateRoleInput) (int, error) {
	const op = "repository.postgres.AssignRole"

	var id int

	query := `UPDATE 
					user_account
				SET 
					role=$2
				WHERE id =$1
				RETURNING id`

	row := r.dbpool.QueryRow(context.Background(), query, userId, role.Role)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s:%s: %w", op, scanErr, err)
	}

	return id, nil
}
