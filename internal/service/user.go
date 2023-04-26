package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/jwt"
	"github.com/lastbyte32/gofemart/internal/storage"
)

const defaultTTL = time.Minute * 60 * 24

type User interface {
	Registration(ctx context.Context, login, password string) (*domain.User, error)
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	CreateToken(login string) (string, error)
	Withdraw(ctx context.Context, userID, orderNumber string, sum float64) error
	GetBalance(ctx context.Context, userID string) (*domain.Balance, error)
}

type user struct {
	store       storage.User
	auth        jwt.TokenManager
	withdrawSrv storage.Withdraw
}

func (s *user) GetBalance(ctx context.Context, userID string) (*domain.Balance, error) {
	accrualAmount, err := s.store.AccrualAmountByUserID(ctx, userID)
	if err != nil {
		fmt.Println("AccrualAmountByUserID err")
		return nil, err
	}
	withdrawalsAmount, err := s.store.WithdrawalsAmountByUserID(ctx, userID)
	if err != nil {
		fmt.Println("WithdrawalsAmountByUserID err")
		return nil, err
	}
	fmt.Printf("accrual:%f \nwithdrawals: %f\n", accrualAmount, withdrawalsAmount)

	return &domain.Balance{
		Current:   accrualAmount - withdrawalsAmount,
		Withdrawn: withdrawalsAmount,
	}, nil
}

func (s *user) Withdraw(ctx context.Context, userID, orderNumber string, sum float64) error {
	balance, err := s.GetBalance(ctx, userID)
	if err != nil {
		return err
	}
	if balance.Current < sum {
		return domain.ErrNotEnoughFunds
	}
	if _, err = s.withdrawSrv.Create(ctx, userID, orderNumber, sum); err != nil {
		return err
	}
	return nil
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

func NewUserService(store storage.User, auth jwt.TokenManager, w storage.Withdraw) User {
	return &user{store: store, auth: auth, withdrawSrv: w}
}
