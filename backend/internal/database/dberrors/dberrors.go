package dberrors

import (
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

// IsUniqueViolation returns true if the error represents a uniqueness/duplicate constraint violation.
// It checks the normalized GORM duplicated key error and the Postgres SQLSTATE 23505 code.
func IsUniqueViolation(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
