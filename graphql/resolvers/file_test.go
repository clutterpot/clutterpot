package resolvers

import (
	"fmt"
	"testing"

	"github.com/clutterpot/clutterpot/auth"
	"github.com/clutterpot/clutterpot/db"
	"github.com/clutterpot/clutterpot/graphql/directives"
	"github.com/clutterpot/clutterpot/graphql/server"
	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/store"
	"github.com/clutterpot/clutterpot/validator"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileResolvers(t *testing.T) {
	auth := auth.New()
	store := store.New(db.Connect())
	gqlServer := handler.NewDefaultServer(server.NewExecutableSchema(server.Config{
		Resolvers:  New(auth, store, validator.New()),
		Directives: directives.New(),
	}))
	middlewares := chi.Middlewares{jwtauth.Verifier(auth.JWT())}
	gqlClient := client.New(middlewares.Handler(gqlServer))

	t.Run("File", func(t *testing.T) { testFileResolvers_File(t, gqlClient, store) })
	t.Run("CreateFile", func(t *testing.T) { testFileResolvers_CreateFile(t, gqlClient, store, auth) })
	t.Run("UpdateFile", func(t *testing.T) { testFileResolvers_UpdateFile(t, gqlClient, store, auth) })
	t.Run("DeleteFile", func(t *testing.T) { testFileResolvers_DeleteFile(t, gqlClient, store, auth) })
}

func testFileResolvers_File(t *testing.T, c *client.Client, s *store.Store) {
	file, err := s.File.Create(model.FileInput{Filename: "test_" + model.NewID()})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.File.Delete(file.ID)) }()

	query := `
		query file($id: ID!) {
			file(id: $id) {
				id
			}
		}
		`

	var resp struct {
		File struct {
			ID string
		}
	}
	err = c.Post(query, &resp, client.Var("id", file.ID))
	require.NoError(t, err, "cannot get file by id")
	assert.Equal(t, file.ID, resp.File.ID, "file id should have been the same")

	err = c.Post(query, &resp, client.Var("id", model.NewID()))
	require.Error(t, err, "'file not found' error should have been returned")
}

func testFileResolvers_CreateFile(t *testing.T, c *client.Client, s *store.Store, a *auth.Auth) {
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
		mutation createFile($input: FileInput!) {
			createFile(input: $input) {
				id
			}
		}
		`

	var resp struct {
		CreateFile struct {
			ID string
		}
	}
	err = c.Post(query, &resp, client.Var("input", map[string]string{
		"filename": "test_" + model.NewID(),
	}), client.AddHeader("Authorization", fmt.Sprintf("Bearer %s", token)))
	require.NoError(t, err, "cannot create file")
	defer func() { require.NoError(t, s.File.Delete(resp.CreateFile.ID)) }()
}

func testFileResolvers_UpdateFile(t *testing.T, c *client.Client, s *store.Store, a *auth.Auth) {
	user, err := s.User.Create(model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	_, token, err := a.NewAccessToken(&auth.Claims{UserID: user.ID, Username: user.Username, Kind: user.Kind})
	require.NoError(t, err)

	file, err := s.File.Create(model.FileInput{Filename: "test_" + model.NewID()})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.File.Delete(file.ID)) }()

	query := `
		mutation updateFile($id: ID!, $input: FileUpdateInput!) {
			updateFile(id: $id, input: $input) {
				filename
			}
		}
		`

	filename := "test_" + model.NewID()
	var resp struct {
		UpdateFile struct {
			Filename string
		}
	}
	err = c.Post(query, &resp, client.Var("id", file.ID), client.Var("input", map[string]string{
		"filename": filename,
	}), client.AddHeader("Authorization", fmt.Sprintf("Bearer %s", token)))
	require.NoError(t, err, "cannot update file")
	assert.Equal(t, filename, resp.UpdateFile.Filename, "filename should have been updated")

	err = c.Post(query, &resp, client.Var("id", model.NewID()), client.Var("input", map[string]string{
		"filename": filename,
	}), client.AddHeader("Authorization", fmt.Sprintf("Bearer %s", token)))
	require.Error(t, err, "'file not found' error should have been returned")
}

func testFileResolvers_DeleteFile(t *testing.T, c *client.Client, s *store.Store, a *auth.Auth) {
	user, err := s.User.Create(model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	_, token, err := a.NewAccessToken(&auth.Claims{UserID: user.ID, Username: user.Username, Kind: user.Kind})
	require.NoError(t, err)

	file, err := s.File.Create(model.FileInput{Filename: "test_" + model.NewID()})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.File.Delete(file.ID)) }()

	query := `
		mutation deleteFile($id: ID!) {
			deleteFile(id: $id) {
				id
			}
		}
		`

	var resp struct {
		DeleteFile struct {
			ID string
		}
	}
	err = c.Post(query, &resp, client.Var("id", file.ID),
		client.AddHeader("Authorization", fmt.Sprintf("Bearer %s", token)))
	require.NoError(t, err, "cannot delete file")
}
