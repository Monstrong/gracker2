package redisclient

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/monstrong/gracker2/auth-service/internal/authconfig"
)

func NewClient(ctx context.Context, cfg *authconfig.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	
	if err := client.Ping(ctx).Err(); err != nil {

		_ = client.Close()
		return nil, err
	}
	return client, nil
}

func Close(client *redis.Client) error {
	return client.Close()
}