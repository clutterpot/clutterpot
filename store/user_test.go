package store

import (
	"testing"

	"github.com/clutterpot/clutterpot/db"
	"github.com/clutterpot/clutterpot/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserStore(t *testing.T) {
	store := New(db.Connect())

	t.Run("Create", func(t *testing.T) { testUserStore_Create(t, store) })
	t.Run("Update", func(t *testing.T) { testUserStore_Update(t, store) })
	t.Run("Delete", func(t *testing.T) { testUserStore_Delete(t, store) })
	t.Run("GetByID", func(t *testing.T) { testUserStore_GetByID(t, store) })
}

func testUserStore_Create(t *testing.T, s *Store) {
	user, err := s.User.Create(&model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err, "cannot create user")
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()
}

func testUserStore_Update(t *testing.T, s *Store) {
	user, err := s.User.Create(&model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	_, err = s.User.Update(user.ID, &model.UserUpdateInput{})
	assert.EqualError(t, err, "user input cannot be empty", "'user input cannot be empty' error should have been returned")

	username := "test_" + model.NewID()

	updatedUser, err := s.User.Update(user.ID, &model.UserUpdateInput{
		Username: &username,
	})
	require.NoError(t, err, "cannot update user")
	assert.Equal(t, username, updatedUser.Username, "username should have been updated")
}

func testUserStore_Delete(t *testing.T, s *Store) {
	user, err := s.User.Create(&model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)

	require.NoError(t, s.User.Delete(user.ID), "cannot delete user")
}

func testUserStore_GetByID(t *testing.T, s *Store) {
	user, err := s.User.Create(&model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, s.User.Delete(user.ID)) }()

	_, err = s.User.GetByID(user.ID)
	require.NoError(t, err, "cannot get user by id")

	_, err = s.User.GetByID(model.NewID())
	require.Error(t, err, "'no rows in result set' error should have been returned")
}
