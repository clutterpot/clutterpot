package pagination

import (
	"fmt"

	"github.com/clutterpot/clutterpot/model"

	sq "github.com/Masterminds/squirrel"
)

func PaginateQuery(query sq.SelectBuilder, after, before *string, first, last *int, sort string, order model.Order) (sq.SelectBuilder, error) {
	query, err := setAfterAndBefore(query, after, before, order)
	if err != nil {
		return query, err
	}
	query, err = setLimitSortAndOrder(query, first, last, sort, string(order))
	if err != nil {
		return query, err
	}

	return query, nil
}

func GetPageInfoQuery(query sq.SelectBuilder, startID, endID []string, order model.Order) sq.SelectBuilder {
	return query.Prefix("WITH q AS (").SuffixExpr(sq.Expr("), r AS (?), s AS (?) ?, LATERAL (?) b UNION ALL ?, LATERAL (?) b",
		sq.Select("id").From("q").Where(sq.Eq{"id": startID}),
		sq.Select("id").From("q").Where(sq.Eq{"id": endID}),
		sq.Select("b.*").From("r ra"),
		sq.Select("count(*) != 0").From("q").Where(fmt.Sprintf("id %s ra.id", flipGt(order == model.OrderAsc))).Limit(1),
		sq.Select("b.*").From("s sa"),
		sq.Select("count(*) != 0").From("q").Where(fmt.Sprintf("id %s sa.id", flipLt(order == model.OrderAsc))).Limit(1),
	)).PlaceholderFormat(sq.Dollar)
}

func setAfterAndBefore(query sq.SelectBuilder, after, before *string, order model.Order) (sq.SelectBuilder, error) {
	if after != nil {
		id, err := decodeCursor(*after)
		if err != nil {
			return query, fmt.Errorf("invalid 'after' cursor")
		}
		query = query.Where(fmt.Sprintf("id %s ?", flipLt(order == model.OrderAsc)), id)
	}
	if before != nil {
		id, err := decodeCursor(*before)
		if err != nil {
			return query, fmt.Errorf("invalid 'before' cursor")
		}
		query = query.Where(fmt.Sprintf("id %s ?", flipGt(order == model.OrderAsc)), id)
	}

	return query, nil
}

func setLimitSortAndOrder(query sq.SelectBuilder, first, last *int, sort, order string) (sq.SelectBuilder, error) {
	// "Strongly discouraged" - https://relay.dev/graphql/connections.htm#sel-FAJLFGBEBY22c
	if first != nil && last != nil {
		return query, fmt.Errorf("'first' and 'last' cannot be given at the same time")
	}
	if first != nil && *first > 0 && *first <= 50 {
		query = query.OrderBy(sort + " " + order).Limit(uint64(*first))
	} else if last != nil && *last > 0 && *last <= 50 {
		query = sq.Select("*").FromSelect(
			query.OrderBy(sort+" "+flipOrder(order)).Limit(uint64(*last)), "s",
		).OrderBy(sort + " " + order)
	} else {
		query = query.OrderBy(sort + " " + order).Limit(50)
	}

	return query, nil
}

func flipOrder(order string) string {
	if order == "ASC" {
		return "DESC"
	}

	return "ASC"
}

func flipLt(flip bool) string {
	if flip {
		return ">"
	}

	return "<"
}

func flipGt(flip bool) string {
	if flip {
		return "<"
	}

	return ">"
}
