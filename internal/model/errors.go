package model

import "errors"

var ErrTaskNotFound = errors.New("task not found")

var ErrInvalidTask = errors.New("invalid task")

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) Unwrap() error {
	return ErrInvalidTask
}
