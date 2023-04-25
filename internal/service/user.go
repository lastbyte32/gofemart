package service

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/jwt"
	"github.com/lastbyte32/gofemart/internal/storage"
)

const defaultTTL = time.Minute * 60

type User interface {
	Registration(ctx context.Context, login, password string) (*domain.User, error)
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	CreateToken(login string) (string, error)
}

type user struct {
	store storage.User
	auth  jwt.TokenManager
}

func (s *user) CreateToken(login string) (string, error) {
	return s.auth.NewJWT(login, defaultTTL)
}

func (s *user) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	return s.store.GetByLogin(ctx, login)
}

func (s *user) Registration(ctx context.Context, login, password string) (*domain.User, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Errorf("failed to generate UUID: %v", err)
	}
	user := domain.User{
		ID:       id.String(),
		Login:    login,
		Password: password,
	}
	return s.store.Create(ctx, user)
}

func NewUserService(store storage.User, auth jwt.TokenManager) User {
	return &user{store: store, auth: auth}
}
