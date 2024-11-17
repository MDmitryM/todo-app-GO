package repository

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

func NewRepository() *Repository {
	return &Repository{}
}
