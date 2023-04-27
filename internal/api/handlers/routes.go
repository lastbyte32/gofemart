package handlers

import (
	"github.com/go-chi/chi"
)

func (h *baseHandler) Routes(router *chi.Mux) {
	router.Post("/api/user/register", h.registration)
	router.Post("/api/user/login", h.authentication)
	router.Group(func(r chi.Router) {
		r.Use(h.services.JWTMiddleware)
		r.Post("/api/user/orders", h.uploadOrder)
		r.Post("/api/user/balance/withdraw", h.withdraw)
		r.Get("/api/user/orders", h.getOrders)
		r.Get("/api/user/balance", h.getBalance)
		r.Get("/api/user/withdrawals", h.getWithdrawals)

	})
}
