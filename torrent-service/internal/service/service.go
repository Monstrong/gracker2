package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/torrent-service/internal/models"
	"github.com/monstrong/gracker2/torrent-service/internal/repository"
)

type TorrentService struct {
	r repository.Repository
}

func NewTorrentService(r repository.Repository) *TorrentService {
	return &TorrentService{r: r}
}

func (s *TorrentService) Create(ctx context.Context, torrent *models.Torrent) (*models.Torrent, error) {
	const op = "service.torrentService.create"

	switch {
	case torrent.Name == "":
		return nil, fmt.Errorf("%s: %w: 'name' is required", op, models.ErrInvalidData)
	case torrent.InfoHash == "":
		return nil, fmt.Errorf("%s: %w: 'info_hash' is required", op, models.ErrInvalidData)
	}
	
	m, err := s.r.Create(ctx, torrent)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return m, nil
}

func (s *TorrentService) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	const op = "service.torrentService.delete"
	
	if id == uuid.Nil {
		return false, fmt.Errorf("%s: %w", op, models.ErrInvalidData)
	}

	isDeleted, err := s.r.Delete(ctx, id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isDeleted, nil
}

func (s *TorrentService) List(ctx context.Context, limit, offset int) ([]*models.Torrent, error) {
	const op = "service.torrentService.list"
	
	ms, err := s.r.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ms, nil
}

func (s *TorrentService) GetByID(ctx context.Context, id uuid.UUID) (*models.Torrent, error) {
	const op = "service.torrentService.getByID"
	
	if id == uuid.Nil {
		return nil, fmt.Errorf("%s: %w", op, models.ErrInvalidData)
	}

	m, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return m, nil
}