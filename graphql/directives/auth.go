package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
)

func auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	token, _, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err = jwt.Validate(token); token == nil || err != nil {
		return nil, err
	}

	return next(ctx)
}
