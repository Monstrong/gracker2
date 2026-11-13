package torrentconfig

import (
	"os"

	"github.com/monstrong/gracker2/pkg/config"
)

type Config struct {
	config.Config `mapstructure:",squash"`
}

const postgresPasswordKey = "POSTGRES_PASSWORD"

func Load() (*Config, error) {
	cfg := &Config{}

	if err := config.LoadInto(cfg); err != nil {
		return nil, err
	}

	cfg.Postgres.Password = os.Getenv(postgresPasswordKey)

	if err := config.Validate(&cfg.Config); err != nil {
		return nil, err
	}

	return cfg, nil
}