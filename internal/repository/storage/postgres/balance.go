package postgres

import (
	"context"

	"github.com/google/uuid"

	"github.com/sreway/gophermart/internal/domain"
)

func (r *repo) GetBalance(ctx context.Context, userID uuid.UUID) (*domain.Balance, error) {
	query := "SELECT id, balance, withdrawn FROM balance WHERE user_id=$1"
	var (
		id        uuid.UUID
		value     float64
		withdrawn float64
	)
	err := r.pool.QueryRow(ctx, query, userID).Scan(&id, &value, &withdrawn)
	if err != nil {
		return nil, domain.NewBalanceError(userID.String(), err)
	}

	return domain.NewBalance(userID, value, withdrawn), nil
}
