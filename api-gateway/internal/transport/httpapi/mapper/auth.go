package mapper

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/api-gateway/internal/transport/httpapi/dto"
	"github.com/monstrong/gracker2/pkg/domain"

	authv1 "github.com/monstrong/gracker2/proto/gen/go/auth/v1"
)
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

func ToRefreshGRPC(req dto.RefreshRequest) *authv1.RefreshRequest {
	return &authv1.RefreshRequest{
		RefreshToken: req.RefreshToken,
	}
}