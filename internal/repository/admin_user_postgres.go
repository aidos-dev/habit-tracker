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
		fmt.Fprintf(os.Stderr, "rowsHabits CollectRows failed: %v\n", err)
		return users, err
	}

	return users, err
}
