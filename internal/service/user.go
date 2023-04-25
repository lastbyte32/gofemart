package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/storage"
)

type User interface {
	Registration(ctx context.Context, login, password string) (*domain.User, error)
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
}

type user struct {
	store storage.User
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

func NewUserService(store storage.User) User {
	return &user{store: store}
}
