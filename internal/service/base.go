package service

import (
	"github.com/lastbyte32/gofemart/internal/jwt"
	"github.com/lastbyte32/gofemart/internal/storage"
)

type Services struct {
	jwt.TokenManager
	User
	Order
}

func New(user storage.User, order storage.Order, withdraw storage.Withdraw, token jwt.TokenManager) *Services {
	return &Services{
		TokenManager: token,
		User:         NewUserService(token, user, withdraw),
		Order:        NewOrderService(order),
	}
}
