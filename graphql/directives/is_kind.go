package directives

import (
	"context"

	"github.com/clutterpot/clutterpot/model"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-chi/jwtauth/v5"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func isKind(ctx context.Context, obj interface{}, next graphql.Resolver, kind model.UserKind) (interface{}, error) {
	_, claims, _ := jwtauth.FromContext(ctx)

	if model.UserKind(claims["knd"].(float64)) < kind {
		return nil, gqlerror.Errorf("permission denied")
	}

	return next(ctx)
}
