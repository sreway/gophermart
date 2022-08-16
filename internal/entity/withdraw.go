package entity

import (
	"encoding/json"
	"time"
)

type (
	WithdrawOrder struct {
		OrderNumber string    `json:"order"`
		Sum         float64   `json:"sum"`
		ProcessedAt time.Time `json:"processed_at"`
	}

	Withdraw struct {
		ID     string `json:"id"`
		UserID uint   `json:"user_id"`
		WithdrawOrder
	}
)

func (wo WithdrawOrder) MarshalJSON() ([]byte, error) {
	type WithdrawOrderAlias WithdrawOrder
	aliasValue := struct {
		WithdrawOrderAlias
		Sum        CustomFloat `json:"sum"`
		UploadedAt string      `json:"processed_at"`
	}{
		WithdrawOrderAlias: WithdrawOrderAlias(wo),
		Sum:                CustomFloat{Value: wo.Sum, Digits: 2},
		UploadedAt:         wo.ProcessedAt.Format(time.RFC3339),
	}
	return json.Marshal(aliasValue)
}
