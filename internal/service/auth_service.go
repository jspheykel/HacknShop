package service

import (
	"context"
	"errors"

	"github.com/jspheykel/HacknShop/internal/handlers"
	"github.com/jspheykel/HacknShop/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Users *handlers.UserHandler
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (int64, error) {
	// Check dup
	existing, err := s.Users.FindByUsername(ctx, username)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	return s.Users.Create(ctx, username, email, string(hash), false)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*models.User, error) {
	u, err := s.Users.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("not registered")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}
