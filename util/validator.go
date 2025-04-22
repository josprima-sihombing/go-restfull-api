package util

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool    `json:"success"`
	Errors  []Error `json:"errors"`
}

func TransformValidationErrors[T any](err error, model T) ErrorResponse {
	var errors []Error

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

			errors = append(errors, Error{
				Field:   jsonTag,
				Message: validationErrorMessage(fieldErr),
			})
		}
	}

	return ErrorResponse{
		Success: false,
		Errors:  errors,
	}
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
