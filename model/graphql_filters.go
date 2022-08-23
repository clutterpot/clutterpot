package model

import (
	"time"

	sq "github.com/Masterminds/squirrel"
)

type ScalarFilter[T any] struct {
	And *ScalarFilter[T]
	Or  *ScalarFilter[T]
	Eq  *T
	Lt  *T
	Gt  *T
	Le  *T
	Ge  *T
	In  []*T
}

type StringFilter = ScalarFilter[string]
type IDFilter = ScalarFilter[string]
type TimeFilter = ScalarFilter[time.Time]
type Int64Filter = ScalarFilter[int64]

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
	} else if f.Le != nil {
		conj = append(conj, sq.LtOrEq{column: *f.Le})
	} else if f.Ge != nil {
		conj = append(conj, sq.GtOrEq{column: *f.Ge})
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
