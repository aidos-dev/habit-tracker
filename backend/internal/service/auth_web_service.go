package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/aidos-dev/habit-tracker/backend/internal/repository"
	"github.com/golang-jwt/jwt"
)

const (
	salt       = "lk6vm9vkf47#b@7kdn4nv"
	signingKey = "436k@5*6lklj4t6^k4$4#(*&$"
	tokenTTL   = 12 * time.Hour
)

type AuthService struct {
	repo repository.User
}

func NewAuthService(repo repository.User) Authorization {
	return &AuthService{repo: repo}
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	const op = "service.auth_web_service.GenerateToken"

	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	claims := &jwt.MapClaims{
		"iss":       "issuer",
		"issuedAt":  time.Now().Unix(),
		"expiresAt": time.Now().Add(tokenTTL).Unix(),
		"data": map[string]any{
			"userId":   user.Id,
			"userRole": user.Role,
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, err
}

func (s *AuthService) ParseToken(accessToken string) (jwt.MapClaims, error) {
	const op = "service.auth_web_service.ParseToken"

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims, err
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
