package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type AdminRolePostgres struct {
	dbpool *pgxpool.Pool
}

func NewAdminRolePostgres(dbpool *pgxpool.Pool) AdminRole {
	return &AdminPostgres{dbpool: dbpool}
}

func (r *AdminPostgres) AssignRole(userId int, role string) (int, error) {
	var id int

	query := `UPDATE 
					user_account tl 
				SET 
					role=$2
				WHERE tl.id =$1
				RETURNING id`

	logrus.Debugf("assignRoleQuerry: %s", query)

	row := r.dbpool.QueryRow(context.Background(), query, userId, role)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
