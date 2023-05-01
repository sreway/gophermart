package domain

import (
	"strconv"
)

type (
	AccrualStatus string
	Accrual       struct {
		number OrderNumber
		status AccrualStatus
		value  float64
	}
)

const (
	AccrualRegistered AccrualStatus = "REGISTERED"
	AccrualProcessing AccrualStatus = "PROCESSING"
	AccrualInvalid    AccrualStatus = "INVALID"
	AccrualProcessed  AccrualStatus = "PROCESSED"
)

func (a *Accrual) Number() *OrderNumber {
	return &a.number
}

func (a *Accrual) Status() AccrualStatus {
	return a.status
}

func (a *Accrual) Value() float64 {
	return a.value
}

func NewAccrualStatus(value string) (*AccrualStatus, error) {
	status := AccrualStatus(value)
	switch status {
	case AccrualRegistered, AccrualProcessing, AccrualInvalid, AccrualProcessed:
		return &status, nil
	default:
		return nil, ErrAccrualStatusInvalid
	}
}

func NewAccrual(number string, status string, value float64) (*Accrual, error) {
	numberInt, err := strconv.Atoi(number)
	if err != nil {
		return nil, err
	}

	orderNumber, err := NewOrderNumber(numberInt)
	if err != nil {
		return nil, err
	}

	accrualStatus, err := NewAccrualStatus(status)
	if err != nil {
		return nil, err
	}

	return &Accrual{
		number: *orderNumber,
		status: *accrualStatus,
		value:  value,
	}, nil
}
