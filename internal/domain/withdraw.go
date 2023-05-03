package domain

import (
	"time"

	"github.com/google/uuid"
)

type Withdraw struct {
	id          uuid.UUID
	userID      uuid.UUID
	order       OrderNumber
	sum         float64
	processedAt time.Time
}

func (w *Withdraw) ID() uuid.UUID {
	return w.id
}

func (w *Withdraw) UserID() uuid.UUID {
	return w.userID
}

func (w *Withdraw) Order() *OrderNumber {
	return &w.order
}

func (w *Withdraw) Sum() float64 {
	return w.sum
}

func (w *Withdraw) ProcessedAt() time.Time {
	return w.processedAt
}

func NewWithdraw(userID uuid.UUID, number OrderNumber, sum float64) *Withdraw {
	id := uuid.New()
	return &Withdraw{
		id:     id,
		userID: userID,
		order:  number,
		sum:    sum,
	}
}

func CreateWithdraw(id, userID uuid.UUID, number int, sum float64, processedAt time.Time) (*Withdraw, error) {
	orderNumber, err := NewOrderNumber(number)
	if err != nil {
		return nil, NewWithdrawError(id.String(), err)
	}

	return &Withdraw{
		id:          id,
		userID:      userID,
		order:       *orderNumber,
		sum:         sum,
		processedAt: processedAt,
	}, nil
}
