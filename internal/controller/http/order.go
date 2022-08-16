package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/pkg/logger"

	"github.com/sreway/gophermart/internal/usecase"
)

type orderRoutes struct {
	order usecase.Order
}

func newOrderRoutes(o usecase.Order) *orderRoutes {
	return &orderRoutes{o}
}

func (or *orderRoutes) orderAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	userID := uint(r.Context().Value(userIDKey).(float64))
	userLogin := r.Context().Value(userLoginKey).(string)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		HandelErrOrder(w, entity.ErrOrderIncorrectData)
		return
	}

	order := new(entity.Order)
	order.Number = string(body)
	order.Status = entity.OrderStatusNew
	order.UserID = userID

	err = or.order.Add(r.Context(), order)

	if err != nil {
		HandelErrOrder(w, entity.NewOrderError(userLogin, order.Number, err))
		return
	}
	logger.Infof("Order success add %s", order.Number)
	w.WriteHeader(http.StatusAccepted)
}

func (or *orderRoutes) orderGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := uint(r.Context().Value(userIDKey).(float64))
	userLogin := r.Context().Value(userLoginKey).(string)
	orders, err := or.order.Get(r.Context(), userID)
	if err != nil {
		HandelErrOrder(w, entity.NewOrderError(userLogin, "all", err))
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	if err = encoder.Encode(&orders); err != nil {
		logger.Error(entity.NewOrderError(userLogin, "all", err))
		return
	}
}

func HandelErrOrder(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.ErrOrderNumberTaken):
		w.WriteHeader(http.StatusConflict)
	case errors.Is(err, entity.ErrOrderAlreadyExist):
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, entity.ErrOrderIncorrectData):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, entity.ErrOrderIncorrectNumber):
		w.WriteHeader(http.StatusUnprocessableEntity)
	case errors.Is(err, entity.ErrUserUnauthorized):
		w.WriteHeader(http.StatusUnauthorized)
	case errors.Is(err, entity.ErrOrderEmptyData):
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	logger.Error(err)
}
