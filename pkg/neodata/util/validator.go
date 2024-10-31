package util

import (
	"sync"

	"github.com/go-playground/validator"
)

var (
	validate     *validator.Validate
	validateOnce sync.Once
)

// GetValidator provides a single instance of the validator.
func GetValidator() *validator.Validate {
	validateOnce.Do(func() {
		validate = validator.New()
		// Here you can add custom validation functions if needed
	})
	return validate
}

// ValidationError represents a single validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatValidationErrors formats the validator errors into a list of custom error messages.
func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errors = append(errors, ValidationError{
				Field:   fieldError.Field(),
				Message: generateErrorMessage(fieldError),
			})
		}
	}
	return errors
}

// generateErrorMessage creates a user-friendly error message.
func generateErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email address"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	default:
		return "Invalid value"
	}
}
