package app

import (
	"github.com/go-playground/validator/v10"
)

// Validator struct
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator constructor
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

// Validate implements Echo's Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
