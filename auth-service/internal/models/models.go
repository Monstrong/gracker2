package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/pkg/domain"
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

type UserPatch struct {
	Name         *string `json:"name"`
	Email        *string `json:"email"`
	AvatarURL    *string `json:"avatar_url"`
	Bio          *string `json:"bio"`
}

type Role = domain.Role

const (
	RoleUnspecified Role = domain.RoleUnspecified
	UserRole        Role = domain.UserRole
	BannedRole      Role = domain.BannedRole
	ModeratorRole   Role = domain.ModeratorRole
	AdminRole       Role = domain.AdminRole
)

type PremiumType = domain.PremiumType

const (
	PremiumUnspecified PremiumType = domain.PremiumUnspecified
	Free               PremiumType = domain.Free
	Premium            PremiumType = domain.Premium
)
