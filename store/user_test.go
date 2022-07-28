package store

import (
	"testing"

	"github.com/clutterpot/clutterpot/db"
	"github.com/clutterpot/clutterpot/model"

	"github.com/stretchr/testify/require"
)

func TestUserStore(t *testing.T) {
	store := New(db.Connect())

	t.Run("Create", func(t *testing.T) { testUserStore_Create(t, store) })
}

func testUserStore_Create(t *testing.T, s *Store) {
	_, err := s.User.Create(&model.UserInput{
		Username: "test_" + model.NewID(),
		Email:    "test_" + model.NewID() + "@example.com",
		Password: "Password1!",
	})
	require.NoError(t, err, "cannot save user")
}
