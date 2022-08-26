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
	paginatedQuery, err := pagination.PaginateQuery(query, after, before, first, last, sort, order)
	if err != nil {
		return nil, err
	}

	q, args, err := paginatedQuery.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while getting all %s", table)
	}

	var c model.Connection[T]
	if err = db.Select(&c.Nodes, q, args...); err != nil || len(c.Nodes) == 0 {
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
	pageInfoQuery := pagination.GetPageInfoQuery(query, []string{(*c.Nodes[0]).GetID()}, []string{(*c.Nodes[len(c.Nodes)-1]).GetID()}, order)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while getting all %s", table)
	}

	q, args, err = pageInfoQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while getting all %s", table)
	}

	// Query HasNextPage and HasPrevious page fields od page info
	var hasPages []bool
	if err = db.Select(&hasPages, q, args...); err != nil {
		return nil, fmt.Errorf("an error occurred while getting all %s", table)
	}

	c.PageInfo = &model.PageInfo{
		HasPreviousPage: hasPages[0],
		HasNextPage:     hasPages[1],
		StartCursor:     &c.Edges[0].Cursor,
		EndCursor:       &c.Edges[len(c.Nodes)-1].Cursor,
	}

	return &c, nil
}

func getAllByKeys[T model.Node, U model.NodeWithKey[T]](db *sqlx.DB, query, subQuery sq.SelectBuilder, join string, table string, keys []string, after, before *string, first, last *int, sort string, order model.Order) ([]*model.Connection[T], []error) {
	// Build paginated query
	paginatedQuery, err := pagination.PaginateQuery(subQuery, after, before, first, last, sort, order)
	if err != nil {
		return nil, []error{err}
	}

	q, args, err := query.JoinClause(sq.Expr(join, paginatedQuery)).OrderBy(sort + " " + string(order)).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, []error{fmt.Errorf("an error occurred while getting all %s", table)}
	}

	var res []*U
	if err = db.Select(&res, q, args...); err != nil {
		return nil, []error{fmt.Errorf("an error occurred while getting all %s", table)}
	}

	ns := make(map[string][]*T)
	for _, r := range res {
		ns[(*r).GetKey()] = append(ns[(*r).GetKey()], (*r).GetNode())
	}

	var startIDs, endIDs []string
	for _, n := range ns {
		startIDs = append(startIDs, (*n[0]).GetID())
		endIDs = append(endIDs, (*n[len(n)-1]).GetID())
	}

	q, args, err = pagination.GetPageInfoQuery(query.JoinClause(sq.Expr(join, subQuery)), startIDs, endIDs, order).ToSql()
	if err != nil {
		return nil, []error{fmt.Errorf("an error occurred while getting all %s", table)}
	}

	hasPages := make([]bool, len(startIDs)+len(endIDs))
	if err := db.Select(&hasPages, q, args...); err != nil {
		return nil, []error{fmt.Errorf("an error occurred while getting all %s", table)}
	}

	hp, i := make(map[string][2]bool), 0
	for k := range ns {
		hp[k] = [2]bool{hasPages[i], hasPages[i+len(startIDs)]}
		i++
	}

	cs := make([]*model.Connection[T], len(keys))
	for i, key := range keys {
		if len(ns[key]) > 0 {
			var c model.Connection[T]
			c.Nodes = ns[key]

			c.Edges = make([]*model.Edge[T], len(c.Nodes))
			for j, n := range c.Nodes {
				c.Edges[j] = &model.Edge[T]{
					Cursor: pagination.EncodeCursor((*n).GetID()),
					Node:   n,
				}
			}

			c.PageInfo = &model.PageInfo{
				HasPreviousPage: hp[key][0],
				HasNextPage:     hp[key][1],
				StartCursor:     &c.Edges[0].Cursor,
				EndCursor:       &c.Edges[len(c.Nodes)-1].Cursor,
			}

			cs[i] = &c
		}
	}

	return cs, nil
}
