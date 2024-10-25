package model

import (
	"errors"
	"strings"
)

// TODO: FIX
type ErrorResponse map[string]interface{}

func NewErrorResponse(errType string, err error) *ErrorResponse {
	return &ErrorResponse{
		errType: err.Error(),
	}
}

func (e *ErrorResponse) Error() error {
	var builder strings.Builder
	for _, value := range *e {
		builder.WriteString(value.(string))
		builder.WriteString("\n")
	}
	return errors.New(builder.String())
}
