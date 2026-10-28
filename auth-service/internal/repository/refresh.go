
package repository

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/monstrong/gracker2/pkg/errs"
	"github.com/redis/go-redis/v9"
)

const userRefreshTokensPrefix = "user-refresh-tokens:"
const refreshPrefix = "refresh:"

type RedisRepository struct {
	c redis.UniversalClient
}
// redis выглядит так:
// refresh:hash_A             → "123"
// refresh:hash_B             → "123"

// user-refresh-tokens:123    → Set{"hash_A", "hash_B"}

// стоит помнить что в SET могут лежать ключи ныне не существующие. нет их автоматического удаления когда сам ключ refresh:hash_A сгорает
// Поскольку при ротации продлевается TTL множества, оно может продолжать существовать очень долго, а устаревшие хеши будут занимать память
// Однако при revokeAllRefreshTokens они удалятся. по истечении ttl самого множества,
// то есть когда пользователь вообще не будет заходить в аккаунт на протяжении месяца ни с одного устройства.
// возможное РЕШЕНИЕ - периодический проход по всей бд и удаление устаревших хэшей из set'ов. 

// При Redis Cluster важно глобально пересмотреть логику работы, так как скрипт берет сразу два ключа, а они могут храниться в разных частях

func refreshKey(tokenHash string) string {
	return refreshPrefix + tokenHash
}

func userRefreshTokensKey(userID string) string {
	return userRefreshTokensPrefix + userID
}

var deleteRefreshScript = redis.NewScript(`
	local tokenKey = KEYS[1]
	local userRefreshTokensPrefix = ARGV[1]
	local tokenHash = ARGV[2]

	local userID = redis.call("GET", tokenKey)

	if userID then
		local sessionsKey = userRefreshTokensPrefix .. userID
		redis.call("SREM", sessionsKey, tokenHash)
	end

	redis.call("DEL", tokenKey)

	return 1
`)

var rotateRefreshScript = redis.NewScript(`
	local oldTokenKey = KEYS[1]
	local newTokenKey = KEYS[2]

	local ttlSeconds = tonumber(ARGV[1])
	local userRefreshTokensPrefix = ARGV[2]
	local oldTokenHash = ARGV[3]
	local newTokenHash = ARGV[4]

	local userID = redis.call("GET", oldTokenKey)

	if not userID then
		return ""
	end

	redis.call("SET", newTokenKey, userID, "EX", ttlSeconds)

	local sessionsKey = userRefreshTokensPrefix .. userID

	redis.call("SREM", sessionsKey, oldTokenHash)
	redis.call("SADD", sessionsKey, newTokenHash)
	redis.call("EXPIRE", sessionsKey, ttlSeconds)

	redis.call("DEL", oldTokenKey)

	return userID
`)

var revokeAllRefreshTokensScript = redis.NewScript(`
	local sessionsKey = KEYS[1]
	local refreshPrefix = ARGV[1]

	local hashes = redis.call("SMEMBERS", sessionsKey)

	for _, hash in ipairs(hashes) do
		redis.call("DEL", refreshPrefix .. hash)
	end

	redis.call("DEL", sessionsKey)

	return #hashes
`)

func (r *RedisRepository) SetRefreshToken(
	ctx context.Context,
	tokenHash string,
	userID string,
	ttl time.Duration,
) error {
	const op = "repo.Redis.SetRefreshToken"

	if ttl <= 0 {
		return fmt.Errorf("%s: ttl must be positive", op)
	}

	_, err := r.c.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, refreshKey(tokenHash), userID, ttl)

		sessionsKey := userRefreshTokensKey(userID)
		pipe.SAdd(ctx, sessionsKey, tokenHash)
		pipe.Expire(ctx, sessionsKey, ttl)

		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RedisRepository) DeleteRefreshToken(
	ctx context.Context,
	tokenHash string,
) error {
	const op = "repo.Redis.DeleteRefreshToken"

	err := deleteRefreshScript.Run(
		ctx,
		r.c,
		[]string{refreshKey(tokenHash)},
		userRefreshTokensPrefix,
		tokenHash,
	).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RedisRepository) RotateRefreshToken(
	ctx context.Context,
	oldTokenHash string,
	newTokenHash string,
	ttl time.Duration,
) (string, error) {
	const op = "repo.Redis.RotateRefreshToken"

	if ttl <= 0 {
		return "", fmt.Errorf("%s: ttl must be positive", op)
	}

	ttlSeconds := int64(math.Ceil(ttl.Seconds()))

	userID, err := rotateRefreshScript.Run(
		ctx,
		r.c,
		[]string{
			refreshKey(oldTokenHash),
			refreshKey(newTokenHash),
		},
		ttlSeconds,
		userRefreshTokensPrefix,
		oldTokenHash,
		newTokenHash,
	).Text()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if userID == "" {
		return "", fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	return userID, nil
}

func (r *RedisRepository) RevokeAllUserRefreshTokens(
	ctx context.Context,
	userID string,
) error {
	const op = "repo.Redis.RevokeAllUserRefreshTokens"

	err := revokeAllRefreshTokensScript.Run(
		ctx,
		r.c,
		[]string{userRefreshTokensKey(userID)},
		refreshPrefix,
	).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
