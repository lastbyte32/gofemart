package storage

import (
	"context"

	"github.com/lastbyte32/gofemart/internal/domain"
)

type Order interface {
	Create(ctx context.Context, order domain.Order) (*domain.Order, error)
	GetByNumber(ctx context.Context, number string) (*domain.Order, error)
	GetOrdersByUserID(ctx context.Context, userID string) ([]*domain.Order, error)
	GetByUserIdAndNumber(ctx context.Context, userID, number string) (*domain.Order, error)
	//GetUnprocessedOrders() ([]domain.Order, error)
	//ProcessOrderAccrual(orderNumber string, status string, accrual float64) error
}
