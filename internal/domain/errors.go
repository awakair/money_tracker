package domain

import "fmt"

type EntityAlreadyExistsError struct {
	EntityType string
	Field      string
	Value      string
}

func (e *EntityAlreadyExistsError) Error() string {
	return fmt.Sprintf("%s with %s '%s' already exists", e.EntityType, e.Field, e.Value)
}

type InvalidArgumentError struct {
	Argument string
	Reason   string
}

func (e *InvalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument '%s': %s", e.Argument, e.Reason)
}

type EntityNotFoundError struct {
	EntityType string
	ID         string
}

func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf("%s with id '%s' not found", e.EntityType, e.ID)
}

type AuthorizationError struct {
}

func (e *AuthorizationError) Error() string {
	return "authorization error"
}
