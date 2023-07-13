package service

import (
	"fmt"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
)

func (s *AuthService) FindTgUser(tgUsername string) (models.GetUser, error) {
	const op = "service.auth_TG_service.FindTgUser"

	var user models.GetUser
	user, err := s.repo.GetUserByTgUsername(tgUsername)
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, err
}
