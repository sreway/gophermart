package order

import (
	"context"
	"errors"
	"os"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"

	"github.com/sreway/gophermart/internal/config"
	"github.com/sreway/gophermart/internal/domain"
	"github.com/sreway/gophermart/internal/usecases/adapters/accrual"
	"github.com/sreway/gophermart/internal/usecases/adapters/storage"
)

type (
	useCase struct {
		storage storage.Order
		logger  *slog.Logger
		accrual accrual.Storage
		tasks   chan *task
	}
)

func (uc *useCase) Add(ctx context.Context, userID uuid.UUID, number domain.OrderNumber) error {
	order := domain.NewOrder(userID, number, domain.OrderNew)
	err := uc.storage.AddOrder(ctx, order)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExist) {
			_, err = uc.storage.GetOrder(ctx, userID, number)
			switch err {
			case nil:
				return domain.NewOrderError(order.ID().String(), domain.ErrAlreadyExist)
			default:
				return domain.NewOrderError(order.ID().String(), domain.ErrOrderNumberTaken)
			}
		}
		uc.logger.Error("failed add order", slog.Any("err", err), slog.Int("number", number.Value()))
		return err
	}

	if len(uc.tasks) == cap(uc.tasks) {
		return ErrTaskBufferFull
	}

	uc.tasks <- NewTask(procAccrual, *order)

	uc.logger.Info("success add order",
		slog.String("id", order.ID().String()),
		slog.Int("number", number.Value()),
	)
	return nil
}

func (uc *useCase) GetMany(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	return uc.storage.GetUserOrders(ctx, userID)
}

func (uc *useCase) Get(ctx context.Context, userID uuid.UUID, number domain.OrderNumber) (*domain.Order, error) {
	return uc.storage.GetOrder(ctx, userID, number)
}

func New(config *config.Orders, repo storage.Order, accrual accrual.Storage) *useCase {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "order")}))

	tasks := make(chan *task, config.MaxTaskQueue)

	return &useCase{
		storage: repo,
		logger:  log,
		accrual: accrual,
		tasks:   tasks,
	}
}
