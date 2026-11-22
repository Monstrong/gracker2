package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/pkg/domain"
	"github.com/monstrong/gracker2/pkg/errs"
	"github.com/monstrong/gracker2/torrent-service/internal/models"
)

const nameMaxLen = models.NameMaxLen
const descriptionMaxLen = models.DescriptionMaxLen

// checkCategoryID is embeded
func (s *TorrentService) validateTorrent(ctx context.Context, torrent *models.Torrent) error {
	if err := validateName(torrent.Name); err != nil {
		return err
	}
	if err := validateDescription(torrent.Description); err != nil {
		return err
	}
	if err := validateInfoHash(torrent.InfoHash); err != nil {
		return err
	}
	if err := validateID(torrent.AuthorID); err != nil {
		return err
	}
	if err := validateID(torrent.ID); err != nil {
		return err
	}
	if err := s.checkCategoryID(ctx, torrent.CategoryID); err != nil {
		return err
	}
	return nil
}

// other roles are blocked by interceptor
func validateRole(role domain.Role) error {
	switch role {
	case domain.UserRole,
		domain.ModeratorRole,
		domain.AdminRole,
		domain.OwnerRole:
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

func validateName(name string) error {
	if name == "" {
		return &errs.ValidationError{Field: "name", Reason: "'name' is required"}
	}
	if len(name) > nameMaxLen {
		return &errs.ValidationError{Field: "name", Reason: "'name' is too long"}
	}
	return nil
}

func validateDescription(description string) error {
	if description == "" {
		return &errs.ValidationError{Field: "description", Reason: "'description' is required"}
	}
	if len(description) > descriptionMaxLen {
		return &errs.ValidationError{Field: "description", Reason: "'description' is too long"}
	}
	return nil
}

func validateInfoHash(infoHash string) error {
	if infoHash == "" {
		return &errs.ValidationError{Field: "infoHash", Reason: "'infoHash' is required"}
	}
	return nil
}

// checkCategoryID is embeded
func (s *TorrentService) validateTorrentPatch(ctx context.Context, patch *models.TorrentPatch) error {
	if patch == nil {
		return &errs.ValidationError{Field: "patch", Reason: "'patch' is required"}
	}
	if patch.CategoryID == nil && patch.Description == nil && patch.Name == nil && patch.InfoHash == nil {
		return &errs.ValidationError{Field: "patch", Reason: "at least one field is required"}
	}
	if patch.CategoryID != nil {
		if err := s.checkCategoryID(ctx, *patch.CategoryID); err != nil {
			return err
		}
	}
	if patch.Name != nil {
		if err := validateName(*patch.Name); err != nil {
			return err
		}
	}
	if patch.Description != nil {
		if err := validateDescription(*patch.Description); err != nil {
			return err
		}
	}
	if patch.InfoHash != nil {
		if err := validateInfoHash(*patch.InfoHash); err != nil {
			return err
		}
	}
	return nil
}

func (s *TorrentService) checkCategoryID(ctx context.Context, categoryID uuid.UUID) error {
	if categoryID == uuid.Nil {
		return &errs.ValidationError{Field: "categoryID", Reason: "'categoryID' is required"}
	}
	ok, err := s.categoriesRepo.IsExists(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("categoriesRepo.IsExists: %w", err)
	}
	if !ok {
		return &errs.ValidationError{Field: "categoryID", Reason: "categoryID doesn't exist"}
	}
	return nil
}
