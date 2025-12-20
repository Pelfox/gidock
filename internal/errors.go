package internal

import "errors"

var (
	// ErrRelationNotFound indicates that the target record for a relation was not found.
	ErrRelationNotFound = errors.New("the target record not found")
	// ErrRecordNotFound indicates that the requested record was not found.
	ErrRecordNotFound = errors.New("the requested record not found")
)
