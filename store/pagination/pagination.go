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

func GetPageInfoQuery(query sq.SelectBuilder, startID, endID string, order model.Order) sq.SelectBuilder {
	return sq.Select("count(*) != 0").From("q").Prefix("WITH q AS (?) (", query).
		Where(flipGt(sq.Gt{"id": startID}, order == model.OrderAsc)).Limit(1).Suffix(
		") UNION ALL (?)", sq.Select("count(*) != 0").From("q").
			Where(flipLt(sq.Lt{"id": endID}, order == model.OrderAsc)).Limit(1),
	)
}

func setAfterAndBefore(query sq.SelectBuilder, after, before *string, order model.Order) (sq.SelectBuilder, error) {
	if after != nil {
		id, err := decodeCursor(*after)
		if err != nil {
			return query, fmt.Errorf("invalid 'after' cursor")
		}
		query = query.Where(flipLt(sq.Lt{"id": id}, order == model.OrderAsc))
	}
	if before != nil {
		id, err := decodeCursor(*before)
		if err != nil {
			return query, fmt.Errorf("invalid 'before' cursor")
		}
		query = query.Where(flipGt(sq.Gt{"id": id}, order == model.OrderAsc))
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

func flipLt(operator sq.Lt, flip bool) any {
	if flip {
		return sq.Gt{"id": operator["id"]}
	}

	return operator
}

func flipGt(operator sq.Gt, flip bool) any {
	if flip {
		return sq.Lt{"id": operator["id"]}
	}

	return operator
}
