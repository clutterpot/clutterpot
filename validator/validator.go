package validator

import "github.com/go-playground/validator/v10"

type Validator struct {
	val *validator.Validate
}

func New() *Validator {
	val := validator.New()
	val.RegisterValidation("printunicode", checkPrintunicode)

	// User
	val.RegisterValidation("username", checkUsername)
	val.RegisterValidation("password", checkPassword)
	val.RegisterValidation("displayname", checkDisplayName)

	// File
	val.RegisterValidation("filename", checkFilename)

	return &Validator{val: val}
}

func (v *Validator) Validate(m any) error {
	return v.val.Struct(m)
}
