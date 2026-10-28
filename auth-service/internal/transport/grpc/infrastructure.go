package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/auth-service/internal/models"
	pb "github.com/monstrong/gracker2/proto/gen/go/auth/v1" // <-- при смене версии изменить
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	s Service
}

func NewServer(s Service) *Server{
	return &Server{s: s}
}

type Service interface {
	Register(ctx context.Context, name, email, password string) (uuid.UUID, error)
	BanUser(ctx context.Context, id uuid.UUID, SenderRole models.Role) (bool, error)
	UnbanUser(ctx context.Context, id uuid.UUID, SenderRole models.Role) (bool, error)
	GetUser(ctx context.Context, id uuid.UUID) (*models.UserPublic, error)
	GetMe(ctx context.Context, targetID uuid.UUID, myID uuid.UUID) (*models.UserPrivate, error)
	AdminGetUser(ctx context.Context, id uuid.UUID, role models.Role) (*models.UserPrivate, error)
	RoleUpdate(ctx context.Context, id uuid.UUID, TargetRole models.Role, SenderRole models.Role) (updated bool, err error)
	PremiumUpdate(ctx context.Context, id uuid.UUID, newPremium models.PremiumType, senderRole models.Role) (bool, error)
	GetRole(ctx context.Context, id uuid.UUID) (models.Role, error)
	DeleteUser(ctx context.Context, id uuid.UUID, senderRole models.Role, senderID uuid.UUID) (bool, error)
	UpdateUser(ctx context.Context, u *models.UserPatch, id uuid.UUID) (*models.UserPrivate, error)
	ChangePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) (bool, error)
	
	Login(ctx context.Context, email, password string) (access, refresh string, err error)
	Logout(ctx context.Context, rawRefresh string) (bool, error)
	LogoutEverywhere(ctx context.Context, id uuid.UUID) (bool, error)
	Refresh(ctx context.Context, rawRefresh string) (access, newRefresh string, err error)

}

