package repository

import (
	"context"

	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type AdminRolePostgres struct {
	dbpool *pgxpool.Pool
}

func NewAdminRolePostgres(dbpool *pgxpool.Pool) AdminRole {
	return &AdminRolePostgres{dbpool: dbpool}
}

func (r *AdminRolePostgres) AssignRole(userId int, role models.UpdateRoleInput) (int, error) {
	var id int

	query := `UPDATE 
					user_account
				SET 
					role=$2
				WHERE id =$1
				RETURNING id`

	logrus.Debugf("assignRoleQuerry: %s", query)

	row := r.dbpool.QueryRow(context.Background(), query, userId, role.Role)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
