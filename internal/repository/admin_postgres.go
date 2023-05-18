package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type AdminPostgres struct {
	dbpool *pgxpool.Pool
}

func NewAdminPostgres(dbpool *pgxpool.Pool) Admin {
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
