package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/sreway/gophermart/internal/domain"
)

func (r *repo) AddUser(ctx context.Context, user *domain.User) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err != nil {
		return err
	}

	query := "INSERT INTO users (id, login, hash_password) VALUES ($1, $2, $3)"
	_, err = tx.Exec(ctx, query, user.ID(), user.Login(), user.HashPassword())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return domain.NewUserError(user.ID().String(), domain.ErrAlreadyExist)
			default:
				return err
			}
		}
		return err
	}

	balance := domain.NewBalance(user.ID(), 0, 0)

	query = "INSERT INTO balance (id, user_id) VALUES ($1, $2)"
	if _, err = tx.Exec(ctx, query, balance.ID(), balance.UserID()); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *repo) GetUser(ctx context.Context, login string) (*domain.User, error) {
	var (
		id           uuid.UUID
		hashPassword string
	)
	query := "SELECT id, hash_password from users where login = $1"
	err := r.pool.QueryRow(ctx, query, login).Scan(&id, &hashPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewUserError("", domain.ErrNotFound)
		}
		return nil, err
	}

	return domain.CreateUser(id, login, hashPassword), nil
}
