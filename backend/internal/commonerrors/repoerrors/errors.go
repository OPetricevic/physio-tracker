package repoerrors

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("conflict")
)
