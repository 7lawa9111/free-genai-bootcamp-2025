package validation

import (
	internalErrors "lang-portal/backend_go/internal/errors"
)

const (
	MaxPageSize     = 100
	DefaultPageSize = 100
	MinPageNumber   = 1
)

// ValidatePagination validates page and limit parameters and returns processed values
func ValidatePagination(page, limit int) (int, int, error) {
	if page < MinPageNumber {
		page = MinPageNumber
	}

	if limit <= 0 || limit > MaxPageSize {
		limit = DefaultPageSize
	}

	return page, limit, nil
}

// ValidateID ensures the ID is positive
func ValidateID(id int64) error {
	if id <= 0 {
		return internalErrors.ErrInvalidInput
	}
	return nil
} 