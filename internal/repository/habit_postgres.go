package repository

import (
	"fmt"
	"strings"

	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type HabitPostgres struct {
	db *sqlx.DB
}

func NewHabitPostgres(db *sqlx.DB) *HabitPostgres {
	return &HabitPostgres{db: db}
}

func (r *HabitPostgres) Create(userId int, habit models.Habit) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var habitId int
	// create a habit
	createHabitQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", habitsTable)
	rowHabit := tx.QueryRow(createHabitQuery, habit.Title, habit.Description)
	if err := rowHabit.Scan(&habitId); err != nil {
		tx.Rollback()
		return 0, err
	}

	// create an empty tracker for a habit
	var trackerId int
	createHabitTrackerQuery := fmt.Sprintf("INSERT INTO %s (habit_id) VALUES ($1) RETURNING id", habitTrackerTable)
	rowTracker := tx.QueryRow(createHabitTrackerQuery, userId)
	err = rowTracker.Scan(&trackerId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// link habit to a user and a tracker to a habit
	createUsersHabitsQuery := fmt.Sprintf("INSERT INTO %s (user_id, habit_id, habit_tracker_id) VALUES ($1, $2, $3)", usersHabitsTable)
	_, err = tx.Exec(createUsersHabitsQuery, userId, habitId, trackerId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return habitId, tx.Commit()
}

func (r *HabitPostgres) GetAll(userId int) ([]models.Habit, error) {
	var habits []models.Habit
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.habit_id WHERE ul.user_id = $1",
		habitsTable, usersHabitsTable)
	err := r.db.Select(&habits, query, userId)

	return habits, err
}

func (r *HabitPostgres) GetById(userId, habitId int) (models.Habit, error) {
	var habit models.Habit

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.habit_id WHERE ul.user_id = $1 AND ul.habit_id = $2",
		habitsTable, usersHabitsTable)
	err := r.db.Get(&habit, query, userId, habitId)

	return habit, err
}

func (r *HabitPostgres) Delete(userId, habitId int) error {
	queryTracker := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.habit_id = ul.habit_id AND ul.user_id=$1 AND ul.habit_id=$2",
		habitTrackerTable, usersHabitsTable)
	_, err := r.db.Exec(queryTracker, userId, habitId)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.habit_id AND ul.user_id=$1 AND ul.habit_id=$2",
		habitsTable, usersHabitsTable)
	_, err = r.db.Exec(query, userId, habitId)

	return err
}

func (r *HabitPostgres) Update(userId, habitId int, input models.UpdateHabitInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.habit_id AND ul.habit_id=$%d AND ul.user_id=$%d",
		habitsTable, setQuery, usersHabitsTable, argId, argId+1)

	args = append(args, habitId, userId)

	logrus.Debugf("updateQuerry: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
