package postgres

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/storage"
)

var _ storage.User = (*userStore)(nil)

const (
	scheme      = "public"
	table       = "users"
	tableScheme = scheme + "." + table
)

const (
	sqlAccrualAmountByUserID     = "SELECT COALESCE(SUM(o.accrual), 0) as sum FROM users u LEFT JOIN orders o ON u.id = o.user_id WHERE u.id = $1 GROUP BY u.id"
	sqlWithdrawalsAmountByUserID = "SELECT COALESCE(SUM(o.sum), 0) as sum FROM users u LEFT JOIN withdrawals o ON u.id = o.user_id WHERE u.id = $1 GROUP BY u.id"

	sqlGetByLogin = "SELECT * FROM " + tableScheme + " WHERE login = $1"

	sqlGetByTelegramID = "SELECT * FROM " + tableScheme + " WHERE telegram_id = $1"

	sqlGetByID = "SELECT * FROM " + tableScheme + " WHERE id = $1"
	//  sqlAllMetrics    = `SELECT * FROM metrics`
	//  sqlUpdateCounter = `UPDATE metrics SET counter = $1 WHERE id = $2`
	//  sqlUpdateGauge   = `UPDATE metrics SET gauge = cast($1 as double precision) WHERE id = $2`
	sqlInsertUser = "INSERT INTO " + tableScheme + " (id, login, password) VALUES (:id, :login, :password)"
	//  sqlInsertGauge   = `INSERT INTO metrics (id, mtype, gauge) VALUES($1,$2,$3)`
)

type userStore struct {
	db *sqlx.DB
}

func (s *userStore) GetByTelegramID(ctx context.Context, telegramID int) (*domain.User, error) {
	var user domain.User
	err := s.db.GetContext(ctx, &user, sqlGetByTelegramID, telegramID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, err
}

func (s *userStore) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	if err := s.db.GetContext(ctx, &user, sqlGetByLogin, login); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *userStore) All(ctx context.Context) ([]*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *userStore) Create(ctx context.Context, u domain.User) (*domain.User, error) {
	_, err := s.db.NamedQueryContext(ctx, sqlInsertUser, u)
	if err != nil {
		return nil, errors.Wrap(err, "store err")
	}
	return &u, nil
}

func (s *userStore) AccrualAmountByUserID(ctx context.Context, userID string) (float64, error) {
	var result struct {
		Sum float64
	}
	if err := s.db.GetContext(ctx, &result, sqlAccrualAmountByUserID, userID); err != nil {
		return 0, err
	}
	return result.Sum, nil
}
func (s *userStore) WithdrawalsAmountByUserID(ctx context.Context, userID string) (float64, error) {
	var result struct {
		Sum float64
	}
	if err := s.db.GetContext(ctx, &result, sqlWithdrawalsAmountByUserID, userID); err != nil {
		return 0, err
	}
	return result.Sum, nil
}
func NewUserStore(db *sqlx.DB) storage.User {
	return &userStore{db: db}
}
