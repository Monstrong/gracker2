package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Grpc     Grpc     `mapstructure:"grpc" yaml:"grpc"`
	Infra    Infra    `mapstructure:"infra" yaml:"infra"`
	Postgres Postgres `mapstructure:"postgres" yaml:"postgres"`
	Logger   Logger   `mapstructure:"logger" yaml:"logger"`

	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout" yaml:"graceful_shutdown_timeout"`
}

type Grpc struct {
	Name    string        `mapstructure:"name" yaml:"name"`
	Port    int           `mapstructure:"port" yaml:"port"`
	Timeout time.Duration `mapstructure:"timeout" yaml:"timeout"`
}
type Infra struct {
	Name    string        `mapstructure:"name" yaml:"name"`
	Port    int           `mapstructure:"port" yaml:"port"`
	Timeout time.Duration `mapstructure:"timeout" yaml:"timeout"`
}

type Postgres struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string
	DBName   string `mapstructure:"dbname" yaml:"dbname"`
	MaxConns int    `mapstructure:"max_conns" yaml:"max_conns"`
	SSLMode  string `mapstructure:"ssl_mode"`
}
type Logger struct {
	Level string `mapstructure:"level" yaml:"level"`
}

const configYamlPath = "config/config.yaml"
const configEnvPath = "config/.env"
const postgresPasswordName = "POSTGRES_PASSWORD"

func Load() (*Config, error) {
	cfg := &Config{}
	v := viper.New()
	v.SetConfigFile(configYamlPath)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config.yaml:%w", err)
	}

	_ = godotenv.Load(configEnvPath)
	v.AutomaticEnv()

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	cfg.Postgres.Password = os.Getenv(postgresPasswordName)

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	switch {
	// grpc
	case c.Grpc.Name == "":
		return fmt.Errorf("invalid grpc name: empty")
	case c.Grpc.Port <= 0:
		return fmt.Errorf("invalid grpc port: %d, must be > 0", c.Grpc.Port)
	case c.Grpc.Timeout <= 0 || c.Grpc.Timeout > time.Minute:
		return fmt.Errorf("invalid grpc timeout: %s, must be > 0 and <= 1m", c.Grpc.Timeout)

	// infra
	case c.Infra.Name == "":
		return fmt.Errorf("invalid infra name: empty")
	case c.Infra.Port <= 0:
		return fmt.Errorf("invalid infra port: %d, must be > 0", c.Infra.Port)
	case c.Infra.Timeout <= 0 || c.Infra.Timeout > time.Minute:
		return fmt.Errorf("invalid infra timeout: %s, must be > 0 and <= 1m", c.Infra.Timeout)

	// logger
	case c.Logger.Level == "":
		return fmt.Errorf("invalid logger level: empty")

	// postgres
	case c.Postgres.SSLMode == "":
		return fmt.Errorf("invalid ssl mode: empty")
	case !isValidSSLMode(c.Postgres.SSLMode):
		return fmt.Errorf(
			"invalid ssl mode: %q, allowed: disable, allow, prefer, require, verify-ca, verify-full",
			c.Postgres.SSLMode,
		)
	case c.Postgres.MaxConns <= 0:
		return fmt.Errorf("invalid max conns: %d, must be > 0", c.Postgres.MaxConns)
	case c.Postgres.Port <= 0:
		return fmt.Errorf("invalid postgres port: %d, must be > 0", c.Postgres.Port)
	case c.Postgres.User == "":
		return fmt.Errorf("empty postgres user")
	case c.Postgres.Password == "":
		return fmt.Errorf("empty postgres password")
	case c.Postgres.DBName == "":
		return fmt.Errorf("empty postgres dbname")
	case c.Postgres.Host == "":
		return fmt.Errorf("empty postgres host")

	// graceful shutdown
	case c.GracefulShutdownTimeout <= 0 || c.GracefulShutdownTimeout > time.Minute:
		return fmt.Errorf("invalid graceful shutdown timeout: %s, must be > 0 and <= 1m", c.GracefulShutdownTimeout)
	default:
		return nil
	}
}

var validSSLModes = map[string]struct{}{
	"disable":     {},
	"allow":       {},
	"prefer":      {},
	"require":     {},
	"verify-ca":   {},
	"verify-full": {},
}

func isValidSSLMode(m string) bool {
	_, ok := validSSLModes[m]
	return ok
}
