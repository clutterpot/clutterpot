package auth

import (
	"context"
	"os"

	"github.com/clutterpot/clutterpot/model"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type Auth struct {
	jwt *jwtauth.JWTAuth
}

type Claims struct {
	UserID   string
	Username string
	Kind     model.UserKind
}

func New() *Auth {
	return &Auth{jwt: jwtauth.New(jwa.HS256.String(), []byte(os.Getenv("JWT_SECRET")), nil)}
}

func (a *Auth) JWT() *jwtauth.JWTAuth {
	return a.jwt
}

func (a *Auth) Decode(token string) (jwt.Token, map[string]any, error) {
	t, _ := a.jwt.Decode(token)

	claims, err := t.AsMap(context.Background())
	if err != nil {
		return nil, nil, err
	}

	return nil, claims, nil

}

func (a *Auth) ForContext(ctx context.Context) *Claims {
	// Already validated in directive
	_, claims, _ := jwtauth.FromContext(ctx)

	return &Claims{
		UserID:   claims["uid"].(string),
		Username: claims["usr"].(string),
		Kind:     claims["knd"].(model.UserKind),
	}
}
