package service

import (
	"github.com/aidos-dev/habit-tracker"
	"github.com/aidos-dev/habit-tracker/pkg/repository"
)

type HabitTrackerService struct {
	repo repository.HabitTracker
}

func NewHabitTrackerService(repo repository.HabitTracker) *HabitTrackerService {
	return &HabitTrackerService{repo: repo}
}

func (s *HabitTrackerService) GetAll(userId int) ([]habit.HabitTracker, error) {
	return s.repo.GetAll(userId)
}

func (s *HabitTrackerService) GetById(userId, habitId int) (habit.HabitTracker, error) {
	return s.repo.GetById(userId, habitId)
}

func (s *HabitTrackerService) Update(userId, habitId int, input habit.UpdateTrackerInput) error {
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
