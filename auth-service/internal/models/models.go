package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID   `db:"id"`
	Name         string      `db:"name"`
	Email        string      `db:"email"`
	PasswordHash string      `db:"password_hash"`
	Role         Role        `db:"role"`
	Premium      PremiumType `db:"premium"`

	AvatarURL string `db:"avatar_url"` // "" = нет аватарки
	Bio       string `db:"bio"`        // "" = нет описания

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Role int8

const (
	RoleUnspecified Role = iota
	UserRole
	BannedRole
	ModeratorRole
	AdminRole
)

type PremiumType int8

const (
	PremiumUnspecified PremiumType = iota
	Free
	Premium
)
