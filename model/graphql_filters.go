package model

import (
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Filter interface {
	GetConj() []sq.Sqlizer
}

type ScalarFilter[T any] struct {
	And *ScalarFilter[T]
	Or  *ScalarFilter[T]
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

func (f *ScalarFilter[T]) GetConj(conj sq.And, column string) sq.And {
	if f == nil {
		return conj
	}

	if f.Eq != nil {
		conj = append(conj, sq.Eq{column: *f.Eq})
	} else if f.Lt != nil {
		conj = append(conj, sq.Lt{column: *f.Lt})
	} else if f.Gt != nil {
		conj = append(conj, sq.Gt{column: *f.Gt})
	} else if f.Leq != nil {
		conj = append(conj, sq.LtOrEq{column: *f.Leq})
	} else if f.Geq != nil {
		conj = append(conj, sq.GtOrEq{column: *f.Geq})
	} else if f.In != nil {
		conj = append(conj, sq.Eq{column: f.In})
	}

	if f.And != nil {
		conj = append(conj, f.And.GetConj(sq.And{}, column)...)
	}
	if f.Or != nil {
		conj = sq.And{sq.Or{conj, f.Or.GetConj(sq.And{}, column)}}
	}

	return conj
}
