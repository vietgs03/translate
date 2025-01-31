package errors

import "fmt"

type ErrorType string

const (
	NotFound       ErrorType = "NOT_FOUND"
	ValidationErr  ErrorType = "VALIDATION_ERROR"
	DatabaseErr    ErrorType = "DATABASE_ERROR"
	Unauthorized   ErrorType = "UNAUTHORIZED"
	InternalError  ErrorType = "INTERNAL_ERROR"
)

type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}

func NewNotFoundError(format string, args ...interface{}) error {
	return AppError{
		Type:    NotFound,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewValidationError(format string, args ...interface{}) error {
	return AppError{
		Type:    ValidationErr,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewDatabaseError(format string, args ...interface{}) error {
	return AppError{
		Type:    DatabaseErr,
		Message: fmt.Sprintf(format, args...),
	}
}

func NewUnauthorizedError(format string, args ...interface{}) error {
	return AppError{
		Type:    Unauthorized,
		Message: fmt.Sprintf(format, args...),
	}
} 