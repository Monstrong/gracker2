package migration

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/monstrong/gracker2/auth-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	l := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User, 
		cfg.Postgres.Password, 
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		l.Error("error creating migration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		l.Error("add command")
		os.Exit(1)
	}
	comand := os.Args[1]
	switch comand {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "step":
		if len(os.Args) > 2 && os.Args[2] != "" {
			steps, parseErr := strconv.Atoi(os.Args[2])
			if parseErr != nil {
				l.Error("error converting steps to int", slog.String("error", parseErr.Error()))
				os.Exit(1)
			}
			err = m.Steps(steps)
		}
	default:
		l.Error("wrong command")
		os.Exit(1)
	}
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			l.Info("database is already up")
		} else {
			l.Error("error doing migration", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}
}