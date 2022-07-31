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

	t.Run("User", func(t *testing.T) { testUserResolvers_User(t, store) })
}

func testUserResolvers_User(t *testing.T, s *store.Store) {
	gqlServer := handler.NewDefaultServer(server.NewExecutableSchema(server.Config{
		Resolvers: New(s),
	}))
	gqlClient := client.New(gqlServer)

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
	err = gqlClient.Post(query, &resp, client.Var("id", user.ID))
	require.NoError(t, err, "cannot get user by id")
	assert.Equal(t, user.ID, resp.User.ID, "user id should have been the same")

	err = gqlClient.Post(query, &resp, client.Var("id", model.NewID()))
	require.Error(t, err, "'no rows in result set' error should have been returned")
}
