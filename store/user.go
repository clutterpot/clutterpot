package store

import (
	"fmt"

	"github.com/clutterpot/clutterpot/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type userStore struct {
	db *sqlx.DB
}

func newUserStore(db *sqlx.DB) *userStore { return &userStore{db: db} }

func (us *userStore) Create(input *model.UserInput) (*model.User, error) {
	var u model.User

	if err := sq.Insert("users").SetMap(sq.Eq{
		"id":       model.NewID(),
		"username": input.Username,
		"password": model.HashPassword(input.Password),
		"email":    input.Email,
		"kind":     model.UserKindUser,
	}).Suffix("RETURNING id, username, email, kind, display_name, bio, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).RunWith(us.db).QueryRow().
		Scan(&u.ID, &u.Username, &u.Email, &u.Kind, &u.DisplayName, &u.Bio, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}

func (us *userStore) Update(id string, input *model.UserUpdateInput) (*model.User, error) {
	var u model.User

	query := sq.Update("users").Set("updated_at", "now()")

	if *input == (model.UserUpdateInput{}) {
		return nil, fmt.Errorf("user update input cannot be empty")
	}
	if input.Username != nil {
		query = query.Set("username", *input.Username)
	}
	if input.Email != nil {
		query = query.Set("email", *input.Email)
	}
	if input.Password != nil {
		query = query.Set("password", model.HashPassword(*input.Password))
	}
	if input.Kind != nil {
		query = query.Set("kind", *input.Kind)
	}
	if input.DisplayName != nil {
		query = query.Set("display_name", *input.DisplayName)
	}
	if input.Bio != nil {
		query = query.Set("bio", *input.Bio)
	}

	if err := query.Where("id = ?", id).
		Suffix("RETURNING id, username, email, kind, display_name, bio, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).RunWith(us.db).QueryRow().
		Scan(&u.ID, &u.Username, &u.Email, &u.Kind, &u.DisplayName, &u.Bio, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}

func (us *userStore) Delete(id string) error {
	if _, err := sq.Delete("users").Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).RunWith(us.db).Query(); err != nil {
		return err
	}

	return nil
}

func (us *userStore) GetByID(id string) (*model.User, error) {
	var u model.User

	if err := sq.Select("id", "username", "email", "kind", "display_name", "bio", "created_at", "updated_at").
		From("users").Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).RunWith(us.db).QueryRow().
		Scan(&u.ID, &u.Username, &u.Email, &u.Kind, &u.DisplayName, &u.Bio, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}
