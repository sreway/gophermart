package usecase

import (
	"context"

	"github.com/sreway/gophermart/internal/entity"
)

type (
	User interface {
		Register(ctx context.Context, user *entity.User) error
		Login(ctx context.Context, user *entity.User) (*entity.User, error)
	}

	UserRepo interface {
		Add(ctx context.Context, user *entity.User) error
		Get(ctx context.Context, login string) (*entity.User, error)
	}

	Order interface {
		Add(ctx context.Context, order *entity.Order) error
		Get(ctx context.Context, userID uint) ([]*entity.Order, error)
	}

	OrderRepo interface {
		Add(ctx context.Context, order *entity.Order) error
		Get(ctx context.Context, number string) (*entity.Order, error)
		GetAll(ctx context.Context, userID uint) ([]*entity.Order, error)
		UpdateStatus(ctx context.Context, number string, status entity.OrderStatus, accrual float64) error
	}

	Balance interface {
		Get(ctx context.Context, userID uint) (*entity.Balance, error)
	}

	BalanceRepo interface {
		Get(ctx context.Context, userID uint) (*entity.Balance, error)
	}

	Withdraw interface {
		Add(ctx context.Context, withdraw *entity.Withdraw) error
		Get(ctx context.Context, userID uint) ([]*entity.WithdrawOrder, error)
	}

	WithdrawRepo interface {
		Add(ctx context.Context, withdraw *entity.Withdraw) error
		GetAll(ctx context.Context, userID uint) ([]*entity.WithdrawOrder, error)
	}

	Accrual interface {
		Get(ctx context.Context, number string) (*entity.Accrual, error)
		UpdateOrderStatus(ctx context.Context, accrual *entity.Accrual) error
	}

	AccrualRepo interface {
		Get(ctx context.Context, number string) (*entity.Accrual, error)
	}

	QueueRepo interface {
		Add(ctx context.Context, number string) error
		Read(ctx context.Context) (string, error)
		Commit(ctx context.Context, msg *entity.QueueMsg) error
	}
)
