package entity

import "errors"

// Error constants for the model package
var (
	// ErrCannotDeleteSystemRole is returned when attempting to delete a system role
	ErrCannotDeleteSystemRole = errors.New("cannot delete system role")

	// Add any other model-related errors here as needed
)
