package resolvers

import (
	"fmt"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserResolvers(t *testing.T) {
	auth := auth.New()
	store := store.New(db.Connect())
	gqlServer := handler.NewDefaultServer(server.NewExecutableSchema(server.Config{
		Resolvers:  New(auth, store, validator.New()),
		Directives: directives.New(),
	}))
	middlewares := chi.Chain(middlewares.New(store, auth)...)
	gqlClient := client.New(middlewares.Handler(gqlServer))

	t.Run("User", func(t *testing.T) { testUserResolvers_User(t, gqlClient, store) })
	t.Run("CreateUser", func(t *testing.T) { testUserResolvers_CreateUser(t, gqlClient, store) })
	t.Run("UpdateUser", func(t *testing.T) { testUserResolvers_UpdateUser(t, gqlClient, store, auth) })
}

func testUserResolvers_User(t *testing.T, c *client.Client, s *store.Store) {
	user, err := s.User.Create(model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	query := `
		query user($id: ID!) {
			user(id: $id) {
				id
			}
		}
		`

	var resp struct {
		User struct {
			ID string
		}
	}
	err = c.Post(query, &resp, client.Var("id", user.ID))
	require.NoError(t, err, "cannot get user by id")
	assert.Equal(t, user.ID, resp.User.ID, "user id should have been the same")

	err = c.Post(query, &resp, client.Var("id", model.NewID()))
	require.Error(t, err, "'user not found' error should have been returned")
}

func testUserResolvers_CreateUser(t *testing.T, c *client.Client, s *store.Store) {
	query := `
		mutation createUser($input: UserInput!) {
			createUser(input: $input) {
				id
			}
		}
		`

	var resp struct {
		CreateUser struct {
			ID string
		}
	}
	err := c.Post(query, &resp, client.Var("input", map[string]string{
		"username": "test_" + model.NewID(),
		"email":    "test_" + model.NewID() + "@example.com",
		"password": "Password1!",
	}))
	require.NoError(t, err, "cannot create user")
	defer func() { require.NoError(t, s.User.Delete(resp.CreateUser.ID)) }()
}

func testUserResolvers_UpdateUser(t *testing.T, c *client.Client, s *store.Store, a *auth.Auth) {
	user, err := s.User.Create(model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	_, token, err := a.NewAccessToken(&auth.Claims{UserID: user.ID, Username: user.Username, Kind: model.UserKindAdmin})
	require.NoError(t, err)

	query := `
		mutation updateUser($id: ID, $input: UserUpdateInput!) {
			updateUser(id: $id, input: $input) {
				username
			}
		}
		`

	username := "test_" + model.NewID()
	var resp struct {
		UpdateUser struct {
			Username string
		}
	}
	err = c.Post(query, &resp, client.Var("input", map[string]string{"username": username}),
		client.AddHeader("Authorization", fmt.Sprintf("Bearer %s", token)))
	require.NoError(t, err, "cannot update user")
	assert.Equal(t, username, resp.UpdateUser.Username, "username should have been updated")

	err = c.Post(query, &resp, client.Var("id", model.NewID()), client.Var("input", map[string]string{
		"username": username,
	}), client.AddHeader("Authorization", fmt.Sprintf("Bearer %s", token)))
	require.Error(t, err, "'user not found' error should have been returned")
}
