package service

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/redis/go-redis/v9"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monstrong/gracker2/auth-service/internal/authconfig"
	"github.com/monstrong/gracker2/auth-service/internal/repository"
	"github.com/monstrong/gracker2/auth-service/pkg/jwtmanager"
)

// DeleteUser TODO: проверка: выполнить только если id == callerID || callerRole == admin
// password validate: if err == bcrypt.ErrPasswordTooLong (>72byte),

type UserService struct {
	m jwtmanager.TokenManager

	urepo repository.UserRepository
	rfrepo repository.RefreshRepository
}

func New(cfg *authconfig.Config, pool *pgxpool.Pool, rclient *redis.Client) *UserService {
	return &UserService{
		m: jwtmanager.NewManager(cfg.JWT.SigningKey, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL),
		urepo: repository.NewPostgresRepository(pool),
		rfrepo: repository.NewRedisRepository(rclient),
	}
}
func hashToken(token string) string {
    h := sha256.Sum256([]byte(token))
    return hex.EncodeToString(h[:])
}

