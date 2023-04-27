package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lastbyte32/gofemart/internal/dto"
)

func (h *baseHandler) authentication(w http.ResponseWriter, r *http.Request) {
	var inputUser dto.Credentials

	if err := json.NewDecoder(r.Body).Decode(&inputUser); err != nil {
		h.ResponseJSONErr(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := inputUser.Validate(); err != nil {
		errStr := fmt.Sprintf("validation failed: %s", err)
		h.ResponseJSONErr(w, http.StatusBadRequest, errStr)
		return
	}

	user, err := h.services.GetUserByLogin(r.Context(), inputUser.Login)
	if err != nil {
		errStr := fmt.Sprintf("get user [%s] failed: %s", user.Login, err)
		h.ResponseJSONErr(w, http.StatusInternalServerError, errStr)
		return
	}

	isPasswordCompare := h.services.CheckPassword(user.Password, inputUser.Password)
	if user == nil || isPasswordCompare != nil {
		h.ResponseJSONErr(w, http.StatusUnauthorized, "credentials don't match")
		return
	}

	token, errT := h.services.GenerateBearerToken(user.Login)
	if errT != nil {
		errStr := fmt.Sprintf("generate token failed: %s", errT)
		h.ResponseJSONErr(w, http.StatusInternalServerError, errStr)
		return
	}

	w.Header().Set(AuthorizationHeader, token)
}
