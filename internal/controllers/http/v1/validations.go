package v1

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var rfc3339Time validator.Func = func(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if _, err := time.Parse(time.RFC3339, value); err != nil {
		return false
	}

	return true
}
