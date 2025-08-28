package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func TranslateValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = field + " is required"
			case "email":
				errors[field] = "invalid email format"
			case "min":
				errors[field] = field + " must be at least " + e.Param() + " characters"
			case "max":
				errors[field] = field + " must be at most " + e.Param() + " characters"
			default:
				errors[field] = "invalid value for " + field
			}
		}
	}
	return errors
}
