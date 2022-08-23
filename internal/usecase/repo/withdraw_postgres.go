package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/pkg/postgres"
)

type WithdrawRepo struct {
	*postgres.Postgres
}

func NewWithdraw(pg *postgres.Postgres) *WithdrawRepo {
	return &WithdrawRepo{pg}
}

func (wr *WithdrawRepo) Add(ctx context.Context, withdraw *entity.Withdraw) error {
	tx, err := wr.Postgres.Pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err != nil {
		return err
	}
	checkBalanceQuery := "SELECT (balance >= $1) FROM balance WHERE user_id = $2"
	balanceValid := false
	err = tx.QueryRow(ctx, checkBalanceQuery, withdraw.Sum, withdraw.UserID).Scan(&balanceValid)
	if err != nil {
		return err
	}
	if !balanceValid {
		return entity.ErrBalanceNotEnough
	}

	updateBalanceQuery := "UPDATE balance SET balance = balance - $1, withdrawn = withdrawn + $1 WHERE user_id = $2"
	_, err = tx.Exec(ctx, updateBalanceQuery, withdraw.Sum, withdraw.UserID)
	if err != nil {
		return err
	}

	addWithdrawQuery := "INSERT INTO withdrawals (user_id, order_number, processed_at, sum) VALUES ($1, $2, $3, $4)"
	if _, err = tx.Exec(ctx, addWithdrawQuery, withdraw.UserID, withdraw.OrderNumber, time.Now(), withdraw.Sum); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (wr *WithdrawRepo) Get(ctx context.Context, userID uint) ([]*entity.WithdrawOrder, error) {
	withdrawals := make([]*entity.WithdrawOrder, 0)
	query := "select order_number, sum, processed_at from withdrawals WHERE user_id=$1 ORDER BY processed_at ASC"
	row, err := wr.Postgres.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		withdrawOrder := entity.WithdrawOrder{}
		if err = row.Scan(&withdrawOrder.OrderNumber, &withdrawOrder.Sum, &withdrawOrder.ProcessedAt); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, &withdrawOrder)
	}

	return withdrawals, nil
}

func (wr *WithdrawRepo) GetPagination(ctx context.Context, userID, startAt, limit uint) (*entity.Withdrawals, error) {
	withdrawals := entity.Withdrawals{}
	query := `select id, order_number, sum, processed_at from withdrawals WHERE id > $2 AND user_id=$1
				ORDER BY processed_at ASC, id ASC LIMIT $3;`
	row, err := wr.Postgres.Pool.Query(ctx, query, userID, startAt, limit)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		w := entity.Withdraw{}

		if err = row.Scan(&w.ID, &w.OrderNumber, &w.Sum, &w.ProcessedAt); err != nil {
			return nil, err
		}

		withdrawals.Items = append(withdrawals.Items, &w)
	}

	withdrawals.PageSize = limit
	return &withdrawals, nil
}
