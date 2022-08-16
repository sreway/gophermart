package balance

import (
	"context"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/internal/usecase"
)

type Balance struct {
	balance usecase.BalanceRepo
}

func New(balance usecase.BalanceRepo) *Balance {
	return &Balance{
		balance: balance,
	}
}

func (bc *Balance) Get(ctx context.Context, userID uint) (*entity.Balance, error) {
	balance, err := bc.balance.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
