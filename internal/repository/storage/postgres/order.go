package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/sreway/gophermart/internal/domain"
)

func (r *repo) AddOrder(ctx context.Context, order *domain.Order) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err != nil {
		return err
	}

	query := "INSERT INTO orders (id, number, user_id, status)  VALUES ($1, $2, $3, $4)"
	_, err = tx.Exec(ctx, query, order.ID(), order.Number().Value(), order.UserID(), order.Status())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return domain.NewOrderError(order.ID().String(), domain.ErrAlreadyExist)
			default:
				return err
			}
		}
		return err
	}
	return tx.Commit(ctx)
}

func (r *repo) GetOrder(ctx context.Context, userID uuid.UUID, number domain.OrderNumber) (*domain.Order, error) {
	var (
		id         uuid.UUID
		status     string
		accrual    float64
		uploadedAt time.Time
	)
	query := "SELECT id, status, accrual, uploaded_at FROM orders WHERE number=$1 and user_id=$2"
	err := r.pool.QueryRow(ctx, query, number.Value(), userID).Scan(&id, &status, &accrual, &uploadedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return domain.CreateOrder(id, userID, number.Value(), status, accrual, uploadedAt)
}

func (r *repo) UpdateOrderStatus(ctx context.Context, number domain.OrderNumber,
	status domain.OrderStatus, accrual float64,
) (*domain.Order, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err != nil {
		return nil, err
	}
	var (
		id         uuid.UUID
		userID     uuid.UUID
		uploadedAt time.Time
	)
	query := "UPDATE orders SET accrual = $1, status = $2 WHERE number = $3 RETURNING id, user_id, uploaded_at"
	err = tx.QueryRow(ctx, query, accrual, status, number.Value()).Scan(&id, &userID, &uploadedAt)
	if err != nil {
		return nil, err
	}

	order, err := domain.CreateOrder(id, userID, number.Value(), status.String(), accrual, uploadedAt)
	if err != nil {
		return nil, err
	}

	return order, tx.Commit(ctx)
}

func (r *repo) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	orders := make([]domain.Order, 0)
	query := "SELECT id, number, status, accrual, uploaded_at " +
		"FROM orders WHERE user_id=$1 ORDER BY uploaded_at ASC"

	row, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var (
			id         uuid.UUID
			number     int
			status     string
			accrual    float64
			uploadedAt time.Time
		)

		if err = row.Scan(
			&id, &number, &status, &accrual, &uploadedAt); err != nil {
			return nil, err
		}
		order, err := domain.CreateOrder(id, userID, number, status, accrual, uploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}

	return orders, row.Err()
}

func (r *repo) RefillBalance(ctx context.Context, userID uuid.UUID, value float64) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err != nil {
		return err
	}

	query := "UPDATE balance SET balance=balance+$1 WHERE user_id = $2"
	_, err = tx.Exec(ctx, query, value, userID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
