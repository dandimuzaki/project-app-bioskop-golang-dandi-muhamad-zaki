package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// validator object struct message
func ValidateErrors(data any) ([]FieldError, error) {
	validate := validator.New()

	err := validate.Struct(data)
	if err == nil {
		return nil, nil
	}

	var errors []FieldError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			var message string
			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", err.Field())
			case "email":
				message = "Please enter a valid email format"
			case "gte":
				message = fmt.Sprintf("%s must be a non-negative number", err.Field())
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
			case "eqfield":
				message = fmt.Sprintf("%s must match %s", err.Field(), err.Param())
			default:
				message = fmt.Sprintf("%s is invalid", err.Field())
			}

			errors = append(errors, FieldError{
				Field:   err.Field(),
				Message: message,
			})
		}
		return errors, err
	}

	// Fallback: return original error if not a validation error
	return nil, err
}
