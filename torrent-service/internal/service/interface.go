package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/pkg/domain"
	"github.com/monstrong/gracker2/torrent-service/internal/models"
)

type Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Torrent, error)
	List(ctx context.Context, limit, offset int) ([]*models.Torrent, error)
	Create(ctx context.Context, torrent *models.Torrent, SenderID uuid.UUID, SenderRole domain.Role) (*models.Torrent, error)
	Delete(ctx context.Context, id uuid.UUID, SenderID uuid.UUID, SenderRole domain.Role) (bool, error)
	Update(ctx context.Context, torrent *models.Torrent,TorrentID, SenderID uuid.UUID, SenderRole domain.Role) (*models.Torrent, error)
	StatusUpdate(ctx context.Context, torrentID uuid.UUID, SenderRole domain.Role) (*models.Torrent, error)
}
