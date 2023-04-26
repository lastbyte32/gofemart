package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/service"
)

const AuthorizationHeader = "Authorization"

type baseHandler struct {
	services *service.Services
	logger   *zerolog.Logger
}

func New(l *zerolog.Logger, s *service.Services) *baseHandler {
	return &baseHandler{
		logger:   l,
		services: s,
	}
}

func (h *baseHandler) getAuthUser(ctx context.Context) (*domain.User, error) {
	id := ctx.Value("userID").(string)
	if id == "" {
		return nil, errors.New("value not found")
	}
	user, err := h.services.GetUserByLogin(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (h *baseHandler) ResponseJson(w http.ResponseWriter, status int, result any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	var payload []byte
	var err error
	payload, err = json.Marshal(result)
	if err != nil {
		w.Write([]byte(`{"error": "marshalError","details": "` + err.Error() + `"}`))
		return
	}
	w.Write(payload)
}

func (h *baseHandler) ResponseJsonErr(w http.ResponseWriter, status int, message string) {
	h.logger.Warn().Msg(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	payload := fmt.Sprintf(`{"error": "%s"}`, message)
	w.Write([]byte(payload))
}
