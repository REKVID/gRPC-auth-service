package auth

import (
	"auth2/internal/store"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store  *store.Store
	jwtKey []byte
}

func NewService(store *store.Store, jwtKey string) *Service {
	return &Service{store: store, jwtKey: []byte(jwtKey)}
}

func (s *Service) Register(email, password string) (uint, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err := s.store.CreateUser(email, string(hash))
	if err != nil {
		return 0, err
	}
	user, _ := s.store.GetUserByEmail(email)
	return user.ID, nil
}

func (s *Service) Login(email, password string) (string, error) {
	user, err := s.store.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", fmt.Errorf("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
	})
	tokenStr, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
