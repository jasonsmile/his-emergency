package errors

import "errors"

var ErrPatientNotFound = errors.New("patient not found")

type DatabaseError struct{ Err error }

func (e *DatabaseError) Error() string { return "database error: " + e.Err.Error() }
func (e *DatabaseError) Unwrap() error { return e.Err }
