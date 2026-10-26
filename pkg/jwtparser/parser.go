package jwtparser

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/monstrong/gracker2/pkg/domain"
)

type TokenManager interface {
	Parse(tokenString string) (*Claims, error) // for API-gateway
}

type Manager struct {
	signingKey      string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type Claims struct {
	UserID  uuid.UUID             `json:"user_id"`
	Role    domain.Role        `json:"role"`
	Premium domain.PremiumType `json:"premium"`
	jwt.RegisteredClaims
}

func NewManager(signingKey string, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *Manager {
	return &Manager{
		signingKey:      signingKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (m *Manager) Parse(tokenString string) (*Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return &claims, nil
}
