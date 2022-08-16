package accrual

import (
	"context"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/internal/usecase"
)

type Accrual struct {
	accrual usecase.AccrualRepo
	order   usecase.OrderRepo
}

func New(accrual usecase.AccrualRepo, order usecase.OrderRepo) *Accrual {
	return &Accrual{
		accrual: accrual,
		order:   order,
	}
}

func (ac *Accrual) Get(ctx context.Context, number string) (*entity.Accrual, error) {
	accrual, err := ac.accrual.Get(ctx, number)
	if err != nil {
		return nil, err
	}

	return accrual, nil
}

func (ac *Accrual) UpdateOrderStatus(ctx context.Context, accrual *entity.Accrual) error {
	err := ac.order.UpdateStatus(ctx, accrual.OrderNumber, entity.OrderStatus(accrual.Status), accrual.Value)
	if err != nil {
		return err
	}

	return nil
}
