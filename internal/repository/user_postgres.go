package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/aidos-dev/habit-tracker/internal/models"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPostgres struct {
	dbpool *pgxpool.Pool
}

func NewUserPostgres(dbpool *pgxpool.Pool) User {
	return &UserPostgres{dbpool: dbpool}
}

func (r *UserPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := `INSERT INTO 
						user_account (user_name, first_name, last_name, email, password_hash) 
						VALUES ($1, $2, $3, $4, $5) 
					RETURNING id`

	row := r.dbpool.QueryRow(context.Background(), query, user.Username, user.FirstName, user.LastName, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User
	query := `SELECT 
					id 
				FROM 
					user_account 
				WHERE user_name=$1 AND password_hash=$2`

	err := r.dbpool.QueryRow(context.Background(), query, username, password).Scan(&user.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return user, err
	}

	return user, err
}
