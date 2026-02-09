package domain

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}
