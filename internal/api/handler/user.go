package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/dto"
	"github.com/lastbyte32/gofemart/internal/jwt"
	"github.com/lastbyte32/gofemart/internal/service"
)

type userHandler struct {
	//	BaseHandler
	service service.User
	auth    jwt.TokenManager
	order   service.Order
}

func NewUserHandler(s service.User, a jwt.TokenManager, o service.Order) *userHandler {
	return &userHandler{
		service: s,
		auth:    a,
		order:   o,
	}
}

func (h *userHandler) Routes(router *chi.Mux) {
	router.Post("/api/user/register", h.registration)
	router.Post("/api/user/login", h.authentication)
	router.Group(func(r chi.Router) {
		r.Use(h.auth.JWTMiddleware)
		r.Post("/api/user/orders", h.uploadOrder)
		r.Get("/api/user/orders", h.getOrders)
		r.Get("/api/user/balance", h.getBalance)
		r.Post("/api/user/balance/withdraw", h.withdraw)

	})

}
func (h *userHandler) withdraw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	user, errGet := h.getAuthUser(r.Context())
	if errGet != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"get user err"}`))
		return
	}

	var inputWithdraw dto.Withdraw
	if err := json.NewDecoder(r.Body).Decode(&inputWithdraw); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(`{"error":"invalid request body"}`))
		return
	}

	if !h.order.IsNumberValid(inputWithdraw.Order) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error":"invalid order number."}`))
		return
	}

	if err := h.service.Withdraw(r.Context(), user.ID, inputWithdraw.Order, inputWithdraw.Sum); err != nil {
		if errors.Is(err, domain.ErrNotEnoughFunds) {
			w.WriteHeader(http.StatusPaymentRequired)
			_, err = w.Write([]byte(`{"error":"ErrNotEnoughFunds"}`))
			return
		}
		return
	}

}

func (h *userHandler) getBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	user, errGet := h.getAuthUser(r.Context())
	if errGet != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"get user err"}`))
		return
	}

	balance, err := h.service.GetBalance(r.Context(), user.ID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"get balance err"}`))
		return
	}

	response, err := json.Marshal(balance)
	if err != nil {
		return
	}

	_, err = w.Write(response)

}

func (h *userHandler) getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	user, errGet := h.getAuthUser(r.Context())
	if errGet != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"get user err"}`))
		return
	}

	orders, err := h.order.GetByUserID(r.Context(), user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"get orders err"}`))
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, err := json.Marshal(orders)
	if err != nil {
		return
	}

	_, err = w.Write(response)
}

func (h *userHandler) uploadOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"something went wrong."}`))
		return
	}

	orderNumber := string(body)
	// не корректный запрос
	if orderNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid request body."}`))
		return
	}
	// не валидный номер заказа
	if !h.order.IsNumberValid(orderNumber) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error":"invalid order number."}`))
		return
	}
	user, errGet := h.getAuthUser(r.Context())
	if errGet != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"get user err"}`))
		return
	}

	err = h.order.Create(r.Context(), orderNumber, user.ID)
	if err != nil && errors.Is(err, domain.ErrDuplicateOrderUploadOtherUser) {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"error":"ErrDuplicateOrderUploadOtherUser"}`))
		return
	}

	if err != nil && errors.Is(err, domain.ErrDuplicateOrderUploadSameUser) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"error":"ErrDuplicateOrderUploadSameUser"}`))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		_, _ = w.Write([]byte(`{"error":"create order err"}`))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *userHandler) getAuthUser(ctx context.Context) (*domain.User, error) {
	id := ctx.Value("userID").(string)
	if id == "" {
		return nil, errors.New("value not found")
	}
	user, err := h.service.GetUserByLogin(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil

}

func (h *userHandler) authentication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var inputUser dto.CreateUser

	// 400
	if err := json.NewDecoder(r.Body).Decode(&inputUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(`{"error":"invalid request body"}`))
		return
	}

	if err := inputUser.Validate(); err != nil {
		//details, _ := json.Marshal(err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(`{"error":"validation error"}`))
		return
	}

	user, err := h.service.GetUserByLogin(r.Context(), inputUser.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(`{"error":"user get err"}`))
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write([]byte(`{"error":"user not found"}`))
		return
	}

	if user.Password != inputUser.Password {
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write([]byte(`{"error":"wrong password"}`))
		return
	}

	token, err := h.service.CreateToken(user.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"jwt create err"}`))
		return
	}
	response := fmt.Sprintf("Bearer %s", token)
	w.Header().Set("Authorization", response)
}

func (h *userHandler) registration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var user dto.CreateUser

	// 400
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(`{"error":"invalid request body"}`))
		return
	}

	// Валидация
	if err := user.Validate(); err != nil {
		//details, _ := json.Marshal(err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(`{"error":"validation error"}`))
		return
	}

	exist, err := h.service.GetUserByLogin(r.Context(), user.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"GetUserByLogin failed"}`))
		return
	}
	if exist != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error":"login exist"}`))
		return
	}

	if _, err := h.service.Registration(r.Context(), user.Login, user.Password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"service create err"}`))
		return
	}
	token, err := h.service.CreateToken(user.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"jwt create err"}`))
		return
	}
	response := fmt.Sprintf("Bearer %s", token)
	w.Header().Set("Authorization", response)
}
