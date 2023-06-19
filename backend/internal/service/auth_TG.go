package service

import "github.com/aidos-dev/habit-tracker/backend/internal/models"

func (s *AuthService) FindTgUser(tgUsername string) (models.GetUser, error) {
	var user models.GetUser
	user, err := s.repo.GetUserByTgUsername(tgUsername)
	if err != nil {
		return user, err
	}

	return user, nil
}
