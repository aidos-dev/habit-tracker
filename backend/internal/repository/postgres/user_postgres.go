package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPostgres struct {
	dbpool *pgxpool.Pool
}

func NewUserPostgres(dbpool *pgxpool.Pool) repository.User {
	return &UserPostgres{dbpool: dbpool}
}

func (r *UserPostgres) CreateUser(user models.User) (int, error) {
	const (
		op         = "repository.postgres.CreateUser"
		userExists = "such user already exists"
	)

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

		/*
			errors.As() is used instead of if err != nil. It is pgx error wrapping
			it gives us ability to see different error codes and handle them accordingly
			https://github.com/jackc/pgx/wiki/Error-Handling
		*/
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == nonUniqueValueCode {
				return 0, fmt.Errorf("%s: %s: %w", op, userExists, err) // for unique_violation error
			}
			return 0, fmt.Errorf("%s: %w", op, err) // for any other errors
		}

	}

	return id, nil
}

func (r *UserPostgres) GetUser(username, password string) (models.User, error) {
	const op = "repository.postgres.GetUser"

	var user models.User
	query := `SELECT 
					id,
					COALESCE(user_name, '') AS user_name,
					COALESCE(tg_user_name, '') AS tg_user_name,
					COALESCE(first_name, '') AS first_name,
					COALESCE(last_name, '') AS last_name,
					COALESCE(email, '') AS email,
					password_hash, 
					role 
				FROM 
					user_account 
				WHERE user_name=$1 AND password_hash=$2`

	userHabit, err := r.dbpool.Query(context.Background(), query, username, password)
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	user, err = pgx.CollectOneRow(userHabit, pgx.RowToStructByName[models.User])
	if err != nil {
		// fmt.Fprintf(os.Stderr, "error from GetUser: Collect One Row failed: %v\n", err)
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *UserPostgres) GetAllUsers() ([]models.GetUser, error) {
	const op = "repository.postgres.GetAllUsers"

	var users []models.GetUser
	query := `SELECT 
					id,
					COALESCE(user_name, '') AS user_name,
					COALESCE(tg_user_name, '') AS tg_user_name,
					COALESCE(first_name, '') AS first_name,
					COALESCE(last_name, '') AS last_name,
					COALESCE(email, '') AS email,
					role 
				FROM 
					user_account`

	rowsUsers, err := r.dbpool.Query(context.Background(), query)
	if err != nil {
		return users, fmt.Errorf("%s:%s: %w", op, queryErr, err)
	}

	defer rowsUsers.Close()

	users, err = pgx.CollectRows(rowsUsers, pgx.RowToStructByName[models.GetUser])
	if err != nil {
		return users, fmt.Errorf("%s:%s: %w", op, collectErr, err)
	}

	return users, nil
}

func (r *UserPostgres) GetUserById(userId int) (models.GetUser, error) {
	const op = "repository.postgres.GetUserById"

	var user models.GetUser
	query := `SELECT 
					id,
					COALESCE(user_name, '') AS user_name,
					COALESCE(tg_user_name, '') AS tg_user_name,
					COALESCE(first_name, '') AS first_name,
					COALESCE(last_name, '') AS last_name,
					COALESCE(email, '') AS email,
					role 
				FROM 
					user_account
				WHERE id=$1`

	rowUser, err := r.dbpool.Query(context.Background(), query, userId)
	if err != nil {
		return user, fmt.Errorf("%s:%s: %w", op, queryErr, err)
	}

	defer rowUser.Close()

	user, err = pgx.CollectOneRow(rowUser, pgx.RowToStructByName[models.GetUser])
	if err != nil {
		return user, fmt.Errorf("%s:%s: %w", op, collectErr, err)
	}

	return user, nil
}

func (r *UserPostgres) GetUserByTgUsername(TGusername string) (models.GetUser, error) {
	const op = "repository.postgres.GetUserById"

	var user models.GetUser
	query := `SELECT 
					id,
					COALESCE(user_name, '') AS user_name,
					COALESCE(tg_user_name, '') AS tg_user_name,
					COALESCE(first_name, '') AS first_name,
					COALESCE(last_name, '') AS last_name,
					COALESCE(email, '') AS email,
					role 
				FROM 
					user_account
				WHERE tg_user_name=$1`

	rowUser, err := r.dbpool.Query(context.Background(), query, TGusername)
	if err != nil {
		return user, fmt.Errorf("%s:%s: %w", op, queryErr, err)
	}

	defer rowUser.Close()

	user, err = pgx.CollectOneRow(rowUser, pgx.RowToStructByName[models.GetUser])
	if err != nil {
		return user, fmt.Errorf("%s:%s: %w", op, collectErr, err)
	}

	return user, nil
}

func (r *UserPostgres) DeleteUser(userId int) (int, error) {
	const op = "repository.postgres.DeleteUser"

	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
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
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return checkUserId, tx.Commit(context.Background())
}
