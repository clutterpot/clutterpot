package store

import (
	"testing"

	"github.com/clutterpot/clutterpot/db"
	"github.com/clutterpot/clutterpot/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionStore(t *testing.T) {
	store := New(db.Connect())

	t.Run("Create", func(t *testing.T) { testSessionStore_Create(t, store) })
	t.Run("SoftDelete", func(t *testing.T) { testSessionStore_SoftDelete(t, store) })
	t.Run("Delete", func(t *testing.T) { testSessionStore_Delete(t, store) })
	t.Run("GetByID", func(t *testing.T) { testSessionStore_GetByID(t, store) })
}

func testSessionStore_Create(t *testing.T, s *Store) {
	user, err := s.User.Create(model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	session, err := s.Session.Create(model.SessionInput{UserID: user.ID})
	require.NoError(t, err, "cannot create session")
	defer func() { require.NoError(t, s.Session.Delete(session.ID)) }()
}

func testSessionStore_SoftDelete(t *testing.T, s *Store) {
	user, err := s.User.Create(model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	session, err := s.Session.Create(model.SessionInput{UserID: user.ID})
	require.NoError(t, err)

	_, err = s.Session.SoftDelete(session.ID)
	require.NoError(t, err, "cannot soft delete session")

	_, err = s.Session.SoftDelete(model.NewID())
	require.Error(t, err, "an error should have been returned")
	assert.EqualError(t, err, "session not found", "'session not found' error should have been returned")
}

func testSessionStore_Delete(t *testing.T, s *Store) {
	user, err := s.User.Create(model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	session, err := s.Session.Create(model.SessionInput{UserID: user.ID})
	require.NoError(t, err)

	require.NoError(t, s.Session.Delete(session.ID), "cannot delete session")
}

func testSessionStore_GetByID(t *testing.T, s *Store) {
	user, err := s.User.Create(model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	session, err := s.Session.Create(model.SessionInput{UserID: user.ID})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.Session.Delete(session.ID)) }()

	_, err = s.Session.GetByID(session.ID)
	require.NoError(t, err, "cannot get session by id")

	_, err = s.Session.GetByID(model.NewID())
	require.Error(t, err, "an error should have been returned")
	assert.EqualError(t, err, "session not found", "'session not found' error should have been returned")
}
