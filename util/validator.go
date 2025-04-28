package util

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var Validate = validator.New()

func TransformValidationErrors[T any](err error, model T) []ErrorField {
	var errors []ErrorField

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		typ := reflect.TypeOf(model)

		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		for _, fieldErr := range validationErrors {
			field, ok := typ.FieldByName(fieldErr.StructField())
			jsonTag := fieldErr.Field()

			if ok {
				tag := field.Tag.Get("json")

				if tag != "" && tag != "-" {
					jsonTag = tag
				}
			}

			errors = append(errors, ErrorField{
				Field:   jsonTag,
				Message: validationErrorMessage(fieldErr),
			})
		}
	}

	return errors
}

func validationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	// Add more cases as needed
	default:
		return fmt.Sprintf("%s is not valid", fe.Field())
	}
}
