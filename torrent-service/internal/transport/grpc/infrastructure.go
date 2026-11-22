package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/pkg/domain"
	pb "github.com/monstrong/gracker2/proto/gen/go/torrent/v1"
	"github.com/monstrong/gracker2/torrent-service/internal/models"
)

type Server struct {
	pb.UnimplementedTorrentServiceServer
	s Service
}

func NewServer(s Service) *Server {
	return &Server{s: s}
}

type Service interface {
	GetByID(ctx context.Context, id uuid.UUID, SenderRole domain.Role) (*models.Torrent, error)
	List(ctx context.Context, limit, offset int) ([]*models.Torrent, error)
	Create(ctx context.Context, torrent *models.Torrent, SenderID uuid.UUID, SenderRole domain.Role) (*models.Torrent, error)
	Delete(ctx context.Context, id uuid.UUID, SenderID uuid.UUID, SenderRole domain.Role) (bool, error)
	Update(ctx context.Context, torrent *models.TorrentPatch, TorrentID, SenderID uuid.UUID, SenderRole domain.Role) (*models.Torrent, error)
	StatusUpdate(ctx context.Context, torrentID uuid.UUID, newStatus int, SenderRole domain.Role) (bool, error)
}
