package response

import (
	"encoding/json"
	"errors"
)

type (
	Error struct {
		Err error `json:"message" swaggertype:"string"`
	}

	internalError struct {
		Message string `json:"message"`
	}
)

var (
	InternalServerError = Error{errors.New("internal server error")}
	MethodNotAllowed    = Error{errors.New("method not allowed")}
	NotFound            = Error{errors.New("not found")}
)

func (e Error) MarshalJSON() ([]byte, error) {
	i := internalError{e.Err.Error()}
	return json.Marshal(i)
}

func ErrorOf(err error) Error {
	return Error{Err: err}
}
