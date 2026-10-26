package errs

import "errors"

var (
	ErrNotFound           = errors.New("record not found")
	ErrInternal           = errors.New("internal server error")
	ErrAlreadyExists      = errors.New("record already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNameAlreadyExists  = errors.New("name already exists")  // for auth
	ErrEmailAlreadyExists = errors.New("email already exists") // for auth
	ErrPermissionDenied   = errors.New("permission denied")
	ErrUserPatchEmpty     = errors.New("user patch is empty")
	ErrPasswordMismatch   = errors.New("password mismatch")
	ErrNoUpdates          = errors.New("no updates")
)

// insted of ErrInvalidData
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return e.Reason
}
