-- name: CreateTodo :one
insert into todos (id, description, done)
values ($1, $2, $3)
returning *;

-- name: GetTodos :many
select * from todos;

-- name: MarkTodoAsDone :one
update todos
set done = true
where id = $1
returning *;
