package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/sreway/gophermart/internal/domain"
)

type (
	User interface {
		AddUser(ctx context.Context, user *domain.User) error
		GetUser(ctx context.Context, login string) (*domain.User, error)
	}

	Balance interface {
		GetBalance(ctx context.Context, userID uuid.UUID) (*domain.Balance, error)
	}

	Order interface {
		AddOrder(ctx context.Context, order *domain.Order) error
		GetOrder(ctx context.Context, userID uuid.UUID, number domain.OrderNumber) (*domain.Order, error)
		GetUserOrders(ctx context.Context, userID uuid.UUID) ([]domain.Order, error)
		UpdateOrderStatus(ctx context.Context, number domain.OrderNumber, status domain.OrderStatus,
			accrual float64) (*domain.Order, error)
		RefillBalance(ctx context.Context, userID uuid.UUID, value float64) error
	}

	Withdraw interface {
		AddWithdraw(ctx context.Context, withdraw *domain.Withdraw) error
		GetWithdraw(ctx context.Context, userID uuid.UUID) ([]domain.Withdraw, error)
	}
)
