package model

type Connection[T any] struct {
	Edges    []*Edge[T] `json:"edges"`
	Nodes    []*T       `json:"nodes"`
	PageInfo *PageInfo  `json:"pageInfo"`
}

type Edge[T any] struct {
	Cursor string `json:"cursor"`
	Node   *T     `json:"node"`
}
