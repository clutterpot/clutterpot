package validator

import (
	"context"
	"reflect"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type Validator struct {
	val *validator.Validate
	uni *ut.UniversalTranslator
}

func New() *Validator {
	en := en.New()
	v := Validator{
		val: validator.New(),
		uni: ut.New(en, en),
	}

	v.registerTagName()
	v.registerValidations()
	v.registerTranslations()

	return &v
}

func (v *Validator) Validate(ctx context.Context, m any) error {
	err := v.val.Struct(m)
	if err != nil {
		trans, _ := v.uni.GetTranslator("en")
		for _, e := range err.(validator.ValidationErrors) {
			path := append(graphql.GetPath(ctx), ast.PathName("input"), ast.PathName(e.Field()))
			graphql.AddError(ctx, gqlerror.ErrorPathf(path, e.Translate(trans)))
		}
	}

	return err
}

func (v *Validator) registerTagName() {
	v.val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func (v *Validator) registerTranslations() {
	trans, _ := v.uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v.val, trans)

	v.val.RegisterTranslation("printunicode", trans, func(ut ut.Translator) error {
		return ut.Add("printunicode", "{0} must contain only printable unicode characters", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("printunicode", fe.Field())
		return t
	})

	v.val.RegisterTranslation("username", trans, func(ut ut.Translator) error {
		return ut.Add("username", "{0} must contain only alphanumeric characters or underscore", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("username", fe.Field())
		return t
	})

	v.val.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add("password", "{0} must contain at least one uppercase letter, one lowercase letter, one digit and one special character", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})

	v.val.RegisterTranslation("displayname", trans, func(ut ut.Translator) error {
		return ut.Add("displayname", "{0} must contain only unicode letters or numbers divided by the space character", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("displayname", fe.Field())
		return t
	})
}

func (v *Validator) registerValidations() {
	v.val.RegisterValidation("printunicode", checkPrintunicode)

	// User
	v.val.RegisterValidation("username", checkUsername)
	v.val.RegisterValidation("password", checkPassword)
	v.val.RegisterValidation("displayname", checkDisplayName)

	// File
	v.val.RegisterValidation("filename", checkFilename)
}
