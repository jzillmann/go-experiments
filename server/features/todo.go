package features

import (
	"context"
	"sync"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type Todo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoReference struct {
	Id int `path:"id" minimum:"1" doc:"ID of the Todo" example:"23"`
}

type TodoResponse struct {
	Body *Todo
}

type TodosResponse struct {
	Body []*Todo
}

var (
	todos      = make([]*Todo, 0)
	currentID  = 0
	todosMutex sync.Mutex
)

func TodoRoutes(api huma.API, basePath string) {

	// GET ALL
	huma.Get(
		api,
		basePath, func(ctx context.Context, input *struct{}) (*TodosResponse, error) {
			todosMutex.Lock()
			defer todosMutex.Unlock()

			time.Sleep(1 * time.Second)
			resp := &TodosResponse{
				Body: todos,
			}
			return resp, nil
		})

	// GET ONE
	huma.Get(api, basePath+"/{id}", func(ctx context.Context, input *TodoReference) (*TodoResponse, error) {
		todosMutex.Lock()
		defer todosMutex.Unlock()
		index := findIndex(input.Id)
		if index == -1 {
			return nil, huma.Error404NotFound("todo not found")
		}
		resp := &TodoResponse{
			Body: todos[index],
		}
		return resp, nil
	})

	// CREATE ONE
	type CreateInput struct {
		Body struct {
			Title string `json:"title" minLength:"3" doc:"The todo text" example:"Buy sugar"`
		}
	}
	huma.Post(api, basePath, func(ctx context.Context, input *CreateInput) (*TodoResponse, error) {
		todosMutex.Lock()
		defer todosMutex.Unlock()

		currentID++
		todo := &Todo{
			Id:    currentID,
			Title: input.Body.Title,
		}
		todos = append(todos, todo)
		return &TodoResponse{
			Body: todo,
		}, nil
	})

	// UPDATE ONE
	type UpdateInput struct {
		TodoReference
		Body struct {
			Title     *string `json:"title,omitempty" minLength:"3"`
			Completed *bool   `json:"completed,omitempty"`
		}
	}
	huma.Patch(api, basePath+"/{id}", func(ctx context.Context, input *UpdateInput) (*TodoResponse, error) {
		todosMutex.Lock()
		defer todosMutex.Unlock()

		index := findIndex(input.Id)
		if index == -1 {
			return nil, huma.Error404NotFound("todo not found")
		}

		if input.Body.Title != nil {
			todos[index].Title = *input.Body.Title
		}
		if input.Body.Completed != nil {
			todos[index].Completed = *input.Body.Completed
		}
		return &TodoResponse{
			Body: todos[index],
		}, nil
	})

	// DELETE ONE
	huma.Delete(api, basePath+"/{id}", func(ctx context.Context, input *TodoReference) (*struct{}, error) {
		todosMutex.Lock()
		defer todosMutex.Unlock()
		index := findIndex(input.Id)
		if index == -1 {
			return nil, huma.Error404NotFound("todo not found")
		}
		time.Sleep(500 * time.Millisecond)
		todos = append(todos[:index], todos[index+1:]...)
		return nil, nil
	})
}

func findIndex(id int) int {
	for index, todo := range todos {
		if todo.Id == id {
			return index
		}
	}
	return -1
}
