package storage

import (
	"context"

	"github.com/lastbyte32/gofemart/internal/domain"
)

type Withdraw interface {
	GetByUserID(ctx context.Context, userID string) ([]*domain.Withdraw, error)
	Create(ctx context.Context, userID, orderNumber string, sum float64) (*domain.Withdraw, error)
}
