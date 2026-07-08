package models

import "errors"

var (
	ErrNotFound = errors.New("record not found")
	ErrInvalidData = errors.New("invalid data")
	ErrInternal = errors.New("internal server error")
)