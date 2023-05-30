package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type HabitPostgres struct {
	dbpool *pgxpool.Pool
}

func NewHabitPostgres(dbpool *pgxpool.Pool) Habit {
	return &HabitPostgres{dbpool: dbpool}
}

func (r *HabitPostgres) Create(userId int, habit models.Habit) (int, error) {
	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return 0, err
	}

	var habitId int
	// create a habit
	createHabitQuery := `INSERT INTO 
								habit (title, description) 
								VALUES ($1, $2) 
							RETURNING id`

	rowHabit := tx.QueryRow(context.Background(), createHabitQuery, habit.Title, habit.Description)
	if err := rowHabit.Scan(&habitId); err != nil {
		tx.Rollback(context.Background())
		return 0, err
	}

	// create an empty tracker for a habit
	var trackerId int
	createHabitTrackerQuery := `INSERT INTO 
										habit_tracker (habit_id) 
										VALUES ($1) 
									RETURNING id`

	rowTracker := tx.QueryRow(context.Background(), createHabitTrackerQuery, habitId)
	err = rowTracker.Scan(&trackerId)
	if err != nil {
		tx.Rollback(context.Background())
		return 0, err
	}

	// link habit to a user and a tracker to a habit
	createUsersHabitsQuery := `INSERT INTO 
										user_habit (user_id, habit_id, habit_tracker_id) 
										VALUES ($1, $2, $3)`

	_, err = tx.Exec(context.Background(), createUsersHabitsQuery, userId, habitId, trackerId)
	if err != nil {
		tx.Rollback(context.Background())
		return 0, err
	}

	return habitId, tx.Commit(context.Background())
}

func (r *HabitPostgres) GetAll(userId int) ([]models.Habit, error) {
	var habits []models.Habit
	query := `SELECT 
					tl.id, 
					tl.title, 
					tl.description 
				FROM 
					habit tl INNER JOIN user_habit ul on tl.id = ul.habit_id 
				WHERE ul.user_id = $1`

	rowsHabits, err := r.dbpool.Query(context.Background(), query, userId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return habits, err
	}

	defer rowsHabits.Close()

	habits, err = pgx.CollectRows(rowsHabits, pgx.RowToStructByName[models.Habit])
	if err != nil {
		fmt.Fprintf(os.Stderr, "rowsHabits CollectRows failed: %v\n", err)
		return habits, err
	}

	return habits, err
}

func (r *HabitPostgres) GetById(userId, habitId int) (models.Habit, error) {
	var habit models.Habit

	query := `SELECT 
					tl.id, 
					tl.title, 
					tl.description 
				FROM 
					habit tl INNER JOIN user_habit ul on tl.id = ul.habit_id 
				WHERE ul.user_id = $1 AND ul.habit_id = $2`

	rowHabit, err := r.dbpool.Query(context.Background(), query, userId, habitId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error from GetById: QueryRow failed: %v\n", err)
		return habit, err
	}

	defer rowHabit.Close()

	habit, err = pgx.CollectOneRow(rowHabit, pgx.RowToStructByName[models.Habit])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error from GetById: Collect One Row failed: %v\n", err)
		return habit, err
	}

	return habit, err
}

func (r *HabitPostgres) Delete(userId, habitId int) error {
	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return err
	}

	queryTracker := `DELETE FROM 
							habit_tracker tl USING user_habit ul 
						WHERE tl.id = ul.habit_tracker_id AND ul.user_id=$1 AND ul.habit_id=$2
						RETURNING tl.id`

	var checkTrackerId int

	rowTracker := tx.QueryRow(context.Background(), queryTracker, userId, habitId)
	err = rowTracker.Scan(&checkTrackerId)
	if err != nil {
		tx.Rollback(context.Background())
		fmt.Printf("err: repository: habit_postgres.go: Delete: rowTracker.Scan: habit tracker doesn't exist: %v\n", err)
		return err
	}

	query := `DELETE FROM 
					habit tl USING user_habit ul 
				WHERE tl.id = ul.habit_id AND ul.user_id=$1 AND ul.habit_id=$2
				RETURNING tl.id`

	var checkHabitId int

	rowHabit := tx.QueryRow(context.Background(), query, userId, habitId)
	err = rowHabit.Scan(&checkHabitId)
	if err != nil {
		tx.Rollback(context.Background())
		fmt.Printf("err: repository: habit_postgres.go: Delete: rowHabit.Scan: habit doesn't exist: %v\n", err)
		return err
	}

	return tx.Commit(context.Background())
}

func (r *HabitPostgres) Update(userId, habitId int, input models.UpdateHabitInput) error {
	query := `UPDATE 
					habit tl 
				SET 
					title=COALESCE($3, title), 
					description=COALESCE($4, description)
				FROM user_habit ul 
					WHERE tl.id = ul.habit_id AND ul.user_id=$1 AND ul.habit_id=$2
					RETURNING tl.id`

	logrus.Debugf("updateQuerry: %s", query)

	var checkHabitId int

	rowHabit := r.dbpool.QueryRow(context.Background(), query, userId, habitId, input.Title, input.Description)
	err := rowHabit.Scan(&checkHabitId)
	if err != nil {

		fmt.Printf("err: repository: habit_postgres.go: Update: rowHabit.Scan: habit doesn't exist: %v\n", err)
		return err
	}

	return err
}
