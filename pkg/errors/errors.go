package errors

import "errors"

type APIError struct {
	Code        int    `json:"-"`
	Description string `json:"description"`
	Err         error  `json:"-"`
}

func (e *APIError) Error() string {
	return e.Description
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func New(msg string, code int) *APIError {
	return &APIError{
		Code:        code,
		Description: msg,
	}
}

func Wrap(err error, msg string, code int) *APIError {
	return &APIError{
		Code:        code,
		Description: msg,
		Err:         err,
	}
}

func IsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}
