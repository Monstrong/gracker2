package gatewayconfig

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/monstrong/gracker2/pkg/config"
	"github.com/spf13/viper"
)

const configYamlPath = "config/config.yaml"
const configEnvPath = "config/.env"

const secretKey = "JWT_SECRET"

type Config struct {
	Logger config.Logger `mapstructure:"logger" yaml:"logger"`

	Services   Services `mapstructure:"services" yaml:"services"`
	Server     Server   `mapstructure:"server" yaml:"server"`
	JWT_Secret string

	ChackerTimeout          time.Duration `mapstructure:"chacker_timeout" yaml:"chacker_timeout"`
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout" yaml:"graceful_shutdown_timeout"`
}

type Services struct {
	Auth    Auth    `mapstructure:"auth" yaml:"auth"`
	Torrent Torrent `mapstructure:"torrent" yaml:"torrent"`
	Tracker Tracker `mapstructure:"tracker" yaml:"tracker"`
}

type Auth struct {
	Addr      string        `mapstructure:"addr" yaml:"addr"`
	Timeout   time.Duration `mapstructure:"timeout" yaml:"timeout"`
	InfraAddr string        `mapstructure:"infra_addr" yaml:"infra_addr"`
}

type Torrent struct {
	Addr      string        `mapstructure:"addr" yaml:"addr"`
	Timeout   time.Duration `mapstructure:"timeout" yaml:"timeout"`
	InfraAddr string        `mapstructure:"infra_addr" yaml:"infra_addr"`
}

type Tracker struct {
	Addr      string        `mapstructure:"addr" yaml:"addr"`
	Timeout   time.Duration `mapstructure:"timeout" yaml:"timeout"`
	InfraAddr string        `mapstructure:"infra_addr" yaml:"infra_addr"`
}

type Server struct {
	Name    string        `mapstructure:"name" yaml:"name"`
	Port    int           `mapstructure:"port" yaml:"port"`
	Timeout time.Duration `mapstructure:"timeout" yaml:"timeout"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configYamlPath)
	cfg := &Config{}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config.yaml: %w", err)
	}

	if err := godotenv.Load(configEnvPath); errors.Is(err, os.ErrNotExist) {
		fmt.Println("no .env file found")
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	cfg.JWT_Secret = os.Getenv(secretKey)

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if cfg.ChackerTimeout == 0 {
		cfg.ChackerTimeout = 3 * time.Second
	}

	return cfg, nil
}

func (c *Config) validate() error {
	switch {
	// services
	case c.Services.Auth.Addr == "":
		return fmt.Errorf("invalid auth_addr: empty")
	case c.Services.Torrent.Addr == "":
		return fmt.Errorf("invalid torrent_addr: empty")
	case c.Services.Tracker.Addr == "":
		return fmt.Errorf("invalid tracker_addr: empty")
	case c.Services.Auth.Timeout <= 0 || c.Services.Auth.Timeout > time.Minute:
		return fmt.Errorf("invalid auth_timeout: %s, must be > 0 and <= 1m", c.Services.Auth.Timeout)
	case c.Services.Torrent.Timeout <= 0 || c.Services.Torrent.Timeout > time.Minute:
		return fmt.Errorf("invalid torrent_timeout: %s, must be > 0 and <= 1m", c.Services.Torrent.Timeout)
	case c.Services.Tracker.Timeout <= 0 || c.Services.Tracker.Timeout > time.Minute:
		return fmt.Errorf("invalid tracker_timeout: %s, must be > 0 and <= 1m", c.Services.Tracker.Timeout)
	case c.Services.Auth.InfraAddr == "":
		return fmt.Errorf("invalid auth_infra_addr: empty")
	case c.Services.Torrent.InfraAddr == "":
		return fmt.Errorf("invalid torrent_infra_addr: empty")
	case c.Services.Tracker.InfraAddr == "":
		return fmt.Errorf("invalid tracker_infra_addr: empty")

	// chacker_timeout
	case c.ChackerTimeout < 0 || c.ChackerTimeout > time.Minute:
		return fmt.Errorf("invalid chacker_timeout: %s, must be > 0 and <= 1m", c.ChackerTimeout)
	// graceful_shutdown_timeout
	case c.GracefulShutdownTimeout <= 0 || c.GracefulShutdownTimeout > time.Minute:
		return fmt.Errorf("invalid graceful_shutdown_timeout: %s, must be > 0 and <= 1m", c.GracefulShutdownTimeout)

	// server
	case c.Server.Name == "":
		return fmt.Errorf("invalid server name: empty")
	case c.Server.Port <= 0:
		return fmt.Errorf("invalid server port: %d, must be > 0", c.Server.Port)
	case c.Server.Timeout <= 0 || c.Server.Timeout > time.Minute:
		return fmt.Errorf("invalid server timeout: %s, must be > 0 and <= 1m", c.Server.Timeout)

	// logger
	case c.Logger.Level == "":
		return fmt.Errorf("invalid logger level: empty")

	default:
		return nil
	}
}
