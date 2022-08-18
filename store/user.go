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
		return nil, fmt.Errorf("an error occurred while creating a user")
	}

	return &u, nil
}

func (us *userStore) Update(id string, input model.UserUpdateInput) (*model.User, error) {
	var u model.User

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

	if err := query.Where("id = ?", id).
		Suffix("RETURNING id, username, email, kind, display_name, bio, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).RunWith(us.db).QueryRow().
		Scan(&u.ID, &u.Username, &u.Email, &u.Kind, &u.DisplayName, &u.Bio, &u.CreatedAt, &u.UpdatedAt); err != nil {
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
	var u model.User

	if err := sq.Select("id", "username", "email", "kind", "display_name", "bio", "created_at", "updated_at").
		From("users").Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).RunWith(us.db).QueryRow().
		Scan(&u.ID, &u.Username, &u.Email, &u.Kind, &u.DisplayName, &u.Bio, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("an error occurred while getting a user by id")
	}

	return &u, nil
}

// Only for auth purposes
func (us *userStore) GetByEmail(email string) (*model.AuthUser, error) {
	var u model.AuthUser

	if err := sq.Select("id", "username", "email", "password", "kind", "display_name", "bio", "created_at", "updated_at").
		From("users").Where("email = ?", email).
		PlaceholderFormat(sq.Dollar).RunWith(us.db).QueryRow().
		Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.Kind, &u.DisplayName, &u.Bio, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid email or password")
		}

		return nil, fmt.Errorf("an error occurred while getting a user by email")
	}

	return &u, nil
}

func (us *userStore) GetAll(after, before *string, first, last *int, sort *model.UserSort, order *model.Order) ([]*model.User, model.PageInfo, error) {
	var urs []*model.User
	var p model.PageInfo

	// Base query
	query := sq.Select("id", "username", "email", "kind", "display_name", "bio", "created_at", "updated_at").
		From("users")

	// Build paginated query
	paginatedQuery, err := pagination.PaginateQuery(query, after, before, first, last, string(*sort), *order)
	if err != nil {
		return nil, p, err
	}

	// Query all users (with pagination settings)
	rows, err := paginatedQuery.PlaceholderFormat(sq.Dollar).RunWith(us.db).Query()
	if err != nil {
		// return nil, fmt.Errorf("an error occurred while getting all users")
		return nil, p, err
	}
	defer rows.Close()

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Kind, &u.DisplayName, &u.Bio, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, p, fmt.Errorf("an error occurred while getting all users")
		}
		urs = append(urs, &u)
	}

	if len(urs) == 0 {
		return nil, p, nil
	}

	// Query HasNextPage and HasPrevious page fields od page info
	rows, err = pagination.GetPageInfoQuery(query, urs[0].ID, urs[len(urs)-1].ID, *order).
		PlaceholderFormat(sq.Dollar).RunWith(us.db).Query()
	if err != nil {
		return nil, p, err /*fmt.Errorf("an error occurred while getting all users")*/
	}

	rows.Next()
	rows.Scan(&p.HasPreviousPage)
	rows.Next()
	rows.Scan(&p.HasNextPage)

	// Set start and end cursors in page info
	startCursor := pagination.EncodeCursor(urs[0].ID)
	endCursor := pagination.EncodeCursor(urs[len(urs)-1].ID)
	p.StartCursor = &startCursor
	p.EndCursor = &endCursor

	return urs, p, nil
}
