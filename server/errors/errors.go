package errors

import "errors"

type BusinessError struct {
	Code    int
	Message string
	Err     error
}

func (e *BusinessError) Error() string { return e.Message }
func (e *BusinessError) Unwrap() error { return e.Err }

func NewBusinessError(code int, message string) *BusinessError {
	return &BusinessError{Code: code, Message: message}
}

var ErrPatientNotFound = errors.New("patient not found")

type DatabaseError struct{ Err error }

func (e *DatabaseError) Error() string { return "database error: " + e.Err.Error() }
func (e *DatabaseError) Unwrap() error { return e.Err }
