package auth

import (
	"os"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwa"
)

type Auth struct {
	JWT *jwtauth.JWTAuth
}

func New() *Auth {
	return &Auth{JWT: jwtauth.New(jwa.HS256.String(), []byte(os.Getenv("JWT_SECRET")), nil)}
}
