package entity

import (
	"encoding/json"
)

type (
	AccrualStatus string
	Accrual       struct {
		OrderNumber string        `json:"order"`
		Status      AccrualStatus `json:"status"`
		Value       float64       `json:"accrual"`
	}
)

const (
	AccrualStatusStatusRegistered AccrualStatus = "REGISTERED"
	AccrualStatusProcessing       AccrualStatus = "PROCESSING"
	AccrualStatusInvalid          AccrualStatus = "INVALID"
	AccrualStatusProcessed        AccrualStatus = "PROCESSED"
)

func AllowedAccrualStatus() []AccrualStatus {
	return []AccrualStatus{
		AccrualStatusStatusRegistered, AccrualStatusProcessing, AccrualStatusInvalid,
		AccrualStatusProcessed,
	}
}

func (as AccrualStatus) String() string {
	for _, v := range AllowedAccrualStatus() {
		if v == as {
			return string(v)
		}
	}
	return "Invalid status name"
}

func (as AccrualStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(as.String())
}
