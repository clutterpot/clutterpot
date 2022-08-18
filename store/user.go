package store

import (
	"database/sql"
	"fmt"

	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/store/pagination"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type userStore struct {
	db *sqlx.DB
}

func newUserStore(db *sqlx.DB) *userStore {
	return &userStore{db: db}
}

func (us *userStore) Create(input model.UserInput) (*model.User, error) {
	query, args, err := sq.Insert("users").SetMap(sq.Eq{
		"id":       model.NewID(),
		"username": input.Username,
		"password": model.HashPassword(input.Password),
		"email":    input.Email,
		"kind":     model.UserKindUser,
	}).Suffix("RETURNING id, username, email, kind, display_name, bio, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating a user")
	}

	var u model.User
	if err = us.db.Get(&u, query, args...); err != nil {
		return nil, fmt.Errorf("an error occurred while updating a user")
	}

	return &u, nil
}

func (us *userStore) Update(id string, input model.UserUpdateInput) (*model.User, error) {
	query := sq.Update("users").Set("updated_at", "now()")

	if input == (model.UserUpdateInput{}) {
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

	updateQuery, args, err := query.Where("id = ?", id).
		Suffix("RETURNING id, username, email, kind, display_name, bio, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while updating a user")
	}

	var u model.User
	if err = us.db.Get(&u, updateQuery, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("an error occurred while updating a user")
	}

	return &u, nil
}

func (us *userStore) Delete(id string) error {
	if _, err := sq.Delete("users").Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).RunWith(us.db).Query(); err != nil {
		return fmt.Errorf("an error occurred while deleting a user")
	}

	return nil
}

func (us *userStore) GetByID(id string) (*model.User, error) {
	query, args, err := sq.Select("id", "username", "email", "kind", "display_name", "bio", "created_at", "updated_at").
		From("users").Where("id = ?", id).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while getting a user by id")
	}

	var u model.User
	if err = us.db.Get(&u, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("an error occurred while getting a user by id")
	}

	return &u, nil
}

// Only for auth purposes
func (us *userStore) GetByEmail(email string) (*model.AuthUser, error) {
	query, args, err := sq.Select("*").From("users").
		Where("email = ?", email).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while getting a user by email")
	}

	var u model.AuthUser
	if err = us.db.Get(&u, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid email or password")
		}

		return nil, fmt.Errorf("an error occurred while getting a user by email")
	}

	return &u, nil
}

func (us *userStore) GetAll(after, before *string, first, last *int, sort *model.UserSort, order *model.Order) ([]*model.User, *model.PageInfo, error) {
	// Base query
	query := sq.Select("id", "username", "email", "kind", "display_name", "bio", "created_at", "updated_at").
		From("users")

	// Build paginated query
	paginatedQuery, args, err := pagination.PaginateQuery(query, after, before, first, last, string(*sort), *order)
	if err != nil {
		return nil, nil, err
	}

	var urs []*model.User
	if err = us.db.Select(&urs, paginatedQuery, args...); err != nil || len(urs) == 0 {
		if err == sql.ErrNoRows || len(urs) == 0 {
			return nil, nil, fmt.Errorf("users not found")
		}

		return nil, nil, fmt.Errorf("an error occurred while getting all users")
	}

	// Build page info query
	pageInfoQuery, args, err := pagination.GetPageInfoQuery(query, urs[0].ID, urs[len(urs)-1].ID, *order)
	if err != nil {
		return nil, nil, fmt.Errorf("an error occurred while getting all users")
	}

	// Query HasNextPage and HasPrevious page fields od page info
	var hasPages []bool
	if err = us.db.Select(&hasPages, pageInfoQuery, args...); err != nil {
		return nil, nil, fmt.Errorf("an error occurred while getting all users")
	}

	startCursor := pagination.EncodeCursor(urs[0].ID)
	endCursor := pagination.EncodeCursor(urs[len(urs)-1].ID)
	p := model.PageInfo{
		HasPreviousPage: hasPages[0],
		HasNextPage:     hasPages[1],
		StartCursor:     &startCursor,
		EndCursor:       &endCursor,
	}

	return urs, &p, nil
}
