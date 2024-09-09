package model

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type requestValidator interface {
	validate() error
}

func validateError(e error) error {
	if err, ok := e.(validator.ValidationErrors); ok {
		var errors []string
		for _, err := range err {
			errors = append(errors, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
		}

		return fmt.Errorf("validate: %s", strings.Join(errors, ", "))
	}

	return fmt.Errorf("validate: %w", e)
}

func Validate(s any) error {
	if rv, ok := s.(requestValidator); ok {
		if err := rv.validate(); err != nil {
			return err
		}
	}

	if err := validate.Struct(s); err != nil {
		return validateError(err)
	}

	return nil
}

func notEmpty(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr && !field.IsNil() {
		field = field.Elem()
	}
	return field.IsValid() && !field.IsZero()
}

func init() {
	_ = validate.RegisterValidation("notempty", notEmpty)
}
