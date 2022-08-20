package store

import (
	"database/sql"
	"fmt"

	"github.com/clutterpot/clutterpot/model"
	"github.com/clutterpot/clutterpot/store/pagination"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	User    *userStore
	File    *fileStore
	Tag     *tagStore
	Session *sessionStore
}

func New(db *sqlx.DB) *Store {
	return &Store{
		User:    newUserStore(db),
		File:    newFileStore(db),
		Tag:     newTagStore(db),
		Session: newSessionStore(db),
	}
}

func getAll[T model.Node](db *sqlx.DB, query sq.SelectBuilder, table string, after, before *string, first, last *int, sort string, order model.Order) (*model.Connection[T], error) {
	// Build paginated query
	paginatedQuery, args, err := pagination.PaginateQuery(query, after, before, first, last, sort, order)
	if err != nil {
		return nil, err
	}

	var c model.Connection[T]
	if err = db.Select(&c.Nodes, paginatedQuery, args...); err != nil || len(c.Nodes) == 0 {
		if err == sql.ErrNoRows || len(c.Nodes) == 0 {
			return nil, fmt.Errorf("%s not found", table)
		}

		return nil, fmt.Errorf("an error occurred while getting all %s", table)
	}

	c.Edges = make([]*model.Edge[T], len(c.Nodes))
	for i, n := range c.Nodes {
		c.Edges[i] = &model.Edge[T]{
			Cursor: pagination.EncodeCursor((*n).GetID()),
			Node:   n,
		}
	}

	// Build page info query
	pageInfoQuery, args, err := pagination.GetPageInfoQuery(query, (*c.Nodes[0]).GetID(), (*c.Nodes[len(c.Nodes)-1]).GetID(), order)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while getting all %s", table)
	}

	// Query HasNextPage and HasPrevious page fields od page info
	var hasPages []bool
	if err = db.Select(&hasPages, pageInfoQuery, args...); err != nil {
		return nil, fmt.Errorf("an error occurred while getting all %s", table)
	}

	startCursor := pagination.EncodeCursor((*c.Nodes[0]).GetID())
	endCursor := pagination.EncodeCursor((*c.Nodes[len(c.Nodes)-1]).GetID())
	c.PageInfo = &model.PageInfo{
		HasPreviousPage: hasPages[0],
		HasNextPage:     hasPages[1],
		StartCursor:     &startCursor,
		EndCursor:       &endCursor,
	}

	return &c, nil
}
