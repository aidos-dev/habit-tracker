package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminUserPostgres struct {
	dbpool *pgxpool.Pool
	User
}

func NewAdminUserPostgres(dbpool *pgxpool.Pool) AdminUser {
	return &AdminUserPostgres{dbpool: dbpool}
}

func (r *AdminUserPostgres) GetAllUsers() ([]models.GetUser, error) {
	var users []models.GetUser
	query := `SELECT 
					id,
					user_name, 
					first_name, 
					last_name, 
					email,
					role 
				FROM 
					user_account`

	rowsUsers, err := r.dbpool.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return users, err
	}

	defer rowsUsers.Close()

	users, err = pgx.CollectRows(rowsUsers, pgx.RowToStructByName[models.GetUser])
	if err != nil {
		fmt.Fprintf(os.Stderr, "rowsUsers CollectRows failed: %v\n", err)
		return users, err
	}

	return users, err
}

func (r *AdminUserPostgres) GetUserById(userId int) (models.GetUser, error) {
	var user models.GetUser
	query := `SELECT 
					id,
					user_name, 
					first_name, 
					last_name, 
					email,
					role 
				FROM 
					user_account
				WHERE id=$1`

	rowUser, err := r.dbpool.Query(context.Background(), query, userId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return user, err
	}

	defer rowUser.Close()

	user, err = pgx.CollectOneRow(rowUser, pgx.RowToStructByName[models.GetUser])
	if err != nil {
		fmt.Fprintf(os.Stderr, "rowUser CollectOneRow failed: %v\n", err)
		return user, err
	}

	return user, err
}
