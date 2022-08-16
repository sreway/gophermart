package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/pkg/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (ur *UserRepo) Add(ctx context.Context, user *entity.User) error {
	tx, err := ur.Postgres.Pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	if err != nil {
		return err
	}

	userID, err := createUser(ctx, tx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return entity.ErrUserAlreadyExist
			default:
				return err
			}
		}
		return err
	}

	err = createBalance(ctx, tx, userID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepo) Get(ctx context.Context, login string) (*entity.User, error) {
	user := entity.User{}
	query := "SELECT id, login, password_hash FROM users WHERE login=$1"
	err := ur.Postgres.Pool.QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}

	return &user, err
}

func createUser(ctx context.Context, tx pgx.Tx, user *entity.User) (uint, error) {
	query := "INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id"
	err := tx.QueryRow(ctx, query, user.Login, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func createBalance(ctx context.Context, tx pgx.Tx, userID uint) error {
	query := "INSERT INTO balance (user_id) VALUES ($1)"
	if _, err := tx.Exec(ctx, query, userID); err != nil {
		return err
	}
	return nil
}
