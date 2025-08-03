package utils

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewCustomValidator(validator *validator.Validate) *CustomValidator {
	return &CustomValidator{validator: validator}
}
