package internal

import "errors"

var (
	// ErrRelationNotFound indicates that the target record for a relation was not found.
	ErrRelationNotFound = errors.New("the target record not found")
	// ErrRecordNotFound indicates that the requested record was not found.
	ErrRecordNotFound = errors.New("the requested record not found")
	// ErrNoContainer indicates that the service has no attached container to it.
	ErrNoContainer = errors.New("service has no associated container")
	// ErrNoFields indicates that no fields were provided for an update operation.
	ErrNoFields = errors.New("no fields to update")
)
