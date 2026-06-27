package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      `mapstructure:"app" yaml:"app"`
	Postgres `mapstructure:"postgres" yaml:"postgres"`
	Logger   `mapstructure:"logger" yaml:"logger"`
}

type App struct {
	Name string `mapstructure:"name" yaml:"name"`
	Port int    `mapstructure:"port" yaml:"port"`
}
type Postgres struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string
	DBName   string `mapstructure:"dbname" yaml:"dbname"`
	MaxCon   int    `mapstructure:"max_conns" yaml:"max_conns"`
}
type Logger struct {
	Level string `mapstructure:"level" yaml:"level"`
}

const configYamlPath = "config/config.yaml"
const configEnvPath = "config/.env"

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

	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")

	return cfg, nil
}
