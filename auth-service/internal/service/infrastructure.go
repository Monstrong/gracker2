package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monstrong/gracker2/auth-service/internal/authconfig"
	"github.com/monstrong/gracker2/auth-service/internal/models"
	"github.com/monstrong/gracker2/auth-service/internal/repository"
	"github.com/monstrong/gracker2/auth-service/pkg/hasher"
	"github.com/monstrong/gracker2/auth-service/pkg/jwtmanager"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	m TokenManager
	h Hasher

	urepo  UserRepository
	rfrepo RefreshRepository
}

func New(cfg *authconfig.Config, pool *pgxpool.Pool, rclient *redis.Client) *UserService {
	return &UserService{
		m:      jwtmanager.NewManager(cfg.JWT.SigningKey, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL),
		h:      hasher.NewSimpleHasher(cfg.Hasher.Cost),
		urepo:  repository.NewPostgresRepository(pool),
		rfrepo: repository.NewRedisRepository(rclient),
	}
}

type Hasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) error
}

type TokenManager interface {
	NewToken(userID uuid.UUID, role models.Role, premium models.PremiumType) (string, error)
	NewRefreshToken() (string, error)
	RefreshTTL() time.Duration
}

type UserRepository interface {
	RegisterNewUser(ctx context.Context, user *models.User) (uuid.UUID, error)
	GetUserPublic(ctx context.Context, id uuid.UUID) (*models.UserPublic, error)
	GetUserPrivate(ctx context.Context, id uuid.UUID) (*models.UserPrivate, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) (bool, error)
	GetRole(ctx context.Context, id uuid.UUID) (models.Role, error)
	UpdateUserRole(ctx context.Context, id uuid.UUID, role models.Role) (updated bool, err error)
	UpdateUserPremium(ctx context.Context, id uuid.UUID, premium models.PremiumType) (updated bool, err error)
	GetPasswordHash(ctx context.Context, id uuid.UUID) (string, error)
	ChangePassword(ctx context.Context, id uuid.UUID, newPasswordHash string) (updated bool, err error)
	UpdateUser(ctx context.Context, u *models.UserPatch, id uuid.UUID) (*models.UserPrivate, error)
}

type RefreshRepository interface {
	SetRefreshToken(ctx context.Context, tokenHash, userID string, ttl time.Duration) error
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
	RotateRefreshToken(ctx context.Context, oldTokenHash, newTokenHash string, ttl time.Duration) (string, error)
}
