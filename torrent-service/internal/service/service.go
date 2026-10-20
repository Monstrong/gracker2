package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/torrent-service/internal/models"
	"github.com/monstrong/gracker2/pkg/errs"
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
		return nil, &errs.ValidationError{Field: "name", Reason: "'name' is required"}
	case torrent.InfoHash == "":
		return nil, &errs.ValidationError{Field: "info_hash", Reason: "'info_hash' is required"}
	case torrent.AuthorID == uuid.Nil:
		return nil, &errs.ValidationError{Field: "author_id", Reason: "'author_id' is required"}
	case torrent.CategoryID == uuid.Nil:
		return nil, &errs.ValidationError{Field: "category_id", Reason: "'category_id' is required"}
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
		return false, &errs.ValidationError{Field: "id", Reason: "invalid torrent id format"}
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
		return nil, &errs.ValidationError{Field: "id", Reason: "invalid torrent id format"}
	}

	m, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return m, nil
}
