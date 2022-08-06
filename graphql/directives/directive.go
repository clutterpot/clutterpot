package directives

import (
	"github.com/clutterpot/clutterpot/graphql/server"
)

func New() server.DirectiveRoot {
	return server.DirectiveRoot{
		Auth: auth,
	}
}
