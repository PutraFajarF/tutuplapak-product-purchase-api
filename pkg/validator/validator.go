package validator

import (
	"regexp"

	v10 "github.com/go-playground/validator/v10"
)

// E.164 simple regex: + then 8..15 digits (relaxed)
var e164 = regexp.MustCompile(`^\+[1-9]\d{7,14}$`)

type CustomValidator struct{ v *v10.Validate }

func NewValidator() *CustomValidator {
	v := v10.New()
	_ = v.RegisterValidation("e164", func(fl v10.FieldLevel) bool {
		if fl.Field().Kind().String() != "string" {
			return false
		}
		s := fl.Field().String()
		return e164.MatchString(s)
	})
	return &CustomValidator{v: v}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.v.Struct(i)
}
