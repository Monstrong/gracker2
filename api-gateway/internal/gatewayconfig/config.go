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

	Services Services `mapstructure:"services" yaml:"services"`
	Server Server `mapstructure:"server" yaml:"server"`
	JWT_Secret string
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout" yaml:"graceful_shutdown_timeout"`
}

type Services struct {
	AuthAddr string `mapstructure:"auth_addr" yaml:"auth_addr"`
	TorrentAddr string `mapstructure:"torrent_addr" yaml:"torrent_addr"`
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

	return cfg, nil
}

func (c *Config) validate() error {
	switch {
	// services
	case c.Services.AuthAddr == "":
		return fmt.Errorf("invalid auth_addr: empty")
	case c.Services.TorrentAddr == "":
		return fmt.Errorf("invalid torrent_addr: empty")

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