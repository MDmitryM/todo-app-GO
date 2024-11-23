package service

import (
	"github.com/MDmitryM/todo-app-GO"
	"github.com/MDmitryM/todo-app-GO/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
}

type TodoList interface {
}

type ListItem interface {
}

type Service struct {
	Authorization
	TodoList
	ListItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos),
	}
}
