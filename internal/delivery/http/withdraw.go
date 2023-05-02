package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"

	"github.com/sreway/gophermart/internal/domain"
)

func (d *delivery) withdrawAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := uuid.Parse(r.Context().Value(ctxUserID{}).(string))
	if err != nil {
		d.logger.Error("failed parse user id from context", slog.Any("err", err),
			slog.String("handler", "withdrawAdd"))
		handelWithdrawErr(w, domain.ErrUserUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	withdraw := new(withdrawRequest)

	if err = decoder.Decode(withdraw); err != nil {
		d.logger.Error("failed decode request body",
			slog.String("handler", "withdrawAdd"),
			slog.Any("err", err))
		handelWithdrawErr(w, ErrDecodeBody)
		return
	}

	numberValue, err := strconv.Atoi(withdraw.Order)
	if err != nil {
		d.logger.Error("failed get order number from request body", slog.Any("err", err),
			slog.String("handler", "withdrawAdd"))
		handelWithdrawErr(w, domain.ErrIncorrectData)
		return
	}

	orderNumber, err := domain.NewOrderNumber(numberValue)
	if err != nil {
		d.logger.Error("order number incorrect", slog.Any("err", err),
			slog.String("handler", "withdrawAdd"))
		handelWithdrawErr(w, err)
		return
	}

	err = d.withdraw.Add(r.Context(), userID, *orderNumber, withdraw.Sum)
	if err != nil {
		d.logger.Error("failed add withdraw", slog.Any("err", err),
			slog.String("handler", "withdrawAdd"))
		handelWithdrawErr(w, err)
		return
	}
}

func (d *delivery) withdrawGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := uuid.Parse(r.Context().Value(ctxUserID{}).(string))
	if err != nil {
		d.logger.Error("failed parse user id from context", slog.Any("err", err),
			slog.String("handler", "withdrawGet"))
		handelOrderErr(w, domain.ErrUserUnauthorized)
		return
	}

	var withdraws withdrawsResponse
	withdraws, err = d.withdraw.Get(r.Context(), userID)
	if err != nil {
		d.logger.Error("failed get user withdraws", slog.Any("err", err),
			slog.String("handler", "withdrawGet"))
		handelOrderErr(w, err)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	if err = encoder.Encode(&withdraws); err != nil {
		d.logger.Error("failed encode withdraws", slog.Any("err", err),
			slog.String("handler", "withdrawGet"))
		handelOrderErr(w, ErrEncodeData)
		return
	}
}

func handelWithdrawErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrUserUnauthorized):
		w.WriteHeader(http.StatusUnauthorized)
	case errors.Is(err, domain.ErrNotFound):
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, ErrDecodeBody):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, domain.ErrIncorrectData):
		w.WriteHeader(http.StatusUnprocessableEntity)
	case errors.Is(err, domain.ErrOrderNumberInvalid):
		w.WriteHeader(http.StatusUnprocessableEntity)
	case errors.Is(err, domain.ErrBalanceNotEnough):
		w.WriteHeader(http.StatusPaymentRequired)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
