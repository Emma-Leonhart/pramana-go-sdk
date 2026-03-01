package pramana

import "fmt"

// PramanaError is returned when a Pramana OGM constraint is violated,
// such as attempting to assign an ID to an object that already has one.
type PramanaError struct {
	Message string
}

func (e *PramanaError) Error() string {
	return fmt.Sprintf("pramana: %s", e.Message)
}

// NewPramanaError creates a new PramanaError with the given message.
func NewPramanaError(msg string) *PramanaError {
	return &PramanaError{Message: msg}
}
