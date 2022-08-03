package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func checkFilename(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString("^[^/>|:&\\p{Cc}\\p{Zl}\\p{Zp}]*$", fl.Field().String())
	return ok
}
