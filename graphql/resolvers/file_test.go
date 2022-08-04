package resolvers

import (
	"testing"

	"github.com/clutterpot/clutterpot/db"
	"github.com/clutterpot/clutterpot/graphql/server"
	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/store"
	"github.com/clutterpot/clutterpot/validator"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileResolvers(t *testing.T) {
	store := store.New(db.Connect())
	val := validator.New()
	gqlServer := handler.NewDefaultServer(server.NewExecutableSchema(server.Config{
		Resolvers: New(store, val),
	}))
	gqlClient := client.New(gqlServer)

	t.Run("File", func(t *testing.T) { testFileResolvers_File(t, gqlClient, store) })
	t.Run("CreateFile", func(t *testing.T) { testFileResolvers_CreateFile(t, gqlClient, store) })
	t.Run("UpdateFile", func(t *testing.T) { testFileResolvers_UpdateFile(t, gqlClient, store) })
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
	require.Error(t, err, "'no rows in result set' error should have been returned")
}

func testFileResolvers_CreateFile(t *testing.T, c *client.Client, s *store.Store) {
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
	err := c.Post(query, &resp, client.Var("input", map[string]string{
		"filename": "test_" + model.NewID(),
	}))
	require.NoError(t, err, "cannot create file")
	defer func() { require.NoError(t, s.File.Delete(resp.CreateFile.ID)) }()
}

func testFileResolvers_UpdateFile(t *testing.T, c *client.Client, s *store.Store) {
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
	}))
	require.NoError(t, err, "cannot update file")
	assert.Equal(t, filename, resp.UpdateFile.Filename, "filename should have been updated")
}
