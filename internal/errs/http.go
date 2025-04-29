package errs

import "errors"

type DefaultError struct {
	Message string `json:"message"`
}

type InvalidParamsError struct {
	Message string `json:"message"`
	Params string `json:"params"`
}

var (
	ErrInvalidParams = errors.New("invalid fields")
	ErrSintaxError = errors.New("sintax error")
)


