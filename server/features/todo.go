package features

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Todo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Bind implements the render.Binder interface
func (t *Todo) Bind(r *http.Request) error {
	// As a simple validation example: ensure that the todo's Title is not empty
	if t.Title == "" {
		return errors.New("title of the todo cannot be empty")
	}
	return nil
}

var (
	todos      = make([]*Todo, 0)
	currentID  = 0
	todosMutex sync.Mutex
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", CreateTodo)
	router.Put("/{id}", UpdateTodo)
	router.Delete("/{id}", DeleteTodo)
	router.Get("/", GetTodos)
	return router
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	todosMutex.Lock()
	defer todosMutex.Unlock()

	var todo Todo
	if err := render.Bind(r, &todo); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	currentID++
	todo.Id = currentID
	todos = append(todos, &todo)
	render.Status(r, http.StatusCreated)
	render.Render(w, r, &TodoRenderer{Todo: &todo})
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todosMutex.Lock()
	defer todosMutex.Unlock()

	todoId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(todoId)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	index := findIndex(id)
	if index == -1 {
		render.Render(w, r, ErrNotFound())
		return
	}
	time.Sleep(500 * time.Millisecond)
	todos = append(todos[:index], todos[index+1:]...)
	render.Status(r, http.StatusNoContent)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todosMutex.Lock()
	defer todosMutex.Unlock()

	todoId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(todoId)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	index := findIndex(id)
	if index == -1 {
		render.Render(w, r, ErrNotFound())
		return
	}

	var updateData struct {
		Title     *string `json:"title,omitempty"`
		Completed *bool   `json:"completed,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "JSON decode error")
		return
	}

	if updateData.Title != nil {
		todos[index].Title = *updateData.Title
	}
	if updateData.Completed != nil {
		todos[index].Completed = *updateData.Completed
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, todos[index])
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	todosMutex.Lock()
	defer todosMutex.Unlock()

	time.Sleep(1 * time.Second)

	if err := render.RenderList(w, r, NewTodoListResponse(todos)); err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

type TodoRenderer struct {
	*Todo
}

func (t *TodoRenderer) Render(w http.ResponseWriter, r *http.Request) error {
	// Any transformations or additional data can be set here
	return nil
}

func NewTodoListResponse(todos []*Todo) []render.Renderer {
	list := []render.Renderer{}
	for _, todo := range todos {
		list = append(list, &TodoRenderer{Todo: todo})
	}
	return list
}

func findIndex(id int) int {
	for index, todo := range todos {
		if todo.Id == id {
			return index
		}
	}
	return -1
}
