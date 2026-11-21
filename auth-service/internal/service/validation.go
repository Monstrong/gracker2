package service

import (
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/auth-service/internal/models"
	"github.com/monstrong/gracker2/pkg/errs"
)

const (
	nameMinLen = 3
	nameMaxLen = 32

	emailMaxLen = 254 // RFC 5321

	passwordMinLen = 8
	passwordMaxLen = 72 // bcrypt лимит в БАЙТАХ

	bioMaxLen       = 500  // в рунах
	avatarURLMaxLen = 2048 // в байтах
)

var nameRe = regexp.MustCompile(`^[\p{L}\p{N}_.-]+$`)

func newValErr(field, reason string) error {
	return &errs.ValidationError{Field: field, Reason: reason}
}

func validateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return newValErr("id", "must not be empty")
	}
	return nil
}

func validateName(name string) error {
	if !utf8.ValidString(name) {
		return newValErr("name", "must be valid utf-8")
	}
	n := utf8.RuneCountInString(name)
	if n < nameMinLen || n > nameMaxLen {
		return newValErr("name", "length must be between 3 and 32 characters")
	}
	if !nameRe.MatchString(name) {
		return newValErr("name", "may contain only letters, digits, '_', '.', '-'")
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return newValErr("email", "required")
	}
	if len(email) > emailMaxLen {
		return newValErr("email", "too long")
	}
	addr, err := mail.ParseAddress(email)
	// ParseAddress принимает "Name <a@b.c>", поэтому сверяем с исходной строкой
	if err != nil || addr.Address != email {
		return newValErr("email", "invalid format")
	}
	at := strings.LastIndexByte(email, '@')
	if !strings.Contains(email[at+1:], ".") {
		return newValErr("email", "invalid domain")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < passwordMinLen {
		return newValErr("password", "must be at least 8 characters")
	}
	if len(password) > passwordMaxLen { // len() = байты, как раз то, что нужно для bcrypt
		return newValErr("password", "must be at most 72 bytes")
	}
	if !utf8.ValidString(password) {
		return newValErr("password", "must be valid utf-8")
	}
	if strings.TrimSpace(password) == "" {
		return newValErr("password", "must not be blank")
	}
	return nil
}

func validateBio(bio string) error {
	if !utf8.ValidString(bio) {
		return newValErr("bio", "must be valid utf-8")
	}
	if utf8.RuneCountInString(bio) > bioMaxLen {
		return newValErr("bio", "too long (max 500 characters)")
	}
	for _, r := range bio {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return newValErr("bio", "contains control characters")
		}
	}
	return nil
}

func validateAvatarURL(raw string) error {
	if raw == "" { // "" = нет аватарки
		return nil
	}
	if len(raw) > avatarURLMaxLen {
		return newValErr("avatar_url", "too long")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return newValErr("avatar_url", "invalid url")
	}
	if u.Scheme != "https" {
		return newValErr("avatar_url", "must use https")
	}
	if u.Host == "" || u.User != nil {
		return newValErr("avatar_url", "invalid host")
	}
	return nil
}

func validateRole(role models.Role) error {
    switch role {
    case models.BannedRole,
        models.UserRole,
        models.ModeratorRole,
        models.AdminRole,
        models.OwnerRole:
        return nil
    case models.RoleUnspecified:
        return newValErr("role", "role is unspecified")
    case models.BillingRole:
        return nil
    default:
        return newValErr("role", "unknown role")
    }
}

func validatePremiumType(p models.PremiumType) error {
    switch p {
    case models.Free, models.Premium:
        return nil
    case models.PremiumUnspecified:
        return newValErr("premium", "premium type is unspecified")
    default:
        return newValErr("premium", "unknown premium type")
    }
}
func validateUserPatch(u *models.UserPatch) error {
	if u == nil {
		return errs.ErrUserPatchEmpty
	}

	if u.AvatarURL == nil &&
		u.Bio == nil &&
		u.Name == nil &&
		u.Email == nil {
		return errs.ErrUserPatchEmpty
	}

	if u.AvatarURL != nil {
		if err := validateAvatarURL(*u.AvatarURL); err != nil {
			return err
		}
	}

	if u.Bio != nil {
		if err := validateBio(*u.Bio); err != nil {
			return err
		}
	}

	if u.Name != nil {
		if err := validateName(*u.Name); err != nil {
			return err
		}
	}

	if u.Email != nil {
		if err := validateEmail(*u.Email); err != nil {
			return err
		}
	}

	return nil
}

