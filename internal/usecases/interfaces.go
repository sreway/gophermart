package usecases

import (
	"context"

	"github.com/google/uuid"

	"github.com/sreway/gophermart/internal/domain"
)

type (
	User interface {
		Register(ctx context.Context, login, password string) (*domain.User, error)
		Login(ctx context.Context, login, password string) (*domain.User, error)
	}

	Order interface {
		Add(ctx context.Context, userID uuid.UUID, number domain.OrderNumber) error
		GetMany(ctx context.Context, userID uuid.UUID) ([]domain.Order, error)
		Get(ctx context.Context, userID uuid.UUID, number domain.OrderNumber) (*domain.Order, error)
	}

	Balance interface {
		Get(ctx context.Context, userID uuid.UUID) (*domain.Balance, error)
	}

	Withdraw interface {
		Add(ctx context.Context, userID uuid.UUID, order domain.OrderNumber, value float64) error
		Get(ctx context.Context, userID uuid.UUID) ([]domain.Withdraw, error)
	}
)
