package jwtmanager

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/monstrong/gracker2/auth-service/internal/models"
	"github.com/monstrong/gracker2/pkg/jwtparser"
)

type TokenManager interface {
	NewToken(userID string, role models.Role, premium models.PremiumType) (string, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey      string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type Claims = jwtparser.Claims

func NewManager(signingKey string, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *Manager {
	return &Manager{
		signingKey:      signingKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (m *Manager) NewToken(userID string, role models.Role, premium models.PremiumType) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID:  userID,
		Role:    role,
		Premium: premium,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	token, err := t.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
