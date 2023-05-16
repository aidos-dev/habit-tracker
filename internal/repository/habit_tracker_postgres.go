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

type HabitTrackerPostgres struct {
	dbpool *pgxpool.Pool
}

func NewHabitTrackerPostgres(dbpool *pgxpool.Pool) HabitTracker {
	return &HabitTrackerPostgres{dbpool: dbpool}
}

func (r *HabitTrackerPostgres) GetById(userId, habitId int) (models.HabitTracker, error) {
	var habitTracker models.HabitTracker

	query := `SELECT 
					tl.id, 
					tl.habit_id, 
					COALESCE(tl.unit_of_messure, '-') as unit_of_messure, 
					COALESCE(tl.goal, '-') as goal,
					COALESCE(tl.frequency, '-') as frequency,
					tl.start_date,
					COALESCE(tl.end_date, CURRENT_DATE) as end_date,
					COALESCE(tl.counter, 0) as counter,
					tl.done 
				FROM 
					habit_tracker tl INNER JOIN user_habit ul on tl.id = ul.habit_tracker_id 
				WHERE ul.user_id = $1 AND ul.habit_id = $2`

	/*
		how to add interval to daterime:
		https://www.commandprompt.com/education/postgresql-dateadd-equivalent-how-to-add-interval-to-datetime/
	*/

	rowTracker, err := r.dbpool.Query(context.Background(), query, userId, habitId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error from habit_tracker: GetById: QueryRow failed: %v\n", err)
		return habitTracker, err
	}

	habitTracker, err = pgx.CollectOneRow(rowTracker, pgx.RowToStructByName[models.HabitTracker])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error from habit_tracker: GetById: Collect One Row failed: %v\n", err)
		return habitTracker, err
	}

	return habitTracker, err
}

func (r *HabitTrackerPostgres) GetAll(userId int) ([]models.HabitTracker, error) {
	var trackers []models.HabitTracker
	query := `SELECT 
					tl.id, 
					tl.habit_id, 
					COALESCE(tl.unit_of_messure, '-') as unit_of_messure, 
					COALESCE(tl.goal, '-') as goal,
					COALESCE(tl.frequency, '-') as frequency,
					tl.start_date,
					COALESCE(tl.end_date, CURRENT_DATE) as end_date,
					COALESCE(tl.counter, 0) as counter,
					tl.done 
				FROM 
					habit_tracker tl INNER JOIN user_habit ul on tl.id = ul.habit_tracker_id 
				WHERE 
					ul.user_id = $1`

	rowsTrackers, err := r.dbpool.Query(context.Background(), query, userId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "habit_tracker: GetAll: QueryRow failed: %v\n", err)
		return trackers, err
	}

	defer rowsTrackers.Close()

	trackers, err = pgx.CollectRows(rowsTrackers, pgx.RowToStructByName[models.HabitTracker])
	if err != nil {
		fmt.Fprintf(os.Stderr, "habit_tracker: GetAll: CollectRows failed: %v\n", err)
		return trackers, err
	}

	return trackers, err
}

func (r *HabitTrackerPostgres) Update(userId, habitId int, input models.UpdateTrackerInput) error {
	query := `UPDATE 
					habit_tracker tl 
				SET 
					unit_of_messure=COALESCE($3, unit_of_messure),
					goal=COALESCE($4, goal),
					frequency=COALESCE($5, frequency),
					start_date=COALESCE($6, start_date),
					end_date=COALESCE($7, end_date),
					counter=COALESCE($8, counter),
					done=COALESCE($9, done) 
				FROM user_habit ul 
					WHERE tl.id = ul.habit_tracker_id AND ul.habit_id=$2 AND ul.user_id=$1`

	logrus.Debugf("updateQuerry: %s", query)

	_, err := r.dbpool.Exec(context.Background(), query, userId, habitId, input.UnitOfMessure, input.Goal, input.Frequency, input.StartDate, input.EndDate, input.Counter, input.Done)

	return err
}

////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////

/*
method Create temporarily commented out because I decided to create trackers
together with creating a habit. So any tracker is always linked to a certain habit.
This method is not deleted because it might work for some future functionality
*/

// func (r *HabitTrackerPostgres) Create(userHabitId int, tracker habit.HabitTracker) (int, error) {
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return 0, err
// 	}

// 	var trackerId int
// 	createTrackerQuery := fmt.Sprintf("INSERT INTO %s (unit_of_messure, goal, frequency, start_date, end_date, counter, done) values ($1, $2, $3, $4, $5, $6, $7) RETURNING id", habitTrackerTable)

// 	row := tx.QueryRow(createTrackerQuery, tracker.UnitOfMessure, tracker.Goal, tracker.Frequency, tracker.StartDate, tracker.EndDate, tracker.Counter, tracker.Done)
// 	err = row.Scan(&trackerId)
// 	if err != nil {
// 		tx.Rollback()
// 		return 0, err
// 	}

// 	return trackerId, tx.Commit()
// }

////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////

/*
method Delete is commented out for the same reasons as method Create
*/

// func (r *HabitTrackerPostgres) Delete(userId, habitId int) error {
// 	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.habit_id AND ul.user_id=$1 AND ul.habit_id=$2",
// 		habitTrackerTable, usersHabitsTable)
// 	_, err := r.db.Exec(query, userId, habitId)

// 	return err
// }
