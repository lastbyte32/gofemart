package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lastbyte32/gofemart/internal/dto"
)

func (h *baseHandler) registration(w http.ResponseWriter, r *http.Request) {
	var user dto.Credentials

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.ResponseJSONErr(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := user.Validate(); err != nil {
		errStr := fmt.Sprintf("validation failed: %s", err)
		h.ResponseJSONErr(w, http.StatusBadRequest, errStr)
		return
	}

	exist, err := h.services.GetUserByLogin(r.Context(), user.Login)
	if err != nil {
		errStr := fmt.Sprintf("get user [%s] failed: %s", user.Login, err)
		h.ResponseJSONErr(w, http.StatusInternalServerError, errStr)
		return
	}

	if exist != nil {
		errStr := fmt.Sprintf("login [%s] is already use", user.Login)
		h.ResponseJSONErr(w, http.StatusConflict, errStr)
		return
	}

	if _, err := h.services.Registration(r.Context(), user.Login, user.Password); err != nil {
		errStr := fmt.Sprintf("registration [%s] failed: %s", user.Login, err)
		h.ResponseJSONErr(w, http.StatusInternalServerError, errStr)
		return
	}

	token, errT := h.services.GenerateBearerToken(user.Login)
	if errT != nil {
		errStr := fmt.Sprintf("generate token failed: %s", errT)
		h.ResponseJSONErr(w, http.StatusInternalServerError, errStr)
		return
	}
	w.Header().Set(AuthorizationHeader, token)
	h.logger.Info().Str("login", user.Login).Msg("user created successfully")
}
