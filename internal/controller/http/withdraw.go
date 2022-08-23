package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/internal/usecase"
	"github.com/sreway/gophermart/pkg/logger"
)

type withdrawRoutes struct {
	withdraw usecase.Withdraw
}

func newWithdrawRoutes(wc usecase.Withdraw) *withdrawRoutes {
	return &withdrawRoutes{wc}
}

func (wr *withdrawRoutes) withdrawAdd(w http.ResponseWriter, r *http.Request) {
	var wo entity.WithdrawOrder
	w.Header().Set("Content-Type", "application/json")
	userID := uint(r.Context().Value(userIDKey).(float64))
	userLogin := r.Context().Value(userLoginKey).(string)

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&wo); err != nil {
		HandelErrOBalance(w, entity.NewErrWithdraw(userLogin, "-", err))
		return
	}

	withdraw := new(entity.Withdraw)

	withdraw.UserID = userID
	withdraw.OrderNumber = wo.OrderNumber
	withdraw.Sum = wo.Sum

	err := wr.withdraw.Add(r.Context(), withdraw)
	if err != nil {
		HandelErrWithdraw(w, entity.NewErrWithdraw(userLogin, wo.OrderNumber, err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (wr *withdrawRoutes) withdrawGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := uint(r.Context().Value(userIDKey).(float64))
	userLogin := r.Context().Value(userLoginKey).(string)
	withdrawals, err := wr.withdraw.Get(r.Context(), userID)
	if err != nil {
		HandelErrWithdraw(w, entity.NewErrWithdraw(userLogin, "all", err))
		return
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	if err = encoder.Encode(&withdrawals); err != nil {
		logger.Error(entity.NewErrWithdraw(userLogin, "all", err))
		return
	}
}

func (wr *withdrawRoutes) withdrawGetPagination(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := uint(r.Context().Value(userIDKey).(float64))
	userLogin := r.Context().Value(userLoginKey).(string)

	nextPageToken := r.URL.Query().Get("next_page_token")
	pageSize := r.URL.Query().Get("page_size")

	withdrawals, err := wr.withdraw.GetPagination(r.Context(), userID, nextPageToken, pageSize)
	if err != nil {
		HandelErrWithdraw(w, entity.NewErrWithdraw(userLogin, "all", err))
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	if err = encoder.Encode(&withdrawals); err != nil {
		logger.Error(entity.NewErrWithdraw(userLogin, "all", err))
		return
	}
}

func HandelErrWithdraw(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.ErrOrderNotFound):
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, entity.ErrOrderIncorrectNumber):
		w.WriteHeader(http.StatusUnprocessableEntity)
	case errors.Is(err, entity.ErrBalanceNotEnough):
		w.WriteHeader(http.StatusPaymentRequired)
	case errors.Is(err, entity.ErrWithdrawEmptyData):
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	logger.Error(err)
}
