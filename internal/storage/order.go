package storage

import (
	"context"

	"github.com/lastbyte32/gofemart/internal/domain"
)

type Order interface {
	Create(ctx context.Context, order domain.Order) (*domain.Order, error)
	GetByNumber(ctx context.Context, number string) (*domain.Order, error)
	GetOrdersByUserID(ctx context.Context, userID string) ([]*domain.Order, error)
	GetOrdersUnprocessed(ctx context.Context) ([]*domain.Order, error)
	UpdateOrder(ctx context.Context, info *domain.AccrualOrderInfo) error
}
