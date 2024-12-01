package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MDmitryM/todo-app-GO"
	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (t *TodoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListItemQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (t *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem

	queryGetAllItems := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id INNER JOIN %s ul on ul.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2",
		todoItemsTable, listsItemsTable, userListsTable)

	if err := t.db.Select(&items, queryGetAllItems, listId, userId); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New("no items found in list")
	}

	return items, nil
}

func (t *TodoItemPostgres) GetById(userId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem

	queryGetById := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id INNER JOIN %s ul on ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2",
		todoItemsTable, listsItemsTable, userListsTable)

	if err := t.db.Get(&item, queryGetById, itemId, userId); err != nil {
		return item, err
	}

	// if item.Id == 0 {
	// 	return item, errors.New("no items found in list")
	// }

	return item, nil
}

func (t *TodoItemPostgres) Delete(userId, itemId int) error {
	deleteQuery := fmt.Sprintf("DELETE FROM %s ti USING %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2",
		todoItemsTable, listsItemsTable, userListsTable)
	del, err := t.db.Exec(deleteQuery, userId, itemId)
	if err != nil {
		return err
	}

	result, err := del.RowsAffected()

	if err != nil {
		return err
	}

	if result == 0 {
		err = errors.New("item is not found to delete")
	}

	return err
}

func (t *TodoItemPostgres) Update(userId, itemId int, input todo.UpdateItemInput) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	updateQuery := fmt.Sprintf("UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d",
		todoItemsTable, setQuery, listsItemsTable, userListsTable, argId, argId+1)

	args = append(args, userId, itemId)

	upd, err := t.db.Exec(updateQuery, args...)
	if err != nil {
		return err
	}

	result, err := upd.RowsAffected()
	if err != nil {
		return err
	}

	if result == 0 {
		err = errors.New("item does not exist")
	}

	return err
}
