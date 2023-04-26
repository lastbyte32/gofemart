package handlers

import (
	"fmt"
	"net/http"
)

func (h *baseHandler) getBalance(w http.ResponseWriter, r *http.Request) {
	user, err := h.getAuthUser(r.Context())
	if err != nil {
		errStr := fmt.Sprintf("get user from ctx failed: %s", err)
		h.ResponseJsonErr(w, http.StatusInternalServerError, errStr)
		return
	}

	balance, err := h.services.GetBalance(r.Context(), user.ID)
	if err != nil {
		errStr := fmt.Sprintf("get balance failed: %s", err)
		h.ResponseJsonErr(w, http.StatusInternalServerError, errStr)
		return
	}

	h.ResponseJson(w, http.StatusOK, balance)
}
