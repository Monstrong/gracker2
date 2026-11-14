package handlers

import (
	"context"
	
	authv1 "github.com/monstrong/gracker2/proto/gen/go/auth/v1"
)


type AuthClient interface {
	Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error)
    Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error)
	Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error)
	Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error)
	LogoutEverywhere(ctx context.Context, req *authv1.LogoutEverywhereRequest) (*authv1.LogoutEverywhereResponse, error)
	RoleUpdate(ctx context.Context, req *authv1.RoleUpdateRequest) (*authv1.RoleUpdateResponse, error)
	PremiumUpdate(ctx context.Context, req *authv1.PremiumUpdateRequest) (*authv1.PremiumUpdateResponse, error)
	GetMe(ctx context.Context, req *authv1.GetMeRequest) (*authv1.GetMeResponse, error)
	GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error)
	AdminGetUser(ctx context.Context, req *authv1.AdminGetUserRequest) (*authv1.AdminGetUserResponse, error)
	UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, req *authv1.DeleteUserRequest) (*authv1.DeleteUserResponse, error)
	BanUser(ctx context.Context, req *authv1.BanUserRequest) (*authv1.BanUserResponse, error)
	UnBanUser(ctx context.Context, req *authv1.UnBanUserRequest) (*authv1.UnBanUserResponse, error)
	ChangePassword(ctx context.Context, req *authv1.ChangePasswordRequest) (*authv1.ChangePasswordResponse, error)
	Close()
}