package jwtmanager

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/monstrong/gracker2/auth-service/internal/models"
)

type TokenManager interface {
	NewToken(userID string, role models.Role, premium models.PremiumType) (string, error)
	NewRefreshToken() (string, error)
	Parse(tokenString string) (*Claims, error)
}

type Manager struct {
	signingKey      string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type Claims struct {
	UserID  string             `json:"user_id"`
	Role    models.Role        `json:"role"`
	Premium models.PremiumType `json:"premium"`
	jwt.RegisteredClaims
}

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

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
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
