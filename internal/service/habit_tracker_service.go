package service

import (
	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/aidos-dev/habit-tracker/internal/repository"
)

type HabitTrackerService struct {
	repo repository.HabitTracker
}

func NewHabitTrackerService(repo repository.HabitTracker) *HabitTrackerService {
	return &HabitTrackerService{repo: repo}
}

func (s *HabitTrackerService) GetAll(userId int) ([]models.HabitTracker, error) {
	return s.repo.GetAll(userId)
}

func (s *HabitTrackerService) GetById(userId, habitId int) (models.HabitTracker, error) {
	return s.repo.GetById(userId, habitId)
}

func (s *HabitTrackerService) Update(userId, habitId int, input models.UpdateTrackerInput) error {
	return s.repo.Update(userId, habitId, input)
}

////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////

/*
temporarily disabled
*/

// func (s *HabitTrackerService) Create(userHabitId int, tracker habit.HabitTracker) (int, error) {
// 	return s.repo.Create(userHabitId, tracker)
// }

////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////

/*
temporarily disabled
*/

// func (s *HabitTrackerService) Delete(userId, habitId int) error {
// 	return s.repo.Delete(userId, habitId)
// }
