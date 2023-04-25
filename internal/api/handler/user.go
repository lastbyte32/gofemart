package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/lastbyte32/gofemart/internal/dto"
	"github.com/lastbyte32/gofemart/internal/service"
)

type userHandler struct {
	//	BaseHandler
	service service.User
	//	auth    service.AuthService
}

func NewUserHandler(s service.User) *userHandler {
	return &userHandler{
		service: s,
	}
}

func (h *userHandler) Routes(router *chi.Mux) {
	router.Post("/api/user/register", h.registration)
	router.Post("/api/user/login", h.authentication)
}

func (h *userHandler) authentication(w http.ResponseWriter, r *http.Request) {
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
	fmt.Printf("LOGIN %s\n", user.Login)

	token, err := h.service.CreateToken(user.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"jwt create err"}`))
		return
	}
	response := fmt.Sprintf(`{"token":"%s"}`, token)
	w.Write([]byte(response))
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
	response := fmt.Sprintf(`{"token":"%s"}`, token)
	w.Write([]byte(response))
}
