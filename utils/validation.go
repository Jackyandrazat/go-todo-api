package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) map[string]string {
	result := map[string]string{}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		result["error"] = err.Error()
		return result
	}

	for _, fieldErr := range validationErrors {
		field := strings.ToLower(fieldErr.Field())

		switch fieldErr.Tag() {
		case "required":
			result[field] = "this field is required"

		case "email":
			result[field] = "must be a valid email address"

		case "min":
			result[field] = "minimum " + fieldErr.Param() + " characters"

		case "max":
			result[field] = "maximum " + fieldErr.Param() + " characters"

		case "gt":
			result[field] = "must be greater than " + fieldErr.Param()

		case "oneof":
			result[field] = "must be one of: " + fieldErr.Param()

		default:
			result[field] = "invalid value"
		}
	}

	return result
}
