package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"

	"github.com/sreway/gophermart/internal/domain"
)

func (d *delivery) orderAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	userID, err := uuid.Parse(r.Context().Value(ctxUserID{}).(string))
	if err != nil {
		d.logger.Error("failed parse user id from context", slog.Any("err", err),
			slog.String("handler", "orderAdd"))
		handelOrderErr(w, domain.ErrUserUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		d.logger.Error("failed decode request body",
			slog.String("handler", "orderAdd"),
			slog.Any("err", err))
		handelUserErr(w, ErrDecodeBody)
		return
	}

	number, err := strconv.Atoi(string(body))
	if err != nil {
		d.logger.Error("failed get order number from request body", slog.Any("err", err),
			slog.String("handler", "orderAdd"))
		handelOrderErr(w, domain.ErrIncorrectData)
		return
	}

	orderNumber, err := domain.NewOrderNumber(number)
	if err != nil {
		d.logger.Error("failed create order number from integer value", slog.Any("err", err))
		handelOrderErr(w, err)
		return
	}

	err = d.order.Add(r.Context(), userID, *orderNumber)
	if err != nil {
		handelOrderErr(w, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (d *delivery) orderGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := uuid.Parse(r.Context().Value(ctxUserID{}).(string))
	if err != nil {
		d.logger.Error("failed parse user id from context", slog.Any("err", err),
			slog.String("handler", "orderGet"))
		handelOrderErr(w, domain.ErrUserUnauthorized)
		return
	}

	var orders ordersResponse
	orders, err = d.order.GetMany(r.Context(), userID)
	if err != nil {
		d.logger.Error("failed get user orders", slog.Any("err", err),
			slog.String("handler", "orderGet"))
		handelOrderErr(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	if err = encoder.Encode(&orders); err != nil {
		d.logger.Error("failed encode orders", slog.Any("err", err),
			slog.String("handler", "orderGet"))
		handelOrderErr(w, ErrEncodeData)
		return
	}
}

func handelOrderErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrUserUnauthorized):
		w.WriteHeader(http.StatusUnauthorized)
	case errors.Is(err, ErrDecodeBody):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, domain.ErrIncorrectData):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, domain.ErrAlreadyExist):
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, domain.ErrOrderNumberTaken):
		w.WriteHeader(http.StatusConflict)
	case errors.Is(err, domain.ErrOrderNumberInvalid):
		w.WriteHeader(http.StatusUnprocessableEntity)
	case errors.Is(err, ErrEncodeData):
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}
