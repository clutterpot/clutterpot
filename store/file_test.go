package store

import (
	"testing"

	"github.com/clutterpot/clutterpot/db"
	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/validator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStore(t *testing.T) {
	store := New(db.Connect(), validator.New())

	t.Run("Create", func(t *testing.T) { testFileStore_Create(t, store) })
	t.Run("Update", func(t *testing.T) { testFileStore_Update(t, store) })
	t.Run("Delete", func(t *testing.T) { testFileStore_Delete(t, store) })
	t.Run("GetByID", func(t *testing.T) { testFileStore_GetByID(t, store) })
}

func testFileStore_Create(t *testing.T, s *Store) {
	file, err := s.File.Create(&model.FileInput{Name: "test"})
	require.NoError(t, err, "cannot create file")
	defer func() { require.NoError(t, s.File.Delete(file.ID)) }()
}

func testFileStore_Update(t *testing.T, s *Store) {
	file, err := s.File.Create(&model.FileInput{Name: "test"})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.File.Delete(file.ID)) }()

	_, err = s.File.Update(file.ID, &model.FileUpdateInput{})
	assert.EqualError(t, err, "file update input cannot be empty", "'file update input cannot be empty' error should have been returned")

	name := "test_" + model.NewID()

	updatedFile, err := s.File.Update(file.ID, &model.FileUpdateInput{Name: &name})
	require.NoError(t, err, "cannot update file")
	assert.Equal(t, name, updatedFile.Name, "file name should have been updated")
	assert.Equal(t, file.CreatedAt, updatedFile.CreatedAt, "file created at should not have been updated")
	assert.NotEqual(t, file.UpdatedAt, updatedFile.UpdatedAt, "file updated at should have been updated")
}

func testFileStore_Delete(t *testing.T, s *Store) {
	file, err := s.File.Create(&model.FileInput{Name: "test"})
	require.NoError(t, err)

	require.NoError(t, s.File.Delete(file.ID), "cannot delete file")
}

func testFileStore_GetByID(t *testing.T, s *Store) {
	file, err := s.File.Create(&model.FileInput{Name: "test"})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.File.Delete(file.ID)) }()

	_, err = s.File.GetByID(file.ID)
	require.NoError(t, err, "cannot get file by id")

	_, err = s.File.GetByID(model.NewID())
	require.Error(t, err, "'no rows in result set' error should have been returned")
}
