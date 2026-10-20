package main

import (
	"errors"
	"log"
	"log/slog"
	"net"
	"net/url"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/monstrong/gracker2/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	l := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dsn := buildDSN(&cfg.Postgres)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		l.Error("error creating migration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		l.Error("add command")
		os.Exit(1)
	}
	command := os.Args[1]
	switch command {
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
		} else {
			l.Error("add number of steps")
			os.Exit(1)
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

func buildDSN(cfg *config.Postgres) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Path:   cfg.DBName,
	}
	q := u.Query()
	q.Set("sslmode", cfg.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}
