package utils

import "github.com/go-playground/validator/v10"

var validate = validator.New()

// ValidateStruct memvalidasi sebuah struct berdasarkan tag 'validate'
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}