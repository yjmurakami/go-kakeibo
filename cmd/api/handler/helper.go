package handler

import (
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// JWT

type Jwt interface {
	NewToken(userId int) (string, error)
	DecodeToken(tokenString string) (int, error)
	Expiration() time.Duration
}

// Cookie

// Validator

func NewValidator() *validator.Validate {
	validate := validator.New()

	// FieldError.Field() で取得する文字列をStructのフィールド名からJSONタグに変更
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}
