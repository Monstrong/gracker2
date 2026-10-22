package authconfig

import (
	"fmt"
	"os"
	"time"

	"github.com/monstrong/gracker2/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	config.Config

	JWT    JWT    `mapstructure:"jwt" yaml:"jwt"`
	Hasher Hasher `mapstructure:"hasher" yaml:"hasher"`
}

type JWT struct {
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl" yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl" yaml:"refresh_token_ttl"`
	SigningKey      string
}

type Hasher struct {
	Cost int `mapstructure:"cost" yaml:"cost"`
}

const secretKey = "JWT_SECRET"
const postgresPasswordKey = "POSTGRES_PASSWORD"

func Load() (*Config, error) {
	cfg := &Config{}

	if err := config.LoadInto(cfg); err != nil {
		return nil, err
	}

	cfg.Postgres.Password = os.Getenv(postgresPasswordKey)
	cfg.JWT.SigningKey = os.Getenv(secretKey)

	if err := config.Validate(&cfg.Config); err != nil {
		return nil, err
	}

	if cfg.JWT.SigningKey == "" {
		return nil, fmt.Errorf("empty JWT secret")
	}
	if cfg.JWT.AccessTokenTTL <= 0 {
		return nil, fmt.Errorf("invalid access token TTL")
	}
	if cfg.JWT.RefreshTokenTTL <= 0 {
		return nil, fmt.Errorf("invalid refresh token TTL")
	}
	if cfg.Hasher.Cost < bcrypt.MinCost ||
		cfg.Hasher.Cost > bcrypt.MaxCost {
		return nil, fmt.Errorf("invalid bcrypt cost: %d", cfg.Hasher.Cost)
	}

	return cfg, nil
}
