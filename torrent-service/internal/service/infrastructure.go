package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/torrent-service/internal/models"
)

type TorrentService struct {
	r Repository
	categoriesRepo CategoriesRepo
}

func NewTorrentService(r Repository, categoriesRepo CategoriesRepo) *TorrentService {
	return &TorrentService{
		r: r, 	
		categoriesRepo: categoriesRepo,
	}
}

type CategoriesRepo interface {
	IsExists(ctx context.Context, categoryID uuid.UUID) (bool, error)
	List(ctx context.Context) ([]*models.Category, error)
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Torrent, error)
	List(ctx context.Context, limit, offset int) ([]*models.Torrent, error)
	Create(ctx context.Context, torrent *models.Torrent) (*models.Torrent, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	Update(ctx context.Context, patch *models.TorrentPatch, id uuid.UUID) (*models.Torrent, error)
	StatusUpdate(ctx context.Context, id uuid.UUID, status int) (bool, error)
	GetTorrentAuthorId(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	GetTorrentStatus(ctx context.Context, id uuid.UUID) (int, error)
}
