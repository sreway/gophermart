package app

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"

	"github.com/sreway/gophermart/config"
	"github.com/sreway/gophermart/pkg/logger"

	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(cfg *config.Config) {
	m, err := migrate.New(cfg.Postgres.MigrateURL, cfg.Postgres.DSN)
	defer func() {
		_, _ = m.Close()
	}()

	if err != nil {
		logger.Panicf("Migrate: %s", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Panicf("Migrate: up error: %s", err)
	}
	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("Migrate: no change")
		return
	}

	logger.Info("Migrate: up success")
}
