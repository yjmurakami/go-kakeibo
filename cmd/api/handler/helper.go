package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
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

func NewJWT(c clock.Clock, key string, expiration time.Duration) jwt {
	return jwt{
		clock:         c,
		signingMethod: jwtgo.SigningMethodHS256,
		key:           []byte(key),
		expiration:    expiration,
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

// JSON

// 参考情報 https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains invalid JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains invalid JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains invalid JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains invalid JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", field)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			// dstがポインターではない場合、invalidUnmarshalError が発生する
			// プログラムに問題があるため、panicする
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return fmt.Errorf("body must only contain a single JSON value")
	}

	return nil
}

func encodeJSON(w http.ResponseWriter, status int, src interface{}, header http.Header) error {
	js, err := json.MarshalIndent(src, "", "\t")
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

// URL parameter

func getParamId(r *http.Request, key string) (int, error) {
	id, err := strconv.Atoi(mux.Vars(r)[key])
	if err != nil || id < 1 {
		return 0, fmt.Errorf("the id is invalid")
	}

	return id, nil
}
