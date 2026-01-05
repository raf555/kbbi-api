package httperr

import (
	"fmt"
	"net/http"
)

type httpError struct {
	inner error
	code  int
	msg   string
}

func (e *httpError) Error() string {
	if e.inner == nil {
		return fmt.Sprintf("http error %d: %s", e.code, e.msg)
	}

	return fmt.Sprintf("http error %d: %s: %s", e.code, e.msg, e.inner.Error())
}

func (e *httpError) Unwrap() error {
	return e.inner
}

func New(code int, message string) error {
	return Wrap(nil, code, message)
}

func Newf(code int, message string, args ...any) error {
	return Wrapf(nil, code, message, args...)
}

func Wrap(err error, code int, message string) error {
	return &httpError{
		inner: err,
		code:  code,
		msg:   message,
	}
}

// Wrapf wraps err with httpStatusCode and a cause from message and args.
// message and args will be an error that wraps err and will be shown to the response.
// err is kept internally.
func Wrapf(err error, code int, message string, args ...any) error {
	return &httpError{
		inner: err,
		code:  code,
		msg:   fmt.Sprintf(message, args...),
	}
}

func (h *httpError) HTTPStatusCode() int {
	return h.code
}

func (h *httpError) HTTPResponseMessage() string {
	return h.msg
}

// HTTPStatusCode returns associated status code from the err.
// If err is nil, it will return [http.StatusOK].
// If err implements HTTPStatusCoder, it will return associated status code.
// Otherwise, [http.StatusInternalServerError] is returned.
func HTTPStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if statuser, ok := err.(interface{ HTTPStatusCode() int }); ok {
		return statuser.HTTPStatusCode()
	}

	return http.StatusInternalServerError
}

func HTTPResponseMessage(err error) (string, bool) {
	if msger, ok := err.(interface{ HTTPResponseMessage() (string, bool) }); ok {
		if msg, ok := msger.HTTPResponseMessage(); ok {
			return msg, true
		}
	}

	return "", false
}
