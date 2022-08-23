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
		ID     uint `json:"id"`
		UserID uint `json:"user_id"`
		WithdrawOrder
	}

	Withdrawals struct {
		Items         []*Withdraw `json:"items"`
		PageSize      uint        `json:"page_size"`
		NextPageToken string      `json:"next_page_token,omitempty"`
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
