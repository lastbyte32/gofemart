package handlers

import (
	"fmt"
	"net/http"
)

func (h *baseHandler) getWithdrawals(w http.ResponseWriter, r *http.Request) {
	user, err := h.getAuthUser(r.Context())
	if err != nil {
		errStr := fmt.Sprintf("get user from ctx failed: %s", err)
		h.ResponseJsonErr(w, http.StatusInternalServerError, errStr)
		return
	}
	withdrawals, err := h.services.Withdrawals(r.Context(), user.ID)
	if err != nil {
		errStr := fmt.Sprintf("get withdrawals failed: %s", err)
		h.ResponseJsonErr(w, http.StatusInternalServerError, errStr)
		return
	}
	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	h.ResponseJson(w, http.StatusOK, withdrawals)

}
