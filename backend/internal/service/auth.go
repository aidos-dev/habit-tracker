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
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
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

	return tokenString, err
}

func (s *AuthService) ParseToken(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	// fmt.Printf("service: auth.go: ParseToken: token content: %v\n", token)

	claims := token.Claims.(jwt.MapClaims)

	// fmt.Printf("service: auth.go: ParseToken: claims content: %v\n", claims)

	return claims, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
