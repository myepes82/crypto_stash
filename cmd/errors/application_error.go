package errors

import "fmt"

type ApplicationError struct {
	message string
}

func NewApplicationError(message string) *ApplicationError {
	return &ApplicationError{message: message}
}

func (e *ApplicationError) Error() string {
	return fmt.Sprintf("Application error: %s", e.message)
}

var (
	NoSecretKeyFileFoundError = &ApplicationError{"no secret key file were found."}
)
