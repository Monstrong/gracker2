package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/torrent-service/internal/models"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Torrent, error)
	List(ctx context.Context, limit, offset int) ([]*models.Torrent, error)
	Create(ctx context.Context, torrent *models.Torrent) (*models.Torrent, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
}