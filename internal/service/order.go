package service

import (
	"context"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/storage"
	"github.com/lastbyte32/gofemart/internal/util/luhn"
)

type Order interface {
	IsNumberValid(number string) bool
	Create(ctx context.Context, number, userId string) error
	GetByUserID(ctx context.Context, userID string) ([]*domain.Order, error)
}

type order struct {
	store storage.Order
}

func (s *order) GetByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	return s.store.GetOrdersByUserID(ctx, userID)
}

func (s *order) Create(ctx context.Context, number, userId string) error {
	newOrder := domain.Order{
		Number: number,
		UserID: userId,
		Status: domain.OrderNew,
	}
	order, err := s.store.GetByNumber(ctx, number)
	if err != nil {
		return err
	}

	if order != nil && order.UserID == userId {
		return domain.ErrDuplicateOrderUploadSameUser
	}

	if order != nil && order.UserID != userId {
		return domain.ErrDuplicateOrderUploadOtherUser
	}

	if _, err := s.store.Create(ctx, newOrder); err != nil {
		return err
	}
	return nil
}

func (s *order) IsNumberValid(number string) bool {
	return luhn.Validation(number)
}

func NewOrderService(store storage.Order) Order {
	return &order{store: store}
}
