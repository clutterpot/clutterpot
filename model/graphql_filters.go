package model

import (
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Filter interface {
	GetConj() sq.And
}

type ScalarFilter[T any] struct {
	Eq  *T
	Lt  *T
	Gt  *T
	Leq *T
	Geq *T
	In  []*T
}

type StringFilter = ScalarFilter[string]
type IDFilter = ScalarFilter[string]
type TimeFilter = ScalarFilter[time.Time]

func (f *ScalarFilter[T]) GetConj(and sq.And, column string) sq.And {
	if f == nil {
		return and
	}
	if f.Eq != nil {
		return append(and, sq.Eq{column: *f.Eq})
	}
	if f.Lt != nil {
		return append(and, sq.Lt{column: *f.Lt})
	}
	if f.Gt != nil {
		return append(and, sq.Gt{column: *f.Gt})
	}
	if f.Leq != nil {
		return append(and, sq.LtOrEq{column: *f.Leq})
	}
	if f.Geq != nil {
		return append(and, sq.GtOrEq{column: *f.Geq})
	}
	if f.In != nil {
		return append(and, sq.Eq{column: f.In})
	}

	return and
}
