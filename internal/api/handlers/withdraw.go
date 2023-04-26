package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/lastbyte32/gofemart/internal/domain"
	"github.com/lastbyte32/gofemart/internal/dto"
)

func (h *baseHandler) withdraw(w http.ResponseWriter, r *http.Request) {
	user, err := h.getAuthUser(r.Context())
	if err != nil {
		errStr := fmt.Sprintf("get user from ctx failed: %s", err)
		h.ResponseJsonErr(w, http.StatusInternalServerError, errStr)
		return
	}

	var inputWithdraw dto.Withdraw
	if err := json.NewDecoder(r.Body).Decode(&inputWithdraw); err != nil {
		h.ResponseJsonErr(w, http.StatusBadRequest, "invalid request")
		return
	}

	if !h.services.IsOrderNumberValid(inputWithdraw.Order) {
		h.ResponseJsonErr(w, http.StatusUnprocessableEntity, "invalid order number")
		return
	}

	if err := h.services.Withdraw(r.Context(), user.ID, inputWithdraw.Order, inputWithdraw.Sum); err != nil {
		if errors.Is(err, domain.ErrNotEnoughFunds) {
			h.ResponseJsonErr(w, http.StatusPaymentRequired, "not enough funds")
			return
		}
		return
	}

}
