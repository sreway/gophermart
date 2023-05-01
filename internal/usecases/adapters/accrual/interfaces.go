package accrual

import (
	"context"

	"github.com/sreway/gophermart/internal/domain"
)

type Storage interface {
	Get(ctx context.Context, number domain.OrderNumber) (*domain.Accrual, error)
}
