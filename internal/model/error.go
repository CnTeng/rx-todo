package model

import "errors"

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		Error: err.Error(),
	}
}

func (e *ErrorResponse) ToError() error {
	return errors.New(e.Error)
}
