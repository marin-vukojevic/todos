// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: todos.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createTodo = `-- name: CreateTodo :one
insert into todos (id, description, done)
values ($1, $2, $3)
returning id, description, done
`

type CreateTodoParams struct {
	ID          uuid.UUID
	Description string
	Done        bool
}

func (q *Queries) CreateTodo(ctx context.Context, arg CreateTodoParams) (Todo, error) {
	row := q.db.QueryRowContext(ctx, createTodo, arg.ID, arg.Description, arg.Done)
	var i Todo
	err := row.Scan(&i.ID, &i.Description, &i.Done)
	return i, err
}

const getTodos = `-- name: GetTodos :many
select id, description, done from todos
`

func (q *Queries) GetTodos(ctx context.Context) ([]Todo, error) {
	rows, err := q.db.QueryContext(ctx, getTodos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Todo
	for rows.Next() {
		var i Todo
		if err := rows.Scan(&i.ID, &i.Description, &i.Done); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markTodoAsDone = `-- name: MarkTodoAsDone :one
update todos
set done = true
where id = $1
returning id, description, done
`

func (q *Queries) MarkTodoAsDone(ctx context.Context, id uuid.UUID) (Todo, error) {
	row := q.db.QueryRowContext(ctx, markTodoAsDone, id)
	var i Todo
	err := row.Scan(&i.ID, &i.Description, &i.Done)
	return i, err
}