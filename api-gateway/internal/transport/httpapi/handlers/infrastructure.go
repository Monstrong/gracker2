package handlers

import (
	"context"

	authv1 "github.com/monstrong/gracker2/proto/gen/go/auth/v1"
	torrentv1 "github.com/monstrong/gracker2/proto/gen/go/torrent/v1"
)

type AuthClient interface {
	Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error)
	Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error)
	Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error)
	Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error)
	LogoutEverywhere(ctx context.Context, req *authv1.LogoutEverywhereRequest) (*authv1.LogoutEverywhereResponse, error)
	RoleUpdate(ctx context.Context, req *authv1.RoleUpdateRequest) (*authv1.RoleUpdateResponse, error)
	GetMe(ctx context.Context, req *authv1.GetMeRequest) (*authv1.GetMeResponse, error)
	GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error)
	AdminGetUser(ctx context.Context, req *authv1.AdminGetUserRequest) (*authv1.AdminGetUserResponse, error)
	UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, req *authv1.DeleteUserRequest) (*authv1.DeleteUserResponse, error)
	BanUser(ctx context.Context, req *authv1.BanUserRequest) (*authv1.BanUserResponse, error)
	UnBanUser(ctx context.Context, req *authv1.UnBanUserRequest) (*authv1.UnBanUserResponse, error)
	ChangePassword(ctx context.Context, req *authv1.ChangePasswordRequest) (*authv1.ChangePasswordResponse, error)
}

type TorrentClient interface {
	GetTorrent(ctx context.Context, req *torrentv1.GetTorrentRequest) (*torrentv1.GetTorrentResponse, error)
	CreateTorrent(ctx context.Context, req *torrentv1.CreateTorrentRequest) (*torrentv1.CreateTorrentResponse, error)
	DeleteTorrent(ctx context.Context, req *torrentv1.DeleteTorrentRequest) (*torrentv1.DeleteTorrentResponse, error)
	ListTorrents(ctx context.Context, req *torrentv1.ListTorrentsRequest) (*torrentv1.ListTorrentsResponse, error)
	UpdateTorrent(ctx context.Context, req *torrentv1.UpdateTorrentRequest) (*torrentv1.UpdateTorrentResponse, error)
	UpdateStatus(ctx context.Context, req *torrentv1.UpdateStatusRequest) (*torrentv1.UpdateStatusResponse, error)
	ListCategory(ctx context.Context, req *torrentv1.ListCategoryRequest) (*torrentv1.ListCategoryResponse, error)
}
type TrackerClient interface {
}
