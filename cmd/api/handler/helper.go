package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/yjmurakami/go-kakeibo/internal/clock"
)

const (
	cookieNameToken = "token"
)

// JWT

type Jwt interface {
	NewToken(userId int) (string, error)
	DecodeToken(tokenString string) (int, error)
	Expiration() time.Duration
}

type jwt struct {
	clock         clock.Clock
	signingMethod jwtgo.SigningMethod
	key           []byte
	expiration    time.Duration
}

func NewJWT(c clock.Clock, key string, expiration int) jwt {
	return jwt{
		clock:         c,
		signingMethod: jwtgo.SigningMethodHS256,
		key:           []byte(key),
		expiration:    time.Duration(expiration) * time.Minute,
	}
}

func (j jwt) NewToken(userId int) (string, error) {
	now := j.clock.Now()
	claims := jwtgo.StandardClaims{
		Subject:   strconv.Itoa(userId),
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(j.expiration).Unix(),
	}

	token := jwtgo.NewWithClaims(j.signingMethod, claims)
	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j jwt) DecodeToken(tokenString string) (int, error) {
	claims := &jwtgo.StandardClaims{}

	token, err := jwtgo.ParseWithClaims(tokenString, claims, func(token *jwtgo.Token) (interface{}, error) {
		// alg none 攻撃対策
		if token.Method == jwtgo.SigningMethodNone {
			return nil, fmt.Errorf("jwt signing method is not specified")
		}
		if token.Method != j.signingMethod {
			return nil, fmt.Errorf("jwt signing method is invalid")
		}
		return j.key, nil
	})
	if err != nil {
		return 0, errJWTInvalid
	}
	if !token.Valid {
		return 0, errJWTInvalid
	}

	// "sub" (Subject)
	sub, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, errJWTInvalid
	}

	// "iat" (Issued At)
	iat := time.Unix(claims.IssuedAt, 0)
	exp := iat.Add(j.expiration)
	if j.clock.Now().After(exp) {
		return 0, errJWTExpired
	}

	// "exp" (Expiration Time)
	// jwt-goが内部的にexpをチェックしているが、二重でチェックする
	exp = time.Unix(claims.ExpiresAt, 0)
	if j.clock.Now().After(exp) {
		return 0, errJWTExpired
	}

	return sub, nil
}

func (j jwt) Expiration() time.Duration {
	return j.expiration
}

// Cookie

func newCookie(name string, value string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	}
}

func newExpiredCookie(name string) *http.Cookie {
	c := newCookie(name, "", time.Unix(0, 0))
	c.MaxAge = -1
	return c
}

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

func encodeJSON(w http.ResponseWriter, status int, src interface{}, header http.Header) error {
	js, err := json.Marshal(src)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for k, v := range header {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)

	return err
}
