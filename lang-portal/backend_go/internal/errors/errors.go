package errors

import "errors"

var (
	// Common errors
	ErrNotFound        = errors.New("resource not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrDatabaseError   = errors.New("database error")
	ErrInvalidPageSize = errors.New("invalid page size")
	
	// Domain-specific errors
	ErrGroupNotFound      = errors.New("group not found")
	ErrWordNotFound       = errors.New("word not found")
	ErrStudySessionNotFound = errors.New("study session not found")
	ErrActivityNotFound   = errors.New("study activity not found")
) 