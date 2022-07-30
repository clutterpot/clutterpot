package store

import (
	"testing"

	"github.com/clutterpot/clutterpot/db"
	"github.com/clutterpot/clutterpot/model"

	"github.com/stretchr/testify/require"
)

func TestFileStore(t *testing.T) {
	store := New(db.Connect())

	t.Run("Create", func(t *testing.T) { testFileStore_Create(t, store) })
	t.Run("Delete", func(t *testing.T) { testFileStore_Delete(t, store) })
}

func testFileStore_Create(t *testing.T, s *Store) {
	file, err := s.File.Create(&model.FileInput{Name: "test"})
	require.NoError(t, err, "cannot create file")
	defer func() { require.NoError(t, s.File.Delete(file.ID)) }()
}

func testFileStore_Delete(t *testing.T, s *Store) {
	file, err := s.File.Create(&model.FileInput{Name: "test"})
	require.NoError(t, err)

	require.NoError(t, s.File.Delete(file.ID), "cannot delete file")
}
