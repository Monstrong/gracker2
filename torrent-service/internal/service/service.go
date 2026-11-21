package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/pkg/domain"
	"github.com/monstrong/gracker2/pkg/errs"
	"github.com/monstrong/gracker2/torrent-service/internal/models"
	"github.com/monstrong/gracker2/torrent-service/internal/repository"
)

type TorrentService struct {
	r repository.Repository
}

func NewTorrentService(r repository.Repository) *TorrentService {
	return &TorrentService{r: r}
}

func (s *TorrentService) StatusUpdate(ctx context.Context, torrentID uuid.UUID, SenderRole domain.Role) (*models.Torrent, error) {
	
}
func (s *TorrentService) Update(ctx context.Context, torrent *models.Torrent, TorrentID, SenderID uuid.UUID, SenderRole domain.Role) (*models.Torrent, error) {
	if err := validateID(TorrentID); err != nil {
		return nil, err
	}
	if err := validateID(SenderID); err != nil {
		return nil, err
	}
	if err := validateRole(SenderRole); err != nil {
		return nil, err
	}


}

func (s *TorrentService) Create(ctx context.Context, torrent *models.Torrent, SenderID uuid.UUID, SenderRole domain.Role) (*models.Torrent, error) {
	const op = "service.torrentService.create"

	switch {
	case torrent.Name == "":
		return nil, &errs.ValidationError{Field: "name", Reason: "'name' is required"}
	case torrent.InfoHash == "":
		return nil, &errs.ValidationError{Field: "info_hash", Reason: "'info_hash' is required"}
	case torrent.AuthorID == uuid.Nil:
		return nil, &errs.ValidationError{Field: "authorID", Reason: "'authorID' is required"}
	case torrent.CategoryID == uuid.Nil:
		return nil, &errs.ValidationError{Field: "categoryID", Reason: "'categoryID' is required"}
	}

	if err := validateID(SenderID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := validateRole(SenderRole); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if SenderRole < domain.UserRole {
		return nil, errs.ErrPermissionDenied
	}
	if SenderRole == domain.BillingRole {
		panic("billing service is trying to create torrent")
	}

	if SenderID != torrent.AuthorID {
		return nil, errs.ErrPermissionDenied
	}

	m, err := s.r.Create(ctx, torrent)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return m, nil
}

func (s *TorrentService) Delete(ctx context.Context, id uuid.UUID, SenderID uuid.UUID, SenderRole domain.Role) (bool, error) {
	const op = "service.torrentService.delete"

	if id == uuid.Nil {
		return false, &errs.ValidationError{Field: "id", Reason: "invalid torrent id format"}
	}

	if err := validateID(SenderID); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if err := validateRole(SenderRole); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	torrentsAuthorID, err := s.r.GetTorrentAuthorId(ctx, id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if SenderID == torrentsAuthorID {
		isDeleted, err := s.r.Delete(ctx, id)
		if err != nil {
			return false, fmt.Errorf("%s: %w", op, err)
		}
		return isDeleted, nil
	} else {
		if SenderRole < domain.ModeratorRole {
			return false, errs.ErrPermissionDenied
		}
		if SenderRole == domain.BillingRole {
			panic("billing service is trying to delete torrent")
		}
		isDeleted, err := s.r.Delete(ctx, id)
		if err != nil {
			return false, fmt.Errorf("%s: %w", op, err)
		}
		return isDeleted, nil
	}
}

func (s *TorrentService) List(ctx context.Context, limit, offset int) ([]*models.Torrent, error) {
	const op = "service.torrentService.list"

	if limit <= 0 || limit > models.MaxLimit {
		return nil, &errs.ValidationError{Field: "limit", Reason: "invalid limit"}
	}
	if offset < 0 || offset > models.MaxOffset {
		return nil, &errs.ValidationError{Field: "offset", Reason: "invalid offset"}
	}
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

func validateRole(role domain.Role) error {
	switch role {
	case domain.BannedRole,
		domain.UserRole,
		domain.ModeratorRole,
		domain.AdminRole,
		domain.OwnerRole:
		return nil
	case domain.RoleUnspecified:
		return &errs.ValidationError{Field: "role", Reason: "role is unspecified"}
	case domain.BillingRole:
		return nil
	default:
		return &errs.ValidationError{Field: "role", Reason: "unknown role"}
	}
}

func validateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return &errs.ValidationError{Field: "id", Reason: "invalid torrent id format"}
	}
	return nil
}
