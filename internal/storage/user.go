package storage

import (
	"context"

	"github.com/lastbyte32/gofemart/internal/domain"
)

type User interface {
	Create(ctx context.Context, u domain.User) (*domain.User, error)
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
	AccrualAmountByUserID(ctx context.Context, userID string) (float64, error)
	WithdrawalsAmountByUserID(ctx context.Context, userID string) (float64, error)
}
