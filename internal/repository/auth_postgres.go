package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
)

type AuthPostgres struct {
	db *pgx.Conn
}

func NewAuthPostgres(db *pgx.Conn) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
// 	var id int
// 	query := fmt.Sprintf("INSERT INTO %s (user_name, first_name, last_name, email, password_hash) values ($1, $2, $3, $4, $5) RETURNING id", usersTable)

// 	row := r.db.QueryRow(query, user.Username, user.FirstName, user.LastName, user.Email, user.Password)
// 	if err := row.Scan(&id); err != nil {
// 		return 0, err
// 	}

// 	return id, nil
// }

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := "INSERT INTO user_account (user_name, first_name, last_name, email, password_hash) values ($1, $2, $3, $4, $5) RETURNING id"

	row := r.db.QueryRow(context.Background(), query, user.Username, user.FirstName, user.LastName, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User
	query := "SELECT id FROM user_account WHERE user_name=$1 AND password_hash=$2"

	err := r.db.QueryRow(context.Background(), query, username, password).Scan(&user.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return user, err
	}

	return user, err
}
