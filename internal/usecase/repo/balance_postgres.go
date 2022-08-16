package repo

import (
	"context"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/pkg/postgres"
)

type BalanceRepo struct {
	*postgres.Postgres
}

func NewBalanceRepo(pg *postgres.Postgres) *BalanceRepo {
	return &BalanceRepo{pg}
}

func (br *BalanceRepo) Get(ctx context.Context, userID uint) (*entity.Balance, error) {
	balance := entity.Balance{
		UserID: userID,
	}

	query := "SELECT id, balance, withdrawn FROM balance WHERE user_id=$1"
	err := br.Postgres.Pool.QueryRow(ctx, query, userID).Scan(&balance.ID, &balance.Value, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}
