package postgres

import (
	"context"
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/exp/slog"

	"github.com/sreway/gophermart/internal/config"
)

type repo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func (r *repo) migrate(migrateURL string) error {
	m, err := migrate.New(migrateURL, r.pool.Config().ConnConfig.ConnString())
	defer func() {
		_, _ = m.Close()
	}()

	if err != nil {
		return err
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		r.logger.Info("no change", slog.String("func", "migrate"))
		return nil
	}

	if err != nil {
		r.logger.Error("failed apply migrations",
			slog.String("func", "migrate"),
			slog.Any("err", err),
		)
		return err
	}

	return nil
}

func New(ctx context.Context, config *config.Postgres) (*repo, error) {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("repository", "postgres")}))

	poolConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	r := &repo{
		pool:   pool,
		logger: log,
	}

	if len(config.MigrateURL) == 0 {
		return r, nil
	}

	err = r.migrate(config.MigrateURL)
	if err != nil {
		return nil, err
	}

	return r, nil
}
