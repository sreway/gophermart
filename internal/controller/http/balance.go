package http

import (
	"encoding/json"
	"net/http"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/internal/usecase"
	"github.com/sreway/gophermart/pkg/logger"
)

type balanceRoutes struct {
	balance usecase.Balance
}

func newBalanceRoutes(b usecase.Balance) *balanceRoutes {
	return &balanceRoutes{b}
}

func (br *balanceRoutes) balanceGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := uint(r.Context().Value(userIDKey).(float64))
	userLogin := r.Context().Value(userLoginKey).(string)
	orders, err := br.balance.Get(r.Context(), userID)
	if err != nil {
		HandelErrOBalance(w, entity.NewAppError(userLogin, err))
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	if err = encoder.Encode(&orders); err != nil {
		HandelErrOBalance(w, entity.NewAppError(userLogin, err))
		return
	}
}

func HandelErrOBalance(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	logger.Error(err)
}
