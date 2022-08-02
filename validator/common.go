package validator

import (
	"strconv"

	"github.com/go-playground/validator/v10"
)

func checkPrintunicode(fl validator.FieldLevel) bool {
	for _, r := range fl.Field().String() {
		if !strconv.IsPrint(r) {
			return false
		}
	}
	return true
}
