package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPostgres struct {
	dbpool *pgxpool.Pool
}

func NewUserPostgres(dbpool *pgxpool.Pool) repository.User {
	return &UserPostgres{dbpool: dbpool}
}

func (r *UserPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := `INSERT INTO 
						user_account (user_name, tg_user_name, first_name, last_name, email, password_hash) 
						VALUES (
							COALESCE(NULLIF($1, ''), NULL),
							COALESCE(NULLIF($2, ''), NULL),
							COALESCE(NULLIF($3, ''), NULL),
							COALESCE(NULLIF($4, ''), NULL),
							COALESCE(NULLIF($5, ''), NULL),
							COALESCE(NULLIF($6, ''), NULL)
							) 
					RETURNING id`

	row := r.dbpool.QueryRow(context.Background(), query, user.Username, user.TgUsername, user.FirstName, user.LastName, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User
	query := `SELECT 
					id,
					user_name, 
					first_name, 
					last_name, 
					email,
					password_hash, 
					role 
				FROM 
					user_account 
				WHERE user_name=$1 AND password_hash=$2`

	userHabit, err := r.dbpool.Query(context.Background(), query, username, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error from GetUser: QueryRow failed: %v\n", err)
		return user, err
	}

	user, err = pgx.CollectOneRow(userHabit, pgx.RowToStructByName[models.User])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error from GetUser: Collect One Row failed: %v\n", err)
		return user, err
	}

	return user, err
}

func (r *UserPostgres) DeleteUser(userId int) (int, error) {
	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return 0, err
	}

	var checkUserId int

	query := `DELETE FROM
						user_account
					WHERE id=$1
					RETURNING id`

	rowUser := tx.QueryRow(context.Background(), query, userId)
	err = rowUser.Scan(&checkUserId)
	if err != nil {
		tx.Rollback(context.Background())
		fmt.Printf("err: repository: user_postgres.go: DeleteUser: rowUser.Scan: user doesn't exist: %v\n", err)
		return 0, err
	}

	return checkUserId, tx.Commit(context.Background())
}
