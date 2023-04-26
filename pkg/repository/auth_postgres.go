package repository

import (
	"fmt"

	"github.com/aidos-dev/habit-tracker"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user habit.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, firstName, lastName, eMail, password_hash) values ($1, $2, $3, $4, $5) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Username, user.FirstName, user.LastName, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (habit.User, error) {
	var user habit.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}
