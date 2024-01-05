package postgres

import (
	"context"
	"fmt"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HabitPostgres struct {
	dbpool *pgxpool.Pool
}

func NewHabitPostgres(dbpool *pgxpool.Pool) repository.Habit {
	return &HabitPostgres{dbpool: dbpool}
}

func (r *HabitPostgres) Create(userId int, habit models.Habit) (int, error) {
	const op = "repository.postgres.habit_postgres.Create"

	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
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
		return 0, fmt.Errorf("%s:%s: %w", op, habitTable, err)
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
		return 0, fmt.Errorf("%s:%s: %w", op, trackerTable, err)
	}

	// link habit to a user and a tracker to a habit
	createUsersHabitsQuery := `INSERT INTO 
										user_habit (user_id, habit_id, habit_tracker_id) 
										VALUES ($1, $2, $3)`

	_, err = tx.Exec(context.Background(), createUsersHabitsQuery, userId, habitId, trackerId)
	if err != nil {
		tx.Rollback(context.Background())
		return 0, fmt.Errorf("%s:%s: %w", op, userHabitTable, err)
	}

	return habitId, tx.Commit(context.Background())
}

func (r *HabitPostgres) GetAll(userId int) ([]models.Habit, error) {
	const op = "repository.postgres.habit_postgres.GetAll"

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
		return habits, fmt.Errorf("%s:%s: %w", op, queryErr, err)
	}

	defer rowsHabits.Close()

	habits, err = pgx.CollectRows(rowsHabits, pgx.RowToStructByName[models.Habit])
	if err != nil {
		return habits, fmt.Errorf("%s:%s: %w", op, collectErr, err)
	}

	return habits, err
}

func (r *HabitPostgres) GetById(userId, habitId int) (models.Habit, error) {
	const op = "repository.postgres.habit_postgres.GetById"

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
		return habit, fmt.Errorf("%s:%s: %w", op, queryErr, err)
	}

	defer rowHabit.Close()

	habit, err = pgx.CollectOneRow(rowHabit, pgx.RowToStructByName[models.Habit])
	if err != nil {
		return habit, fmt.Errorf("%s:%s: %w", op, collectErr, err)
	}

	return habit, err
}

func (r *HabitPostgres) Delete(userId, habitId int) error {
	const op = "repository.postgres.habit_postgres.Delete"

	tx, err := r.dbpool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
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
		return fmt.Errorf("%s:%s: %w", op, trackerTable, err)
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
		return fmt.Errorf("%s:%s: %w", op, habitTable, err)
	}

	return tx.Commit(context.Background())
}

func (r *HabitPostgres) Update(userId, habitId int, input models.UpdateHabitInput) error {
	const op = "repository.postgres.habit_postgres.Update"

	query := `UPDATE 
					habit tl 
				SET 
					title=COALESCE($3, title), 
					description=COALESCE($4, description)
				FROM user_habit ul 
					WHERE tl.id = ul.habit_id AND ul.user_id=$1 AND ul.habit_id=$2
					RETURNING tl.id`

	var checkHabitId int

	rowHabit := r.dbpool.QueryRow(context.Background(), query, userId, habitId, input.Title, input.Description)
	err := rowHabit.Scan(&checkHabitId)
	if err != nil {
		return fmt.Errorf("%s:%s: %w", op, scanErr, err)
	}

	return err
}
