package models

import "time"

type Habit struct {
	Id          int    `json:"habitId"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type HabitTracker struct {
	Id            int       `json:"trackerId"`
	HabitId       int       `json:"habitId"`
	UnitOfMessure string    `json:"unit_of_messure"`
	Goal          string    `json:"goal"`
	Frequency     string    `json:"frequency"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Counter       int       `json:"counter"`
	Done          bool      `json:"done"`
}
