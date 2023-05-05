package repository

import (
	"fmt"
	"strings"

	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type HabitTrackerPostgres struct {
	db *sqlx.DB
}

func NewHabitTrackerPostgres(db *sqlx.DB) *HabitTrackerPostgres {
	return &HabitTrackerPostgres{db: db}
}

func (r *HabitTrackerPostgres) GetAll(userId int) ([]models.HabitTracker, error) {
	var trackers []models.HabitTracker
	query := fmt.Sprintf("SELECT ti.id, ti.unit_of_messure, ti.goal, ti.frequency, ti.start_date, ti.end_date, ti.counter, ti.done FROM %s tl INNER JOIN %s ul on tl.id = ul.habit_id WHERE ul.user_id = $1",
		habitTrackerTable, usersHabitsTable)

	if err := r.db.Select(&trackers, query, userId); err != nil {
		return nil, err
	}

	return trackers, nil
}

func (r *HabitTrackerPostgres) GetById(userId, habitId int) (models.HabitTracker, error) {
	var habitTracker models.HabitTracker

	query := fmt.Sprintf("SELECT ti.id, ti.unit_of_messure, ti.goal, ti.frequency, ti.start_date, ti.end_date, ti.counter, ti.done FROM %s tl INNER JOIN %s ul on tl.id = ul.habit_id WHERE ul.user_id = $1 AND ul.habit_id = $2",
		habitTrackerTable, usersHabitsTable)
	err := r.db.Get(&habitTracker, query, userId, habitId)

	return habitTracker, err
}

/*
trackerMapStruct and newTrackerMap are created to pass fields of HabitTracker as variables. We can't
use variables to call struct fields. Struct fields can be called only with direct name specification.
Therefore newTrackerMap func builds and returns a map of strings as keys and struct fields as values.
*/
type trackerMapStruct struct {
	trackerMap map[string]any
}

func newTrackerMap() trackerMapStruct {
	var tracker models.UpdateTrackerInput
	updateTrackMap := map[string]any{
		"unit_of_messure": tracker.UnitOfMessure,
		"goal":            tracker.Goal,
		"frequency":       tracker.Frequency,
		"start_date":      tracker.StartDate,
		"end_date":        tracker.EndDate,
		"counter":         tracker.Counter,
		"done":            tracker.Done,
	}

	return trackerMapStruct{
		trackerMap: updateTrackMap,
	}
}

func (r *HabitTrackerPostgres) Update(userId, habitId int, input models.UpdateTrackerInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	updateInputs := newTrackerMap()

	for key := range updateInputs.trackerMap {
		setValues, args, argId = updateAppender(setValues, updateInputs, args, argId, key)
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.habit_id AND ul.habit_id=$%d AND ul.user_id=$%d",
		habitTrackerTable, setQuery, usersHabitsTable, argId, argId+1)

	args = append(args, habitId, userId)

	logrus.Debugf("updateQuerry: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}

/*
updateAppender func helps to avoid code repetition. It checks and handles updates for each feild of
UpdateTrackerInput.
Example of standard implementation of this update handling in habit_postgres.go file (method Update)
*/
func updateAppender(setValues []string, updateInput trackerMapStruct, args []interface{}, argId int, fieldName string) ([]string, []interface{}, int) {
	if updateInput.trackerMap[fieldName] != nil {
		setValues = append(setValues, fmt.Sprintf("%s=$%d", fieldName, argId))
		args = append(args, updateInput.trackerMap[fieldName])
		argId++
	}

	return setValues, args, argId
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
