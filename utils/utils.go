// Package utils provides utility functions and helpers used across the application.
// It includes functions for generating IDs, validating structs, and other reusable utilities.
package utils

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validateInstance *validator.Validate
	once             sync.Once
)

func getValidator() *validator.Validate {
	once.Do(func() {
		validateInstance = validator.New()
	})
	return validateInstance
}

// ValidateStruct validates the fields of a struct based on the tags defined in the struct.
// a is the struct to validate.
func ValidateStruct(a any) error {
	return getValidator().Struct(a)
}
