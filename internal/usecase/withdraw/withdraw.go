package withdraw

import (
	"context"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/internal/usecase"
)

type Withdraw struct {
	withdraw usecase.WithdrawRepo
}

func New(withdraw usecase.WithdrawRepo) *Withdraw {
	return &Withdraw{
		withdraw: withdraw,
	}
}

func (wc *Withdraw) Add(ctx context.Context, withdraw *entity.Withdraw) error {
	err := wc.withdraw.Add(ctx, withdraw)
	if err != nil {
		return err
	}

	return nil
}

func (wc *Withdraw) Get(ctx context.Context, userID uint) ([]*entity.WithdrawOrder, error) {
	withdrawals, err := wc.withdraw.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, entity.ErrWithdrawEmptyData
	}

	return withdrawals, err
}
