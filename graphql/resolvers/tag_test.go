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

func TestTagResolvers(t *testing.T) {
	auth := auth.New()
	store := store.New(db.Connect())
	gqlServer := handler.NewDefaultServer(server.NewExecutableSchema(server.Config{
		Resolvers:  New(auth, store, validator.New()),
		Directives: directives.New(),
	}))
	middlewares := chi.Chain(middlewares.New(store, auth)...)
	gqlClient := client.New(middlewares.Handler(gqlServer))

	t.Run("Tag", func(t *testing.T) { testTagResolvers_Tag(t, gqlClient, store) })
	t.Run("CreateTag", func(t *testing.T) { testTagResolvers_CreateTag(t, gqlClient, store, auth) })
}

func testTagResolvers_Tag(t *testing.T, c *client.Client, s *store.Store) {
	tag, err := s.Tag.Create(model.TagInput{Name: "test_" + model.NewID()})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.Tag.Delete(tag.ID)) }()

	query := `
		query tag($id: ID!) {
			tag(id: $id) {
				id
			}
		}
		`

	var resp struct {
		Tag struct {
			ID string
		}
	}
	err = c.Post(query, &resp, client.Var("id", tag.ID))
	require.NoError(t, err, "cannot get tag by id")
	assert.Equal(t, tag.ID, resp.Tag.ID, "tag id should have been the same")

	err = c.Post(query, &resp, client.Var("id", model.NewID()))
	require.Error(t, err, "'tag not found' error should have been returned")
}

func testTagResolvers_CreateTag(t *testing.T, c *client.Client, s *store.Store, a *auth.Auth) {
	user, err := s.User.Create(model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	_, token, err := a.NewAccessToken(&auth.Claims{UserID: user.ID, Username: user.Username, Kind: user.Kind})
	require.NoError(t, err)

	query := `
		mutation createTag($input: TagInput!) {
			createTag(input: $input) {
				id
			}
		}
		`

	var resp struct {
		CreateTag struct {
			ID string
		}
	}
	err = c.Post(query, &resp, client.Var("input", map[string]string{
		"name": "test_" + model.NewID(),
	}), client.AddHeader("Authorization", fmt.Sprintf("Bearer %s", token)))
	require.NoError(t, err, "cannot create tag")
	defer func() { require.NoError(t, s.Tag.Delete(resp.CreateTag.ID)) }()
}
