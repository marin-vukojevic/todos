package todos

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type TodoHandler struct {
	repository *TodoRepository
}

func NewTodoHandler(repository *TodoRepository) *TodoHandler {
	return &TodoHandler{
		repository: repository,
	}
}

func (handler *TodoHandler) Index(w http.ResponseWriter, r *http.Request) {
	fetched, err := handler.repository.GetAllTodos(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := Index(fetched)
	component.Render(r.Context(), w)
}

func (handler *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	description := r.PostFormValue("description")
	todo, err := handler.repository.CreateTodo(r.Context(), description)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := TodoTableItem(todo)
	component.Render(r.Context(), w)
}

func (handler *TodoHandler) CompleteTodo(w http.ResponseWriter, r *http.Request) {
	uuidString := chi.URLParam(r, "todoUuid")
	parsed, err := uuid.Parse(uuidString)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	todo, err := handler.repository.MarkTodoAsDone(r.Context(), parsed)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := TodoTableItem(todo)
	component.Render(r.Context(), w)
}
