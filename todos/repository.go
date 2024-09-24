package todos

import (
	"context"

	"github.com/google/uuid"
	"github.com/marin-vukojevic/todos/generated/database"
)

type TodoRepository struct {
	queries *database.Queries
}

func NewTodoRepository(queries *database.Queries) *TodoRepository {
	return &TodoRepository{
		queries: queries,
	}
}

func (repo *TodoRepository) GetAllTodos(ctx context.Context) ([]Todo, error) {
	dbTodos, err := repo.queries.GetTodos(ctx)
	if err != nil {
		return []Todo{}, err
	}

	return convertMultipleFromDb(dbTodos), nil
}

func (repo *TodoRepository) CreateTodo(ctx context.Context, description string) (Todo, error) {
	dbTodo, err := repo.queries.CreateTodo(ctx, database.CreateTodoParams{ID: uuid.New(), Description: description, Done: false})
	if err != nil {
		return Todo{}, err
	}

	return convertSingleFromDb(dbTodo), nil
}

func (repo *TodoRepository) MarkTodoAsDone(ctx context.Context, id uuid.UUID) (Todo, error) {
	dbTodo, err := repo.queries.MarkTodoAsDone(ctx, id)
	if err != nil {
		return Todo{}, err
	}

	return convertSingleFromDb(dbTodo), nil
}

func convertMultipleFromDb(dbTodos []database.Todo) []Todo {
	converted := make([]Todo, 0, len(dbTodos))
	for _, dbTodo := range dbTodos {
		converted = append(converted, Todo{Uuid: dbTodo.ID, Description: dbTodo.Description, Done: dbTodo.Done})
	}
	return converted
}

func convertSingleFromDb(dbTodo database.Todo) Todo {
	return Todo{Uuid: dbTodo.ID, Description: dbTodo.Description, Done: dbTodo.Done}
}
