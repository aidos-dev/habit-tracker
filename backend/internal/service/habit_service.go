package service

import (
	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
)

type HabitService struct {
	repo repository.Habit
}

func NewHabitService(repo repository.Habit) Habit {
	return &HabitService{repo: repo}
}

func (s *HabitService) Create(userId int, habit models.Habit) (int, error) {
	return s.repo.Create(userId, habit)
}

func (s *HabitService) GetAll(userId int) ([]models.Habit, error) {
	return s.repo.GetAll(userId)
}

func (s *HabitService) GetById(userId, habitId int) (models.Habit, error) {
	return s.repo.GetById(userId, habitId)
}

func (s *HabitService) Delete(userId, habitId int) error {
	return s.repo.Delete(userId, habitId)
}

func (s *HabitService) Update(userId, habitId int, input models.UpdateHabitInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, habitId, input)
}
