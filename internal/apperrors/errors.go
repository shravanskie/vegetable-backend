package apperrors

import "fmt"

// AppError represents an application-level error with a code and a message.
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Err     error  `json:"-"`
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

// New creates a new AppError.
func New(code int, message string) *AppError {
    return &AppError{Code: code, Message: message}
}

// Wrap wraps an existing error with an AppError code and message.
func Wrap(code int, message string, err error) *AppError {
    return &AppError{Code: code, Message: message, Err: err}
}
