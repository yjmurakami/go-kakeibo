package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	inner = validator.New()
)

type Validator struct {
	inner  *validator.Validate
	errors map[string]string
}

func New() *Validator {

	// FieldError.Field() で取得する文字列をStructのフィールド名からJSONタグに変更
	inner.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	v := &Validator{
		inner:  inner,
		errors: map[string]string{},
	}

	return v
}

func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

func (v *Validator) Errors() map[string]string {
	return v.errors
}

func (v *Validator) AddError(key string, message string) {
	if _, ok := v.errors[key]; !ok {
		v.errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key string, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func InMap(value string, m map[string]bool) bool {
	_, ok := m[value]
	return ok
}

func (v *Validator) ValidateStruct(s interface{}) {
	err := v.inner.Struct(s)
	if err == nil {
		return
	}

	// Server error
	if _, ok := err.(*validator.InvalidValidationError); ok {
		panic(err)
	}

	// Client error
	for _, err := range err.(validator.ValidationErrors) {
		msg := ""

		switch err.Tag() {
		case "datetime":
			msg = fmt.Sprintf("the format must be %s", err.Param())
		case "max":
			if err.Type().Kind() == reflect.String {
				msg = fmt.Sprintf("the length must be less than or equal to %s", err.Param())
			} else if err.Type().Kind() == reflect.Int {
				msg = fmt.Sprintf("the value must be less than or equal to %s", err.Param())
			}
		case "min":
			if err.Type().Kind() == reflect.String {
				msg = fmt.Sprintf("the length must be more than or equal to %s", err.Param())
			} else if err.Type().Kind() == reflect.Int {
				msg = fmt.Sprintf("the value must be more than or equal to %s", err.Param())
			}
		case "oneof":
			msg = fmt.Sprintf("the value must be one of the followings [%s]", err.Param())
		case "required":
			msg = "the value is required"
		default:
			msg = "no message"
		}

		v.AddError(err.Field(), msg)
	}
}
