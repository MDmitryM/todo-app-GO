package repository

import (
	"github.com/MDmitryM/todo-app-GO"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
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
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
