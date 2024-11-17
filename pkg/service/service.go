package service

import "github.com/MDmitryM/todo-app-GO/pkg/repository"

type Authorization interface {
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
	return &Service{}
}
