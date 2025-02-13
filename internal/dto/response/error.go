package response

import (
	"encoding/json"
	"errors"
)

type (
	Error struct {
		Message error `json:"message" swaggertype:"string"`
	}

	internalError struct {
		Message string `json:"message"`
	}
)

func (e Error) MarshalJSON() ([]byte, error) {
	i := internalError{e.Message.Error()}
	return json.Marshal(i)
}

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrMethodNotAllowed    = errors.New("method not allowed")
	ErrNotFound            = errors.New("not found")
)
