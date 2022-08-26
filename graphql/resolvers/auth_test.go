package resolvers

import (
	"testing"

	"github.com/clutterpot/clutterpot/auth"
	"github.com/clutterpot/clutterpot/db"
	"github.com/clutterpot/clutterpot/graphql/directives"
	"github.com/clutterpot/clutterpot/graphql/server"
	"github.com/clutterpot/clutterpot/middlewares"
	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/store"
	"github.com/clutterpot/clutterpot/validator"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestAuthResolvers(t *testing.T) {
	auth := auth.New()
	store := store.New(db.Connect())
	gqlServer := handler.NewDefaultServer(server.NewExecutableSchema(server.Config{
		Resolvers:  New(auth, store, validator.New()),
		Directives: directives.New(),
	}))
	middlewares := chi.Chain(middlewares.New(store, auth)...)
	gqlClient := client.New(middlewares.Handler(gqlServer))

	t.Run("Login", func(t *testing.T) { testAuthResolvers_Login(t, gqlClient, store, auth) })
	t.Run("RefreshAccessToken", func(t *testing.T) { testAuthResolvers_RefreshAccessToken(t, gqlClient, store, auth) })
	t.Run("RevokeRefreshToken", func(t *testing.T) { testAuthResolvers_RevokeRefreshToken(t, gqlClient, store, auth) })
}

func testAuthResolvers_Login(t *testing.T, c *client.Client, s *store.Store, a *auth.Auth) {
	user, err := s.User.Create(model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	query := `
		mutation login($email: String!, $password: String!) {
			login(email: $email, password: $password) {
				accessToken
			}
		}
		`

	var resp struct {
		Login struct {
			AccessToken string
		}
	}
	err = c.Post(query, &resp, client.Var("email", user.Email), client.Var("password", "Password1!"))
	require.NoError(t, err, "cannot sign in with email and password")
}

func testAuthResolvers_RefreshAccessToken(t *testing.T, c *client.Client, s *store.Store, a *auth.Auth) {
	user, err := s.User.Create(model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	session, err := s.Session.Create(model.SessionInput{UserID: user.ID})
	require.NoError(t, err)

	_, token, err := a.NewRefreshToken(session)
	require.NoError(t, err)

	query := `
		mutation refreshAccessToken($refreshToken: String!) {
			refreshAccessToken(refreshToken: $refreshToken) {
				accessToken
			}
		}
		`

	var resp struct {
		RefreshAccessToken struct {
			AccessToken string
		}
	}
	err = c.Post(query, &resp, client.Var("refreshToken", token))
	require.NoError(t, err, "cannot refresh access token with refresh token")
}

func testAuthResolvers_RevokeRefreshToken(t *testing.T, c *client.Client, s *store.Store, a *auth.Auth) {
	user, err := s.User.Create(model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	session, err := s.Session.Create(model.SessionInput{UserID: user.ID})
	require.NoError(t, err)

	_, token, err := a.NewRefreshToken(session)
	require.NoError(t, err)

	query := `
		mutation revokeRefreshToken($refreshToken: String!) {
			revokeRefreshToken(refreshToken: $refreshToken) {
				refreshToken
			}
		}
		`

	var resp struct {
		RevokeRefreshToken struct {
			RefreshToken string
		}
	}
	err = c.Post(query, &resp, client.Var("refreshToken", token))
	require.NoError(t, err, "cannot revoke refresh token")
}
