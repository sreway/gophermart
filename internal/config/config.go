package config

import (
	"time"
	// to do viper
	"github.com/caarlos0/env/v8"
)

type (
	config struct {
		http    *HTTP
		storage *Postgres
		accrual *Accrual
		orders  *Orders
	}

	HTTP struct {
		Address string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8080"`
		Auth    struct {
			TokenTTL   time.Duration `env:"AUTH_TOKEN_TTL" envDefault:"30m"`
			Key        string        `env:"AUTH_KEY" envDefault:"supersecret"`
			CookieName string        `env:"AUTH_COOKIE_NAME" envDefault:"jwt"`
		}
	}

	Postgres struct {
		DSN        string `env:"DATABASE_URI"`
		MigrateURL string `env:"MIGRATE_URL" envDefault:"file://migrations/postgres"`
	}

	Orders struct {
		TaskInterval time.Duration `env:"ORDERS_TASK_INTERVAL" envDefault:"5s"`
		MaxTaskQueue int           `env:"ORDERS_MAX_TASK_QUEUE" envDefault:"100"`
	}

	Accrual struct {
		Address string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://127.0.0.1:8082"`
	}
)

func (c *config) HTTP() *HTTP {
	return c.http
}

func (c *config) Storage() *Postgres {
	return c.storage
}

func (c *config) Accrual() *Accrual {
	return c.accrual
}

func (c *config) Orders() *Orders {
	return c.orders
}

func New() (*config, error) {
	cfg := new(config)
	cfg.http = new(HTTP)
	cfg.storage = new(Postgres)
	cfg.accrual = new(Accrual)
	cfg.orders = new(Orders)

	if err := env.Parse(cfg.http); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg.storage); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg.accrual); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg.orders); err != nil {
		return nil, err
	}

	return cfg, nil
}
