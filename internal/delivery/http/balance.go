package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"

	"github.com/sreway/gophermart/internal/domain"
)

func (d *delivery) balanceGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := uuid.Parse(r.Context().Value(ctxUserID).(string))
	if err != nil {
		d.logger.Error("failed parse user id from context", slog.Any("err", err),
			slog.String("handler", "balanceGet"))
		handelBalanceErr(w, domain.ErrUserUnauthorized)
		return
	}

	var response balanceResponse
	response.Balance, err = d.balance.Get(r.Context(), userID)
	if err != nil {
		d.logger.Error("failed get user balance", slog.Any("err", err),
			slog.String("handler", "balanceGet"))
		handelBalanceErr(w, err)
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	if err = encoder.Encode(&response); err != nil {
		d.logger.Error("failed encode user balance", slog.Any("err", err),
			slog.String("handler", "balanceGet"))
		handelBalanceErr(w, ErrEncodeData)
		return
	}
}

func handelBalanceErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrUserUnauthorized):
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
