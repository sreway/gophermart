package http

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/sreway/gophermart/internal/domain"
)

type (
	balanceResponse struct {
		*domain.Balance
	}
	ordersResponse []domain.Order
	customFloat    struct {
		Value  float64
		Digits int
	}
	withdrawsResponse []domain.Withdraw
)

func (or ordersResponse) MarshalJSON() ([]byte, error) {
	type OrderAlias struct {
		Number     string    `json:"number"`
		UserID     uuid.UUID `json:"user_id"`
		Status     string    `json:"status"`
		Accrual    float64   `json:"accrual,omitempty"`
		UploadedAt string    `json:"uploaded_at"`
	}

	orders := make([]OrderAlias, len(or))
	for idx, item := range or {
		aliasValue := OrderAlias{}
		aliasValue.UserID = item.UserID()
		aliasValue.Number = strconv.Itoa(item.Number().Value())
		aliasValue.Status = item.Status().String()
		aliasValue.Accrual = item.Accrual()
		aliasValue.UploadedAt = item.UploadedAt().Format(time.RFC3339)
		orders[idx] = aliasValue
	}
	return json.Marshal(orders)
}

func (cf customFloat) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.*f", cf.Digits, cf.Value)
	return []byte(s), nil
}

func (br balanceResponse) MarshalJSON() ([]byte, error) {
	type (
		BalanceAlias struct {
			Value     customFloat `json:"current"`
			Withdrawn customFloat `json:"withdrawn"`
		}
	)

	balance := BalanceAlias{}
	balance.Withdrawn = customFloat{Value: br.Withdrawn(), Digits: 2}
	balance.Value = customFloat{Value: br.Value(), Digits: 2}

	return json.Marshal(balance)
}

func (wr withdrawsResponse) MarshalJSON() ([]byte, error) {
	type WithdrawAlias struct {
		Order       string      `json:"order"`
		Sum         customFloat `json:"sum"`
		ProcessedAt string      `json:"processed_at"`
	}

	withdraws := make([]WithdrawAlias, len(wr))
	for idx, item := range wr {
		aliasValue := WithdrawAlias{}
		aliasValue.Order = strconv.Itoa(item.Order().Value())
		aliasValue.Sum = customFloat{Value: item.Sum(), Digits: 2}
		aliasValue.ProcessedAt = item.ProcessedAt().Format(time.RFC3339)
		withdraws[idx] = aliasValue
	}
	return json.Marshal(withdraws)
}
