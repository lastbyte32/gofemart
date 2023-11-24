package handlers

import (
	"fmt"
	"net/http"
)

func (h *baseHandler) getOrders(w http.ResponseWriter, r *http.Request) {
	user, err := h.getAuthUser(r.Context())
	if err != nil {
		errStr := fmt.Sprintf("get user from ctx failed: %s", err)
		h.ResponseJSONErr(w, http.StatusInternalServerError, errStr)
		return
	}
	orders, err := h.services.GetOrdersByUserID(r.Context(), user.ID)
	if err != nil {
		errStr := fmt.Sprintf("get orders failed: %s", err)
		h.ResponseJSONErr(w, http.StatusInternalServerError, errStr)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	h.ResponseJSON(w, http.StatusOK, orders)
}
