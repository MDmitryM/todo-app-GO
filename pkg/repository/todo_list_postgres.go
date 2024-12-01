package repository

import (
	"fmt"
	"strings"

	"github.com/MDmitryM/todo-app-GO"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (t *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return 0, err
	}

	var listId int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&listId); err != nil {
		return 0, err
	}

	createUsersListsQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", userListsTable)
	_, err = tx.Exec(createUsersListsQuery, userId, listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return listId, tx.Commit()
}

func (t *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList

	queryGetAllLists := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id=$1",
		todoListsTable, userListsTable)
	err := t.db.Select(&lists, queryGetAllLists, userId)

	return lists, err
}

func (t *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	queryListQuery := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id=$1 AND ul.list_id=$2",
		todoListsTable, userListsTable)
	err := t.db.Get(&list, queryListQuery, userId, listId)

	return list, err
}

func (t *TodoListPostgres) Delete(userId, listId int) error {

	queryListQuery := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id =$2",
		todoListsTable, userListsTable)
	_, err := t.db.Exec(queryListQuery, userId, listId)

	return err
}

func (t *TodoListPostgres) Update(userId, listId int, input todo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	updateQuery := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListsTable, setQuery, userListsTable, argId, argId+1)

	args = append(args, listId, userId)
	logrus.Debugf("updateQuery: %s", updateQuery)
	logrus.Debugf("args: %s", args)

	_, err := t.db.Exec(updateQuery, args...)

	return err
}
