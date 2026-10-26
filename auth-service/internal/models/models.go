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

type UserPublic struct {
	ID      uuid.UUID   `db:"id"`
	Name    string      `db:"name"`
	Role    Role        `db:"role"`
	Premium PremiumType `db:"premium"`

	AvatarURL string `db:"avatar_url"` // "" = нет аватарки
	Bio       string `db:"bio"`        // "" = нет описания

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserPrivate struct {
	ID      uuid.UUID   `db:"id"`
	Name    string      `db:"name"`
	Email   string      `db:"email"`
	Role    Role        `db:"role"`
	Premium PremiumType `db:"premium"`

	AvatarURL string `db:"avatar_url"` // "" = нет аватарки
	Bio       string `db:"bio"`        // "" = нет описания

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
type UserPatch struct {
	Name      *string `json:"name"`
	Email     *string `json:"email"`
	AvatarURL *string `json:"avatar_url"`
	Bio       *string `json:"bio"`
}

const PrivateFieldsMessageString = "NO_INFO"

type Role = domain.Role

const (
	RoleUnspecified Role = domain.RoleUnspecified
	UserRole        Role = domain.UserRole
	BannedRole      Role = domain.BannedRole
	ModeratorRole   Role = domain.ModeratorRole
	AdminRole       Role = domain.AdminRole
	OwnerRole       Role = domain.OwnerRole
	BillingRole     Role = domain.BillingRole
)

type PremiumType = domain.PremiumType

const (
	PremiumUnspecified PremiumType = domain.PremiumUnspecified
	Free               PremiumType = domain.Free
	Premium            PremiumType = domain.Premium
)
