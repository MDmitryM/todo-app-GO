package service

import (
	"github.com/MDmitryM/todo-app-GO"
	"github.com/MDmitryM/todo-app-GO/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (t *TodoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) {
	_, err := t.listRepo.GetById(userId, listId)
	if err != nil {
		//lsit does not exists or does noe belongs to user
		return 0, err
	}

	return t.repo.Create(listId, item)
}

func (t *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	return t.repo.GetAll(userId, listId)
}

func (t *TodoItemService) GetById(userId, itemId int) (todo.TodoItem, error) {
	return t.repo.GetById(userId, itemId)
}

func (t *TodoItemService) Delete(userId, itemId int) error {
	return t.repo.Delete(userId, itemId)
}

func (t *TodoItemService) Update(userId, itemId int, input todo.UpdateItemInput) error {
	return t.repo.Update(userId, itemId, input)
}
