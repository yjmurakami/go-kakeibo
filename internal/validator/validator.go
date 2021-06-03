package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var inner = validator.New()

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

func (v *Validator) AddError(key, message string) {
	if _, ok := v.errors[key]; !ok {
		v.errors[key] = message
	}
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
		v.AddError(err.Field(), err.Error())
	}
}
