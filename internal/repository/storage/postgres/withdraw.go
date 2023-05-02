package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/sreway/gophermart/internal/domain"
)

func (r *repo) AddWithdraw(ctx context.Context, withdraw *domain.Withdraw) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err != nil {
		return err
	}

	query := "SELECT (balance >= $1) FROM balance WHERE user_id = $2"
	var balanceValid bool

	err = tx.QueryRow(ctx, query, withdraw.Sum(), withdraw.UserID()).Scan(&balanceValid)
	if err != nil {
		return err
	}
	if !balanceValid {
		return domain.NewBalanceError(withdraw.UserID().String(), domain.ErrBalanceNotEnough)
	}

	query = "UPDATE balance SET balance = balance - $1, withdrawn = withdrawn + $1 WHERE user_id = $2"
	_, err = tx.Exec(ctx, query, withdraw.Sum(), withdraw.UserID())
	if err != nil {
		return domain.NewBalanceError(withdraw.UserID().String(), err)
	}

	query = "INSERT INTO withdrawals (id, user_id, order_number, processed_at, sum) VALUES ($1, $2, $3, $4, $5)"
	_, err = tx.Exec(ctx, query, withdraw.ID(), withdraw.UserID(), withdraw.Order().Value(), time.Now(), withdraw.Sum())
	if err != nil {
		return domain.NewWithdrawError(withdraw.ID().String(), err)
	}

	return tx.Commit(ctx)
}

func (r *repo) GetWithdraw(ctx context.Context, userID uuid.UUID) ([]domain.Withdraw, error) {
	withdrawals := make([]domain.Withdraw, 0)
	query := "SELECT id, order_number, sum, processed_at from withdrawals WHERE user_id=$1 ORDER BY processed_at ASC"
	row, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var (
			id          uuid.UUID
			number      int
			sum         float64
			processedAt time.Time
		)

		if err = row.Scan(
			&id, &number, &sum, &processedAt); err != nil {
			return nil, err
		}
		withdraw, err := domain.CreateWithdraw(id, userID, number, sum, processedAt)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, *withdraw)
	}

	return withdrawals, row.Err()
}
