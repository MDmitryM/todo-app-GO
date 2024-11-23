package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
}

type TodoList interface {
}

type ListItem interface {
}

type Repository struct {
	Authorization
	TodoList
	ListItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
