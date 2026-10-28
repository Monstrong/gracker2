package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/auth-service/internal/models"
	"github.com/monstrong/gracker2/pkg/errs"
)



func (s *UserService) Login(ctx context.Context, email, password string) (access, refresh string, err error) {
	const op = "service.User.Login"
	if err := validateEmail(email); err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	u, err := s.urepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return "", "", fmt.Errorf("%s: %w", op, errs.ErrInvalidCredentials) // не раскрываем
		}
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	if u.Role == models.BannedRole {
		return "", "", fmt.Errorf("%s: %w", op, errs.BannedUser)
	}

	if err = s.h.CheckPassword(password, u.PasswordHash); err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	access, err = s.m.NewToken(u.ID, u.Role, u.Premium)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	refresh, err = s.m.NewRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.rfrepo.SetRefreshToken(ctx, hashToken(refresh), u.ID.String(), s.m.RefreshTTL())
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return access, refresh, nil
}

func (s *UserService) Logout(ctx context.Context, rawRefresh string) (bool, error) {
    const op = "service.User.Logout"
    if rawRefresh == "" {
        return false, nil
    }
    if err := s.rfrepo.DeleteRefreshToken(ctx, hashToken(rawRefresh)); err != nil {
        return false, fmt.Errorf("%s: %w", op, err)
    }
    return true, nil
}

func (s *UserService) LogoutEverywhere(ctx context.Context, id uuid.UUID) (bool, error) {
	const op = "service.User.LogoutEverywhere"
	if err := s.rfrepo.RevokeAllUserRefreshTokens(ctx, id.String()); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}

func (s *UserService) Refresh(ctx context.Context, rawRefresh string) (access, newRefresh string, err error) {
    const op = "service.User.Refresh"

    if rawRefresh == "" {
        return "", "", fmt.Errorf("%s: %w", op, &errs.ValidationError{Field: "refresh_token", Reason: "required"})
    }

    newRaw, err := s.m.NewRefreshToken()
    if err != nil {
        return "", "", fmt.Errorf("%s: %w", op, err)
    }

    userIDStr, err := s.rfrepo.RotateRefreshToken(
        ctx, hashToken(rawRefresh), hashToken(newRaw), s.m.RefreshTTL(),
    )
    if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return "", "", fmt.Errorf("%s: %w", op, errs.ErrInvalidCredentials) // → 401
		}
        return "", "", fmt.Errorf("%s: %w", op, err)
    }

    userID, err := uuid.Parse(userIDStr)
    if err != nil {
		_ = s.rfrepo.DeleteRefreshToken(ctx, hashToken(newRaw))
        return "", "", fmt.Errorf("%s: %w", op, err)
    }

    u, err := s.urepo.GetUserPrivate(ctx, userID) // актуальные role/premium
    if err != nil {
		_ = s.rfrepo.DeleteRefreshToken(ctx, hashToken(newRaw))
        return "", "", fmt.Errorf("%s: %w", op, err)
    }

	if u.Role == models.BannedRole {
		_ = s.rfrepo.DeleteRefreshToken(ctx, hashToken(newRaw))
		return "", "", fmt.Errorf("%s: %w", op, errs.BannedUser)
	}

    access, err = s.m.NewToken(u.ID, u.Role, u.Premium)
    if err != nil {
		_ = s.rfrepo.DeleteRefreshToken(ctx, hashToken(newRaw))
        return "", "", fmt.Errorf("%s: %w", op, err)
    }

    return access, newRaw, nil
}
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}