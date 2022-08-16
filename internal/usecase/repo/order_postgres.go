package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/pkg/postgres"
)

type OrderRepo struct {
	*postgres.Postgres
}

func NewOrderRepo(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{pg}
}

func (or *OrderRepo) Add(ctx context.Context, order *entity.Order) error {
	tx, err := or.Postgres.Pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err != nil {
		return err
	}

	query := "INSERT INTO orders (number, user_id, status, accrual)  VALUES ($1, $2, $3, $4)"
	_, err = tx.Exec(ctx, query, order.Number, order.UserID, order.Status, order.Accrual)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return entity.ErrOrderAlreadyExist
			default:
				return err
			}
		}
		return err
	}

	_, err = tx.Exec(ctx, "SELECT pg_notify('new_order', $1)", order.Number)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (or *OrderRepo) Get(ctx context.Context, number string) (*entity.Order, error) {
	order := entity.Order{}
	query := "SELECT id, number, user_id, status, accrual, uploaded_at FROM orders WHERE number=$1"
	err := or.Postgres.Pool.QueryRow(ctx, query, number).Scan(
		&order.ID, &order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (or *OrderRepo) GetAll(ctx context.Context, userID uint) ([]*entity.Order, error) {
	orders := make([]*entity.Order, 0)
	query := "SELECT id, number, user_id, status, accrual, uploaded_at " +
		"FROM orders WHERE user_id=$1 ORDER BY uploaded_at ASC"

	row, err := or.Postgres.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		order := entity.Order{}
		if err = row.Scan(
			&order.ID, &order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (or *OrderRepo) UpdateStatus(ctx context.Context, number string, status entity.OrderStatus, accrual float64) error {
	tx, err := or.Postgres.Pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err != nil {
		return err
	}

	var userID uint
	updateOrderQuery := "UPDATE orders SET accrual = $1, status = $2 WHERE number = $3 RETURNING user_id"
	err = tx.QueryRow(ctx, updateOrderQuery, accrual, status, number).Scan(&userID)
	if err != nil {
		return err
	}
	switch status {
	case entity.OrderStatusInvalid, entity.OrdertStatusProcessed:
		updateUserBalance := "UPDATE balance SET balance=balance+$1 WHERE user_id = $2"
		_, err = tx.Exec(ctx, updateUserBalance, accrual, userID)
		if err != nil {
			return err
		}
		return tx.Commit(ctx)
	default:
		return tx.Commit(ctx)
	}
}
