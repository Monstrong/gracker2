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
	ErrWrongMethod           = errors.New("wrong method")
	ErrBannedUser            = errors.New("user is banned")

)

// insted of ErrInvalidData
type ValidationError struct {
	Field  string
	Reason string
}

type GatewayWrongFieldError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return e.Reason
}

func (e *GatewayWrongFieldError) Error() string {
	return e.Reason
}