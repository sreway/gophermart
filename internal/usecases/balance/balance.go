package balance

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
		storage storage.Balance
		logger  *slog.Logger
	}
)

func (uc *useCase) Get(ctx context.Context, userID uuid.UUID) (*domain.Balance, error) {
	return uc.storage.GetBalance(ctx, userID)
}

func New(repo storage.Balance) *useCase {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "balance")}))
	return &useCase{
		storage: repo,
		logger:  log,
	}
}
