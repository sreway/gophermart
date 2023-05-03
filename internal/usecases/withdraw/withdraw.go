package withdraw

import (
	"context"
	"os"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"

	"github.com/sreway/gophermart/internal/domain"
	"github.com/sreway/gophermart/internal/usecases/adapters/storage"
)

type (
	useCase struct {
		storage storage.Withdraw
		logger  *slog.Logger
	}
)

func (uc *useCase) Add(ctx context.Context, userID uuid.UUID, order domain.OrderNumber, value float64) error {
	withdraw := domain.NewWithdraw(userID, order, value)
	return uc.storage.AddWithdraw(ctx, withdraw)
}

func (uc *useCase) Get(ctx context.Context, userID uuid.UUID) ([]domain.Withdraw, error) {
	return uc.storage.GetWithdraw(ctx, userID)
}

func New(repo storage.Withdraw) *useCase {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "withdraw")}))
	return &useCase{
		storage: repo,
		logger:  log,
	}
}
