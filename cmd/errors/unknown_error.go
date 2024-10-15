package errors

import "fmt"

type UnknownError struct {
	message string
}

func NewUnknownError(message string) *UnknownError {
	return &UnknownError{message: message}
}

func (e *UnknownError) Error() string {
	return fmt.Sprintf("Unknown error: %s", e.message)
}
