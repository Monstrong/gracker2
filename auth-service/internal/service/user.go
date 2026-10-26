package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/auth-service/internal/models"
	"github.com/monstrong/gracker2/pkg/errs"
)

// password validate: if err == bcrypt.ErrPasswordTooLong (>72byte),

// Register
func (s *UserService) Register(ctx context.Context, name, email, password string) (uuid.UUID, error) {
	const op = "service.User.Register"

	if err := validateName(name); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := validateEmail(email); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := validatePassword(password); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	passwordHash, err := s.h.HashPassword(password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.urepo.RegisterNewUser(ctx, &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil

}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*models.UserPublic, error) {
	const op = "service.User.GetUserPublic"
	if err := validateID(id); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	u, err := s.urepo.GetUserPublic(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return u, nil
}

func (s *UserService) GetMe(ctx context.Context, targetID uuid.UUID, myID uuid.UUID) (*models.UserPrivate, error) {
	const op = "service.User.GetUserPrivate"
	if err := validateID(targetID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if targetID != myID {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrPermissionDenied)
	}
	
	u, err := s.urepo.GetUserPrivate(ctx, targetID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return u, nil
}

func (s *UserService) AdminGetUser(ctx context.Context, id uuid.UUID, role models.Role) (*models.UserPrivate, error) {
	const op = "service.User.AdminGetUser"
	if err := validateID(id); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if role != models.AdminRole && role != models.OwnerRole && role != models.ModeratorRole {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrPermissionDenied)
	}
	u, err := s.urepo.GetUserPrivate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return u, nil
}

// RoleUpdate. TargetRole - the role you want to switch to. SenderRole - your role. id - target user id.
func (s *UserService) RoleUpdate(ctx context.Context, id uuid.UUID, TargetRole models.Role, SenderRole models.Role) (updated bool, err error) {
	const op = "service.User.RoleUpdate"
	if err := validateID(id); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if TargetRole == models.OwnerRole || TargetRole == models.RoleUnspecified || TargetRole == models.BillingRole ||
		SenderRole == models.RoleUnspecified || SenderRole == models.BannedRole || SenderRole == models.UserRole ||
		SenderRole == models.BillingRole {
		return false, fmt.Errorf("%s: %w", op, errs.ErrPermissionDenied)
	}

	currentRole, err := s.urepo.GetRole(ctx, id) // get target user role
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	switch currentRole {
	case models.OwnerRole:
		return false, fmt.Errorf("%s: %w", op, errs.ErrPermissionDenied)
	case models.RoleUnspecified:
		return false, fmt.Errorf("%s: %w", op, errs.ErrInternal)
	}

	if SenderRole.Weight() <= currentRole.Weight() || SenderRole.Weight() <= TargetRole.Weight() {
		return false, fmt.Errorf("%s: %w", op, errs.ErrPermissionDenied)
	}

	updated, err = s.urepo.UpdateUserRole(ctx, id, TargetRole)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if !updated {
		return false, fmt.Errorf("%s: %w", op, errs.ErrNoUpdates)
	}
	return true, nil
}

// PremiumUpdate. Billing service is tho only who is able to call this method
func (s *UserService) PremiumUpdate(ctx context.Context, id uuid.UUID, newPremium models.PremiumType, senderRole models.Role) (bool, error) {
	const op = "service.User.PremiumUpdate"
	if err := validateID(id); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if senderRole != models.BillingRole {
		return false, fmt.Errorf("%s: %w", op, errs.ErrPermissionDenied)
	}

	updated, err := s.urepo.UpdateUserPremium(ctx, id, newPremium)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if !updated {
		return false, fmt.Errorf("%s: %w", op, errs.ErrNoUpdates)
	}
	return true, nil
}

func (s *UserService) GetRole(ctx context.Context, id uuid.UUID) (models.Role, error) {
	const op = "service.User.GetRole"
	if err := validateID(id); err != nil {
		return models.RoleUnspecified, fmt.Errorf("%s: %w", op, err)
	}
	role, err := s.urepo.GetRole(ctx, id)
	if err != nil {
		return models.RoleUnspecified, fmt.Errorf("%s: %w", op, err)
	}
	return role, nil
}

// DeleteUser. Only works if callerRole >= Admin || callerID == TargetID
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID, senderRole models.Role, senderID uuid.UUID) (bool, error) {
	const op = "service.User.DeleteUser"
	if err := validateID(id); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if senderRole != models.AdminRole && senderRole != models.OwnerRole && id != senderID {
		return false, fmt.Errorf("%s: %w", op, errs.ErrPermissionDenied)
	}
	deleted, err := s.urepo.DeleteUser(ctx, id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if !deleted {
		return false, fmt.Errorf("%s: %w", op, errs.ErrNoUpdates)
	}
	return true, nil
}

func (s *UserService) UpdateUser(ctx context.Context, u *models.UserPatch, id uuid.UUID) (*models.UserPrivate, error) {
	const op = "service.User.UpdateUser"

	if err := validateID(id); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := validateUserPatch(u); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	nu, err := s.urepo.UpdateUser(ctx, u, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	
	return nu, nil
}

func (s *UserService) ChangePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) (bool, error) {
	const op = "service.User.ChangePassword"

	if err := validateID(id); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if err := validatePassword(newPassword); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	
	oldHash, err := s.urepo.GetPasswordHash(ctx, id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if err = s.h.CheckPassword(oldPassword, oldHash); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	newHash, err := s.h.HashPassword(newPassword)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	updated, err := s.urepo.ChangePassword(ctx, id, newHash)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	if !updated {
		return false, fmt.Errorf("%s: %w", op, errs.ErrNoUpdates)
	}

	return true, nil
}

func validateAvatarURL(url string) error {
	return nil
}

func validateBio(bio string) error {
	return nil
}
func validateID(id uuid.UUID) error {
	return nil
}

func validateName(name string) error {
	return nil
}

func validateEmail(email string) error {
	return nil
}

func validatePassword(password string) error {
	return nil
}

func validateUserPatch(u *models.UserPatch) error {
	if u == nil || (u.AvatarURL == nil && u.Bio == nil && u.Name == nil && u.Email == nil) {
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