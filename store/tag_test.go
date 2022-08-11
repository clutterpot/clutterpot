package store

import (
	"testing"

	"github.com/clutterpot/clutterpot/db"
	"github.com/clutterpot/clutterpot/model"
	"github.com/stretchr/testify/require"
)

func TestTagStore(t *testing.T) {
	store := New(db.Connect())

	t.Run("Create", func(t *testing.T) { testTagStore_Create(t, store) })
	t.Run("Delete", func(t *testing.T) { testTagStore_Delete(t, store) })
}

func testTagStore_Create(t *testing.T, s *Store) {
	tag, err := s.Tag.Create(model.TagInput{Name: "test_" + model.NewID()})
	require.NoError(t, err, "cannot create tag")
	defer func() { require.NoError(t, s.Tag.Delete(tag.ID)) }()
}

func testTagStore_Delete(t *testing.T, s *Store) {
	tag, err := s.Tag.Create(model.TagInput{Name: "test_" + model.NewID()})
	require.NoError(t, err)

	require.NoError(t, s.Tag.Delete(tag.ID), "cannot delete tag")
}
