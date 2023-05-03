package user

import (
	"context"
	"os"

	"golang.org/x/exp/slog"

	"github.com/sreway/gophermart/internal/domain"
	"github.com/sreway/gophermart/internal/usecases/adapters/storage"
)

type (
	useCase struct {
		storage storage.User
		logger  *slog.Logger
	}
)

func (uc *useCase) Register(ctx context.Context, login, password string) (*domain.User, error) {
	hashPassword, err := createHashPassword(password)
	if err != nil {
		return nil, domain.NewUserError("", domain.ErrIncorrectData)
	}

	user := domain.NewUser(login, string(hashPassword))

	err = uc.storage.AddUser(ctx, user)
	if err != nil {
		return nil, err
	}

	uc.logger.Info("success user register",
		slog.String("id", user.ID().String()),
		slog.String("login", user.Login()),
	)

	return user, nil
}

func (uc *useCase) Login(ctx context.Context, login, password string) (*domain.User, error) {
	user, err := uc.storage.GetUser(ctx, login)
	if err != nil {
		uc.logger.Error("failed get user", slog.Any("err", err), slog.String("login", login))
		return nil, err
	}

	if !compareHashPassword(user.HashPassword(), password) {
		return nil, domain.NewUserError(user.ID().String(), domain.ErrUserUnauthorized)
	}

	return user, nil
}

func New(repo storage.User) *useCase {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "user")}))
	return &useCase{
		storage: repo,
		logger:  log,
	}
}
