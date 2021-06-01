package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

// JSON

// 参考情報 https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	// TODO リクエスト最大サイズの設定

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		syntaxError := &json.SyntaxError{}
		unmarshalTypeError := &json.UnmarshalTypeError{}

		switch {
		case errors.As(err, &syntaxError):
			return clientError{
				status: http.StatusBadRequest,
				body:   fmt.Sprintf("request body contains invalid JSON (at position %d)", syntaxError.Offset),
			}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return clientError{
				status: http.StatusBadRequest,
				body:   "request body contains invalid JSON",
			}

		case errors.As(err, &unmarshalTypeError):
			return clientError{
				status: http.StatusBadRequest,
				body:   fmt.Sprintf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset),
			}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return clientError{
				status: http.StatusBadRequest,
				body:   fmt.Sprintf("request body contains unknown field %s", field),
			}

		case errors.Is(err, io.EOF):
			return clientError{
				status: http.StatusBadRequest,
				body:   "request body is empty",
			}

		case err.Error() == "http: request body too large":
			return clientError{
				status: http.StatusBadRequest,
				body:   "request body is too large",
			}

		default:
			// Server Error - json.InvalidUnmarshalError
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return clientError{
			status: http.StatusBadRequest,
			body:   "request body contains more than one JSON object",
		}
	}

	return nil
}

func encodeJSON(w http.ResponseWriter, src interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	return enc.Encode(src)
}
