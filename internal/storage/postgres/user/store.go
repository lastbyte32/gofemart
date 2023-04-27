package user

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/storage"
)

type store struct {
	db *sqlx.DB
}

func (s *store) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	if err := s.db.GetContext(ctx, &user, sqlGetByLogin, login); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *store) Create(ctx context.Context, u domain.User) (*domain.User, error) {
	_, err := s.db.NamedExecContext(ctx, sqlInsertUser, u)
	if err != nil || err == sql.ErrNoRows {
		return nil, errors.Wrap(err, "store err")
	}
	return &u, nil
}

func (s *store) AccrualAmountByUserID(ctx context.Context, userID string) (float64, error) {
	var result struct {
		Sum float64
	}
	if err := s.db.GetContext(ctx, &result, sqlAccrualAmountByUserID, userID); err != nil {
		return 0, err
	}
	return result.Sum, nil
}
func (s *store) WithdrawalsAmountByUserID(ctx context.Context, userID string) (float64, error) {
	var result struct {
		Sum float64
	}
	if err := s.db.GetContext(ctx, &result, sqlWithdrawalsAmountByUserID, userID); err != nil {
		return 0, err
	}
	return result.Sum, nil
}
func NewStore(db *sqlx.DB) storage.User {
	return &store{db: db}
}
