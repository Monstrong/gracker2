package errs

import "errors"

var (
	ErrNotFound      = errors.New("record not found")
	ErrInternal      = errors.New("internal server error")
	ErrAlreadyExists = errors.New("record already exists")
)

// new ErrInvalidData
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return e.Reason
}
