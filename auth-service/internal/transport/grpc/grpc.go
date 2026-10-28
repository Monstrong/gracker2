package grpc

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/auth-service/internal/models"
	"github.com/monstrong/gracker2/pkg/errs"
	pb "github.com/monstrong/gracker2/proto/gen/go/auth/v1" // <-- при смене версии изменить
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const UserIDKey = "x-user-id"
const SenderRoleKey = "x-sender-role"

func (s *Server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	id, err := s.s.Register(ctx, in.GetName(), in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{Id: id.String()}, nil
}

func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	access, refresh, err := s.s.Login(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *Server) Refresh(ctx context.Context, in *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	access, refresh, err := s.s.Refresh(ctx, in.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &pb.RefreshResponse{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *Server) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {

	loggedOut, err := s.s.Logout(ctx, in.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	if !loggedOut {
		// already logged out. можно какую-нибудь ошибку возвращать
		return &pb.LogoutResponse{}, nil
	}

	return &pb.LogoutResponse{}, nil
}

// LogoutEverywhere wants UserID in metadata
func (s *Server) LogoutEverywhere(ctx context.Context, in *pb.LogoutEverywhereRequest) (*pb.LogoutEverywhereResponse, error) {
	MD, _ := metadata.FromIncomingContext(ctx)

	userID, err := checkIDfromMD(MD)
	if err != nil {
		return nil, err
	}

	loggedOut, err := s.s.LogoutEverywhere(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !loggedOut {
		// already logged out. можно какую-нибудь ошибку возвращать
		return &pb.LogoutEverywhereResponse{}, nil
	}

	return &pb.LogoutEverywhereResponse{}, nil
}

// RoleUpdate wants SenderRole in metadata
func (s *Server) RoleUpdate(ctx context.Context, in *pb.RoleUpdateRequest) (*pb.RoleUpdateResponse, error) {
	id, err := uuid.Parse(in.GetTargetUserId())
	if err != nil {
		return nil, &errs.GatewayWrongFieldError{Field: UserIDKey, Reason: "invalid user id format"}
	}
	MD, _ := metadata.FromIncomingContext(ctx)

	SenderRole, err := checkSenderRolefromMD(MD)
	if err != nil {
		return nil, err
	}

	updated, err := s.s.RoleUpdate(ctx, id, models.Role(in.GetNewRole()), SenderRole)
	if err != nil {
		return nil, err
	}
	return &pb.RoleUpdateResponse{Updated: updated}, nil
}

// PremiumUpdate wants SenderRole in metadata
func (s *Server) PremiumUpdate(ctx context.Context, in *pb.PremiumUpdateRequest) (*pb.PremiumUpdateResponse, error) {
	id, err := uuid.Parse(in.GetTargetUserId())
	if err != nil {
		return nil, &errs.GatewayWrongFieldError{Field: UserIDKey, Reason: "invalid user id format"}
	}
	MD, _ := metadata.FromIncomingContext(ctx)

	SenderRole, err := checkSenderRolefromMD(MD)
	if err != nil {
		return nil, err
	}

	updated, err := s.s.PremiumUpdate(ctx, id, models.PremiumType(in.GetNewPremium()), SenderRole)
	if err != nil {
		return nil, err
	}
	return &pb.PremiumUpdateResponse{Updated: updated}, nil
}

// GetMe wants UserID in metadata
func (s *Server) GetMe(ctx context.Context, in *pb.GetMeRequest) (*pb.GetMeResponse, error) {
	MD, _ := metadata.FromIncomingContext(ctx)

	id, err := checkIDfromMD(MD)
	if err != nil {
		return nil, err
	}

	user, err := s.s.GetMe(ctx, id, id)
	if err != nil {
		return nil, err
	}
	return &pb.GetMeResponse{User: userPrivateToPb(user)}, nil
}

func (s *Server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, &errs.GatewayWrongFieldError{Field: UserIDKey, Reason: "invalid user id format"}
	}
	user, err := s.s.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{User: userPublicToPb(user)}, nil
}

// AdminGetUser wants SenderRole in metadata
func (s *Server) AdminGetUser(ctx context.Context, in *pb.AdminGetUserRequest) (*pb.AdminGetUserResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, &errs.GatewayWrongFieldError{Field: UserIDKey, Reason: "invalid user id format"}
	}

	MD, _ := metadata.FromIncomingContext(ctx)

	SenderRole, err := checkSenderRolefromMD(MD) 
	user, err := s.s.AdminGetUser(ctx, id, SenderRole)
	if err != nil {
		return nil, err
	}
	return &pb.AdminGetUserResponse{User: userPrivateToPb(user)}, nil
}

// UpdateUser wants SenderID in metadata
func (s *Server) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	MD, _ := metadata.FromIncomingContext(ctx)

	userID, err := checkIDfromMD(MD)
	if err != nil {
		return nil, err
	}

	patch := &models.UserPatch{
		Name:      in.Name,
		Email:     in.Email,
		AvatarURL: in.AvatarUrl,
		Bio:       in.Bio,
	}

	user, err := s.s.UpdateUser(ctx, patch, userID)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserResponse{User: userPrivateToPb(user)}, nil
}

// DeleteUser wants SenderRole AND SenderID in metadata
func (s *Server) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, &errs.GatewayWrongFieldError{Field: UserIDKey, Reason: "invalid user id format"}
	}

	MD, _ := metadata.FromIncomingContext(ctx)

	SenderRole, err := checkSenderRolefromMD(MD)
	if err != nil {
		return nil, err
	}
	SenderID, err := checkIDfromMD(MD)
	if err != nil {
		return nil, err
	}

	deleted, err := s.s.DeleteUser(ctx, id, SenderRole, SenderID)
	if err != nil {
		return nil, err
	}
	
	return &pb.DeleteUserResponse{Deleted: deleted}, nil
}

// BanUser wants SenderRole in metadata
func (s *Server) BanUser(ctx context.Context, in *pb.BanUserRequest) (*pb.BanUserResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, &errs.GatewayWrongFieldError{Field: UserIDKey, Reason: "invalid user id format"}
	}

	MD, _ := metadata.FromIncomingContext(ctx)

	SenderRole, err := checkSenderRolefromMD(MD)
	if err != nil {
		return nil, err
	}

	banned, err := s.s.BanUser(ctx, id, SenderRole)
	if err != nil {
		return nil, err
	}
	return &pb.BanUserResponse{Banned: banned}, nil
}

// UnBanUser wants SenderRole in metadata
func (s *Server) UnBanUser(ctx context.Context, in *pb.UnBanUserRequest) (*pb.UnBanUserResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, &errs.GatewayWrongFieldError{Field: UserIDKey, Reason: "invalid user id format"}
	}

	MD, _ := metadata.FromIncomingContext(ctx)

	SenderRole, err := checkSenderRolefromMD(MD)
	if err != nil {
		return nil, err
	}


	unbanned, err := s.s.UnbanUser(ctx, id, SenderRole)
	if err != nil {
		return nil, err
	}
	return &pb.UnBanUserResponse{Unbanned: unbanned}, nil
}

// ChangePassword wants UserID in metadata
func (s *Server) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	
	MD, _ := metadata.FromIncomingContext(ctx)

	userID, err := checkIDfromMD(MD)
	if err != nil {
		return nil, err
	}

	updated, err := s.s.ChangePassword(ctx, userID, in.GetOldPassword(), in.GetNewPassword())
	if err != nil {
		return nil, err
	}
	return &pb.ChangePasswordResponse{Updated: updated}, nil
}

func checkIDfromMD(MD metadata.MD) (userID uuid.UUID, err error) {
	vals := MD.Get(UserIDKey)

	if len(vals) > 0 {
		if userID, err = uuid.Parse(vals[0]); err != nil {
			return uuid.Nil, &errs.GatewayWrongFieldError{Field: UserIDKey, Reason: "invalid user id format"}
		}
	} else {
		return uuid.Nil, &errs.GatewayWrongFieldError{Field: UserIDKey, Reason: "empty user id"}
	}

	return
}

func checkSenderRolefromMD(MD metadata.MD) (SenderRole models.Role, err error) {
	vals := MD.Get(SenderRoleKey)

	if len(vals) > 0 {
		roleValue, err := strconv.ParseInt(vals[0], 10, 8)
		if err != nil {
			return models.RoleUnspecified, &errs.GatewayWrongFieldError{Field: SenderRoleKey, Reason: "invalid role format"}
		}
		SenderRole = models.Role(roleValue)
	} else {
		return models.RoleUnspecified, &errs.GatewayWrongFieldError{Field: SenderRoleKey, Reason: "empty role"}
	}

	return
}

func userPrivateToPb(user *models.UserPrivate) *pb.User {
	return &pb.User{
		Id:      user.ID.String(),
		Name:    user.Name,
		Email:   user.Email,
		Role:    pb.Role(user.Role),
		Premium: pb.Premium(user.Premium),
		AvatarUrl: user.AvatarURL,
		Bio:       user.Bio,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func userPublicToPb(user *models.UserPublic) *pb.PublicUser {
	return &pb.PublicUser{
		Id:      user.ID.String(),
		Name:    user.Name,
		Role:    pb.Role(user.Role),
		Premium: pb.Premium(user.Premium),
		AvatarUrl: user.AvatarURL,
		Bio:       user.Bio,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

