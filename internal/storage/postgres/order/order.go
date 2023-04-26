package order

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/storage"
)

type orderStore struct {
	db *sqlx.DB
}

func (s *orderStore) Create(ctx context.Context, order domain.Order) (*domain.Order, error) {
	_, err := s.db.NamedQueryContext(ctx, sqlInsert, order)
	if err != nil {
		return nil, errors.Wrap(err, "store err")
	}
	return &order, nil
}

func (s *orderStore) GetByNumber(ctx context.Context, number string) (*domain.Order, error) {
	var user domain.Order
	if err := s.db.GetContext(ctx, &user, sqlGetByNumber, number); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *orderStore) GetOrdersByUserID(ctx context.Context, userID string) ([]*domain.Order, error) {
	var orders []*domain.Order
	if err := s.db.SelectContext(ctx, &orders, sqlGetByUserID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return orders, nil
}

func NewOrderStore(db *sqlx.DB) storage.Order {
	return &orderStore{db: db}
}
