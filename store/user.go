package store

import (
	"github.com/clutterpot/clutterpot/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type userStore struct {
	db *sqlx.DB
}

func newUserStore(db *sqlx.DB) *userStore {
	return &userStore{db: db}
}

func (us *userStore) Create(user *model.UserInput) (*model.User, error) {
	var u model.User

	if err := sq.Insert("users").
		Columns("id", "username", "password", "email", "kind").
		Values(model.NewID(), user.Username,
			model.HashPassword(user.Password),
			user.Email, model.UserKindUser).
		Suffix("RETURNING id, username, email, kind, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		RunWith(us.db).
		QueryRow().
		Scan(&u.ID, &u.Username, &u.Email, &u.Kind, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}

func (us *userStore) Delete(id string) error {
	if _, err := sq.Delete("users").
		Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).
		RunWith(us.db).
		Query(); err != nil {
		return err
	}

	return nil
}

func (us *userStore) GetByID(id string) (*model.User, error) {
	var u model.User

	if err := sq.Select("id", "username", "email", "kind",
		"display_name", "bio", "created_at", "updated_at").
		From("users").
		Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).
		RunWith(us.db).
		QueryRow().
		Scan(&u.ID, &u.Username, &u.Email, &u.Kind, &u.DisplayName,
			&u.Bio, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}