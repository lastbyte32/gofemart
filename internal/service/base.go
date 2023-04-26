package service

import (
	"context"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/service/jwt"
	"github.com/lastbyte32/gofemart/internal/service/order"
	"github.com/lastbyte32/gofemart/internal/service/user"
	"github.com/lastbyte32/gofemart/internal/storage"
)

type accrualGetter interface {
	GetOrder(ctx context.Context, number string) (*domain.AccrualOrderInfo, error)
}

type Services struct {
	jwt.TokenManager
	user.User
	order.Order
}

func New(u storage.User, o storage.Order, withdraw storage.Withdraw, token jwt.TokenManager) *Services {
	return &Services{
		TokenManager: token,
		User:         user.NewService(token, u, withdraw),
		Order:        order.NewService(o),
	}
}
