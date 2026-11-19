package mapper

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/api-gateway/internal/transport/httpapi/dto"
	"github.com/monstrong/gracker2/pkg/domain"

	authv1 "github.com/monstrong/gracker2/proto/gen/go/auth/v1"
)

func ToRegisterGRPC(req dto.RegisterRequest) *authv1.RegisterRequest {
	return &authv1.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}

func ToRegisterHTTP(resp *authv1.RegisterResponse) (dto.RegisterResponse, error) {
	id, err := uuid.Parse(resp.Id)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("invalid user UUID from auth-service: %w", err)
	}
	return dto.RegisterResponse{
		Id: id,
	}, nil
}
func ToLoginGRPC(req dto.LoginRequest) *authv1.LoginRequest {
	return &authv1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}

func ToLoginHTTP(resp *authv1.LoginResponse) dto.LoginResponse {
	return dto.LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}
}

func ToRefreshGRPC(req dto.RefreshRequest) *authv1.RefreshRequest {
	return &authv1.RefreshRequest{
		RefreshToken: req.RefreshToken,
	}
}

func ToRefreshHTTP(resp *authv1.RefreshResponse) dto.RefreshResponse {
	return dto.RefreshResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}
}

func ToLogoutGRPC(req dto.LogoutRequest) *authv1.LogoutRequest {
	return &authv1.LogoutRequest{
		RefreshToken: req.RefreshToken,
	}
}

// logout response is empty

// logoutEverywhere request is empty
// logoutEverywhere response is empty

func ToRoleUpdateGRPC(req dto.RoleUpdateRequest) *authv1.RoleUpdateRequest {
	return &authv1.RoleUpdateRequest{
		TargetUserId: req.TargetUserId.String(),
		NewRole:      authv1.Role(req.NewRole),
	}
}

func ToRoleUpdateHTTP(resp *authv1.RoleUpdateResponse) dto.RoleUpdateResponse {
	return dto.RoleUpdateResponse{
		Updated: resp.Updated,
	}
}

// getme request is empty

func ToGetMeHTTP(resp *authv1.GetMeResponse) (dto.GetMeResponse, error) {
	if resp == nil || resp.User == nil {
		return dto.GetMeResponse{}, fmt.Errorf("empty user response")
	}

	id, err := uuid.Parse(resp.User.Id)
	if err != nil {
		return dto.GetMeResponse{}, fmt.Errorf("invalid user UUID: %w", err)
	}

	return dto.GetMeResponse{
		User: dto.User{
			Id:        id,
			Name:      resp.User.Name,
			Email:     resp.User.Email,
			Role:      domain.Role(resp.User.Role),
			Premium:   domain.PremiumType(resp.User.Premium),
			AvatarUrl: resp.User.AvatarUrl,
			Bio:       resp.User.Bio,
			CreatedAt: resp.User.CreatedAt.AsTime(),
			UpdatedAt: resp.User.UpdatedAt.AsTime(),
		},
	}, nil
}

func ToGetUserGRPC(req dto.GetUserRequest) (*authv1.GetUserRequest, error) {
	id, err := uuid.Parse(req.UserId.String())
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	return &authv1.GetUserRequest{
		Id: id.String(),
	}, nil
}
func ToGetUserHTTP(resp *authv1.GetUserResponse) (dto.GetUserResponse, error) {
	if resp == nil || resp.User == nil {
		return dto.GetUserResponse{}, fmt.Errorf("empty user response")
	}

	id, err := uuid.Parse(resp.User.Id)
	if err != nil {
		return dto.GetUserResponse{}, fmt.Errorf("invalid user UUID: %w", err)
	}

	return dto.GetUserResponse{
		User: dto.PublicUser{
			Id:        id,
			Name:      resp.User.Name,
			Role:      domain.Role(resp.User.Role),
			Premium:   domain.PremiumType(resp.User.Premium),
			AvatarUrl: resp.User.AvatarUrl,
			Bio:       resp.User.Bio,
			CreatedAt: resp.User.CreatedAt.AsTime(),
			UpdatedAt: resp.User.UpdatedAt.AsTime(),
		},
	}, nil
}
func ToAdminGetUserGRPC(req dto.AdminGetUserRequest) (*authv1.AdminGetUserRequest, error) {
	id, err := uuid.Parse(req.Id.String())
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}
	return &authv1.AdminGetUserRequest{
		Id: id.String(),
	}, nil
}

func ToAdminGetUserHTTP(resp *authv1.AdminGetUserResponse) (dto.AdminGetUserResponse, error) {
	if resp == nil || resp.User == nil {
		return dto.AdminGetUserResponse{}, fmt.Errorf("empty user response")
	}

	id, err := uuid.Parse(resp.User.Id)
	if err != nil {
		return dto.AdminGetUserResponse{}, fmt.Errorf("invalid user UUID: %w", err)
	}

	return dto.AdminGetUserResponse{
		User: dto.User{
			Id:        id,
			Name:      resp.User.Name,
			Email:     resp.User.Email,
			Role:      domain.Role(resp.User.Role),
			Premium:   domain.PremiumType(resp.User.Premium),
			AvatarUrl: resp.User.AvatarUrl,
			Bio:       resp.User.Bio,
			CreatedAt: resp.User.CreatedAt.AsTime(),
			UpdatedAt: resp.User.UpdatedAt.AsTime(),
		},
	}, nil
}
func ToUpdateUserGRPC(req dto.UpdateUserRequest) *authv1.UpdateUserRequest {
    return &authv1.UpdateUserRequest{
        Name:      req.Name,
        Email:     req.Email,
        AvatarUrl: req.AvatarURL,
        Bio:       req.Bio,
    }
}
func ToUpdateUserHTTP(resp *authv1.UpdateUserResponse) (dto.UpdateUserResponse, error) {
    id, err := uuid.Parse(resp.User.Id)
    if err != nil {
        return dto.UpdateUserResponse{}, fmt.Errorf("invalid user UUID: %w", err)
    }
    return dto.UpdateUserResponse{
        User: dto.User{
            Id:        id,
            Name:      resp.User.Name,
            Email:     resp.User.Email,
            Role:      domain.Role(resp.User.Role),
            Premium:   domain.PremiumType(resp.User.Premium),
            AvatarUrl: resp.User.AvatarUrl,
            Bio:       resp.User.Bio,
            CreatedAt: resp.User.CreatedAt.AsTime(),
            UpdatedAt: resp.User.UpdatedAt.AsTime(),
        },
    }, nil
}
func ToDeleteUserGRPC(req dto.DeleteUserRequest) (*authv1.DeleteUserRequest, error) {
	id, err := uuid.Parse(req.UserId.String())
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}

	return &authv1.DeleteUserRequest{
		Id: id.String(),
	}, nil
}
func ToDeleteUserHTTP(resp *authv1.DeleteUserResponse) dto.DeleteUserResponse {
	return dto.DeleteUserResponse{Deleted: resp.Deleted}
}
func ToBanUserGRPC(req dto.BanUserRequest) (*authv1.BanUserRequest, error) {
	id, err := uuid.Parse(req.UserId.String())
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}

	return &authv1.BanUserRequest{
		Id: id.String(),
	}, nil
}
func ToBanUserHTTP(resp *authv1.BanUserResponse) dto.BanUserResponse {
	return dto.BanUserResponse{Banned: resp.Banned}
}
func ToUnBanUserGRPC(req dto.UnBanUserRequest) (*authv1.UnBanUserRequest, error) {
	id, err := uuid.Parse(req.UserId.String())
	if err != nil {
		return nil, fmt.Errorf("invalid user UUID: %w", err)
	}

	return &authv1.UnBanUserRequest{
		Id: id.String(),
	}, nil
}
func ToUnBanUserHTTP(resp *authv1.UnBanUserResponse) dto.UnBanUserResponse {
	return dto.UnBanUserResponse{Unbanned: resp.Unbanned}
}
func ToChangePasswordGRPC(req dto.ChangePasswordRequest) *authv1.ChangePasswordRequest {
	return &authv1.ChangePasswordRequest{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
}
func ToChangePasswordHTTP(resp *authv1.ChangePasswordResponse) dto.ChangePasswordResponse {
    return dto.ChangePasswordResponse{Updated: resp.Updated}
}
