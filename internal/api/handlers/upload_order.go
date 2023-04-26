package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/lastbyte32/gofemart/internal/domain"
)

func (h *baseHandler) uploadOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.ResponseJsonErr(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	orderNumber := string(body)
	if orderNumber == "" {
		h.ResponseJsonErr(w, http.StatusBadRequest, "empty order number")
		return
	}

	if !h.services.IsOrderNumberValid(orderNumber) {
		h.ResponseJsonErr(w, http.StatusUnprocessableEntity, "invalid order number")
		return
	}
	user, err := h.getAuthUser(r.Context())
	if err != nil {
		errStr := fmt.Sprintf("get user from ctx failed: %s", err)
		h.ResponseJsonErr(w, http.StatusInternalServerError, errStr)
		return
	}

	err = h.services.UploadOrder(r.Context(), orderNumber, user.ID)
	if err != nil && errors.Is(err, domain.ErrDuplicateOrderUploadOtherUser) {
		h.ResponseJsonErr(w, http.StatusConflict, "duplicate order upload other user")
		return
	}

	if err != nil && errors.Is(err, domain.ErrDuplicateOrderUploadSameUser) {
		h.ResponseJsonErr(w, http.StatusConflict, "duplicate order upload same user")
		return

	}
	if err != nil {
		errStr := fmt.Sprintf("upload order failed: %s", err)
		h.ResponseJsonErr(w, http.StatusInternalServerError, errStr)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
