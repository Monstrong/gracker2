package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/pkg/domain"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Id uuid.UUID `json:"id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RoleUpdateRequest struct {
	TargetUserId uuid.UUID   `json:"target_user_id"`
	NewRole      domain.Role `json:"new_role"`
}

type RoleUpdateResponse struct {
	Updated bool `json:"updated"`
}

type PremiumUpdateRequest struct {
	TargetUserId uuid.UUID          `json:"target_user_id"`
	NewPremium   domain.PremiumType `json:"new_premium"`
}

type PremiumUpdateResponse struct {
	Updated bool `json:"updated"`
}

type GetMeRequest struct{}
type GetMeResponse struct {
	User User `json:"user"`
}

type GetUserRequest struct {
	UserId uuid.UUID `json:"id"`
}
type GetUserResponse struct {
	User PublicUser `json:"user"`
}
type DeleteUserRequest struct {
	UserId uuid.UUID `json:"id"`
}
type DeleteUserResponse struct {
	Deleted bool `json:"deleted"`
}

type AdminGetUserRequest struct {
	Id uuid.UUID `json:"id"`
}
type AdminGetUserResponse struct {
	User User `json:"user"`
}

type UpdateUserRequest struct {
	Name      *string `json:"name"`
	Email     *string `json:"email"`
	AvatarURL *string `json:"avatar_url"`
	Bio       *string `json:"bio"`
}

type UpdateUserResponse struct {
	User User `json:"user"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordResponse struct {
	Changed bool `json:"changed"`
}
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutResponse struct{}
type BanUserRequest struct {
	UserId uuid.UUID `json:"id"`
}
type BanUserResponse struct {
	Banned bool `json:"banned"`
}

type UnBanUserRequest struct {
	UserId uuid.UUID `json:"id"`
}

type UnBanUserResponse struct {
	Unbanned bool `json:"unbanned"`
}
type User struct {
	Id        uuid.UUID          `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Role      domain.Role        `json:"role"`
	Premium   domain.PremiumType `json:"premium"`
	AvatarUrl string             `json:"avatar_url"`
	Bio       string             `json:"bio"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type PublicUser struct {
	Id        uuid.UUID          `json:"id"`
	Name      string             `json:"name"`
	Role      domain.Role        `json:"role"`
	Premium   domain.PremiumType `json:"premium"`
	AvatarUrl string             `json:"avatar_url"`
	Bio       string             `json:"bio"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
