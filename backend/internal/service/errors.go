package service

import (
	"errors"
	"fmt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// MLServiceError is returned when the ML HTTP API responds with a non-success status.
// UpstreamStatus is the ML service status code; Message is taken from JSON "detail" when present.
type MLServiceError struct {
	UpstreamStatus int
	Message        string
}

func (e *MLServiceError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("ml service returned status %d", e.UpstreamStatus)
}
