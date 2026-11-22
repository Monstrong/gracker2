package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/pkg/domain"
	"github.com/monstrong/gracker2/pkg/errs"
	"github.com/monstrong/gracker2/torrent-service/internal/models"
)

func (s *TorrentService) CategoriesList(ctx context.Context) ([]*models.Category, error) {
	const op = "service.torrent.listcategories"

    categories, err := s.categoriesRepo.List(ctx)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return categories, nil
}

func (s *TorrentService) StatusUpdate(ctx context.Context, torrentID uuid.UUID, newStatus int, SenderRole domain.Role) (bool, error) {
	const op = "service.torrentService.statusUpdate"
	if err := validateRole(SenderRole); err != nil {
		return false, err
	}
	if SenderRole.Weight() < domain.ModeratorRole.Weight() {
		return false, errs.ErrPermissionDenied
	}
	if newStatus > 2 || newStatus < -1 {
		return false, &errs.ValidationError{Field: "status", Reason: "invalid status"}
	}
	ok, err := s.r.StatusUpdate(ctx, torrentID, newStatus)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return ok, nil
}
func (s *TorrentService) Update(ctx context.Context, patch *models.TorrentPatch, TorrentID, SenderID uuid.UUID, SenderRole domain.Role) (*models.Torrent, error) {
	const op = "service.torrentService.update"
	if err := s.validateTorrentPatch(ctx, patch); err != nil {
		return nil, err
	}
	if err := validateID(TorrentID); err != nil {
		return nil, err
	}
	if err := validateID(SenderID); err != nil {
		return nil, err
	}
	if err := validateRole(SenderRole); err != nil {
		return nil, err
	}
	status, err := s.r.GetTorrentStatus(ctx, TorrentID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if status < 0 && SenderRole.Weight() < domain.ModeratorRole.Weight(){
		return nil, errs.ErrPermissionDenied
	}


	torrentsAuthorID, err := s.r.GetTorrentAuthorId(ctx, TorrentID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if SenderID == torrentsAuthorID {
		torrent, err := s.r.Update(ctx, patch, TorrentID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return torrent, nil
	} else {
		if SenderRole.Weight() < domain.ModeratorRole.Weight() {
			return nil, errs.ErrPermissionDenied
		}
		torrent, err := s.r.Update(ctx, patch, TorrentID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return torrent, nil
	}

}

func (s *TorrentService) Create(ctx context.Context, torrent *models.Torrent, SenderID uuid.UUID, SenderRole domain.Role) (*models.Torrent, error) {
	const op = "service.torrentService.create"

	if err := s.validateTorrent(ctx, torrent); err != nil {
		return nil, err
	}

	if err := validateID(SenderID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := validateRole(SenderRole); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if SenderRole.Weight() < domain.UserRole.Weight() {
		return nil, errs.ErrPermissionDenied
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

	if err := validateID(id); err != nil {
		return false, err
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
		if SenderRole.Weight() < domain.ModeratorRole.Weight() {
			return false, errs.ErrPermissionDenied
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


func (s *TorrentService) GetByID(ctx context.Context, id uuid.UUID, SenderRole domain.Role) (*models.Torrent, error) {
	const op = "service.torrentService.getByID"

	if err := validateID(id); err != nil {
		return nil, err
	}
	if err := validateRole(SenderRole); err != nil {
		return nil, err
	}
	status, err := s.r.GetTorrentStatus(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if status < 0 && SenderRole.Weight() < domain.ModeratorRole.Weight(){
		return nil, errs.ErrPermissionDenied
	}

	m, err := s.r.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return m, nil
}