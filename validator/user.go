package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func checkUsername(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString("^[A-Za-z0-9_]+$", fl.Field().String())
	return ok
}

func checkPassword(fl validator.FieldLevel) bool {
	containsLowercase := strings.ContainsAny(fl.Field().String(), "abcdefghijklmnopqrstuvwxyz")
	containsUppercase := strings.ContainsAny(fl.Field().String(), "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	containsDigit := strings.ContainsAny(fl.Field().String(), "0123456789")
	containsSpecial := strings.ContainsAny(fl.Field().String(), " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")

	if !containsLowercase || !containsUppercase || !containsDigit || !containsSpecial {
		return false
	}

	return true
}

func checkDisplayName(fl validator.FieldLevel) bool {
	ok, _ := regexp.MatchString("^[\\p{L}\\p{N}]+([ ][\\p{L}\\p{N}]+)*$", fl.Field().String())
	return ok
}
