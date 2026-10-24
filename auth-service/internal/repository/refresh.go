package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/monstrong/gracker2/pkg/errs"
	"github.com/redis/go-redis/v9"
)

type RefreshRepository interface {
	SetRefreshToken(ctx context.Context, tokenHash, userID string, ttl time.Duration) error
	GetRefreshToken(ctx context.Context, tokenHash string) (string, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	RotateRefreshToken(ctx context.Context, oldTokenHash, newTokenHash string, ttl time.Duration) (string, error)
}

type RedisRepository struct {
	c *redis.Client
}

func NewRedisRepository(c *redis.Client) *RedisRepository {
	return &RedisRepository{c: c}
}

var rotateRefreshScript = redis.NewScript(`
local userID = redis.call("GET", KEYS[1])
if not userID then
    return ""
end
redis.call("SET", KEYS[2], userID, "EX", ARGV[1])
redis.call("DEL", KEYS[1])
return userID
`)

func (r *RedisRepository) SetRefreshToken(ctx context.Context, tokenHash, userID string, ttl time.Duration) error {
	if ttl <= 0 {
		return fmt.Errorf("set refresh token: ttl must be positive")
	}
	if err := r.c.Set(ctx, tokenHash, userID, ttl).Err(); err != nil {
		return fmt.Errorf("set refresh token: %w", err)
	}
	return nil
}

func (r *RedisRepository) GetRefreshToken(ctx context.Context, tokenHash string) (string, error) {
	userID, err := r.c.Get(ctx, tokenHash).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", errs.ErrNotFound
		}
		return "", fmt.Errorf("get refresh token: %w", err)
	}
	return userID, nil
}

func (r *RedisRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	if err := r.c.Del(ctx, tokenHash).Err(); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}

// Если ключа нет (истёк) — ErrNotFound → 401.
// Фронтенд при 401 на refresh — редиректит на логин.

func (r *RedisRepository) RotateRefreshToken(ctx context.Context, oldTokenHash, newTokenHash string, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		return "", fmt.Errorf("rotate refresh token: ttl must be positive")
	}
	ttlSeconds := int64(math.Ceil(ttl.Seconds()))

	res, err := rotateRefreshScript.Run(ctx, r.c, []string{oldTokenHash, newTokenHash}, ttlSeconds).Result()
	if err != nil {
		return "", fmt.Errorf("rotate refresh token: %w", err)
	}

	userID, ok := res.(string)
	if !ok || userID == "" {
		return "", errs.ErrNotFound
	}
	return userID, nil
}