package withdraw

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/storage"
)

type store struct {
	db *sqlx.DB
}

func (s *store) GetByUserID(ctx context.Context, userID string) ([]*domain.Withdraw, error) {
	var withdrawals []*domain.Withdraw
	if err := s.db.SelectContext(ctx, &withdrawals, sqlGetByUserID, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return withdrawals, nil
}

func (s *store) Create(ctx context.Context, userID, orderNumber string, sum float64) (*domain.Withdraw, error) {
	withdraw := domain.Withdraw{
		UserID:      userID,
		OrderNumber: orderNumber,
		Sum:         sum,
	}
	_, err := s.db.NamedExecContext(ctx, sqlInsert, withdraw)
	if err != nil {
		return nil, errors.Wrap(err, "store err")
	}
	return &withdraw, nil
}

func NewStore(db *sqlx.DB) storage.Withdraw {
	return &store{db: db}
}
