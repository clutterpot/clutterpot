package resolvers

import (
	"testing"

	"github.com/clutterpot/clutterpot/db"
	"github.com/clutterpot/clutterpot/graphql/server"
	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/store"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserResolvers(t *testing.T) {
	store := store.New(db.Connect())
	gqlServer := handler.NewDefaultServer(server.NewExecutableSchema(server.Config{
		Resolvers: New(store),
	}))
	gqlClient := client.New(gqlServer)

	t.Run("User", func(t *testing.T) { testUserResolvers_User(t, gqlClient, store) })
	t.Run("CreateUser", func(t *testing.T) { testUserResolvers_CreateUser(t, gqlClient, store) })
}

func testUserResolvers_User(t *testing.T, c *client.Client, s *store.Store) {
	user, err := s.User.Create(&model.UserInput{
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
	require.Error(t, err, "'no rows in result set' error should have been returned")
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
