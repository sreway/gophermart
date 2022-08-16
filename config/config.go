package config

import (
	"time"

	"github.com/caarlos0/env/v6"
)

var (
	DefaultAccrualAddress    string
	DefaultHTTPAddress       string
	DefaultHTTPCompressLevel = 5
	DefaultHTTPCompressTypes = []string{
		"text/html",
		"text/plain",
		"application/json",
	}
	DefaultPostgresDSN           string
	DefaultPostgresMigrateURL    = "file://migrations/"
	DefaultPostgresListenChannel = "new_order"
	DefaultJWTKey                = "some_secret" // not set as default, need move env
	DefaultJWTLiveTime           = 30 * time.Minute
)

type (
	Config struct {
		Server   ServerConfig
		Postgres PostgresConfig
		Accrual  AccrualConfig
	}

	ServerConfig struct {
		Auth AuthConfig
		HTTP HTTPConfig
	}

	AuthConfig struct {
		JWT JWTConfig
	}

	JWTConfig struct {
		TokenTTL time.Duration
		Key      string
	}

	HTTPConfig struct {
		Address       string `env:"RUN_ADDRESS"`
		CompressLevel int
		CompressTypes []string
	}

	PostgresConfig struct {
		DSN           string `env:"DATABASE_URI"`
		MigrateURL    string
		ListenChannel string `env:"DATABASE_LISTEN_CHANNEL"`
	}

	AccrualConfig struct {
		Address string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			HTTP: HTTPConfig{
				Address:       DefaultHTTPAddress,
				CompressLevel: DefaultHTTPCompressLevel,
				CompressTypes: DefaultHTTPCompressTypes,
			},
			Auth: AuthConfig{
				JWTConfig{
					Key:      DefaultJWTKey,
					TokenTTL: DefaultJWTLiveTime,
				},
			},
		},
		Postgres: PostgresConfig{
			DSN:           DefaultPostgresDSN,
			MigrateURL:    DefaultPostgresMigrateURL,
			ListenChannel: DefaultPostgresListenChannel,
		},
		Accrual: AccrualConfig{
			Address: DefaultAccrualAddress,
		},
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
