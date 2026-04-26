package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"restapi/todo"
)

type HTTPHandlers struct {
	todolist *todo.List
}

func NewHTTPHandlers(todolist *todo.List) *HTTPHandlers {
	return &HTTPHandlers{
		todolist: todolist,
	}
}

/*
pattern: /tasks
method: POST
info: 	 JSON in HTTP body

succeed:
	- status code: 201
	- response body: JSON with created task
failed:
	- status code: 400, 409, 500, ...
	- response body: JSON with error message + time

*/

func (h *HTTPHandlers) HandlerCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errorDTO := ErrorDTO{
			Message: "invalid request body",
			Time:    time.Now(),
		}

		http.Error(w, errorDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := taskDTO.ValidateForCreate(); err != nil {
		errorDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errorDTO.ToString(), http.StatusBadRequest)
		return
	}

	todoTask := todo.NewTask(taskDTO.Title, taskDTO.Description)
	if err := h.todolist.AddTask(todoTask); err != nil {
		errorDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, todo.ErrTaskAlreadyExists) {
			http.Error(w, errorDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errorDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(todoTask, "", "   ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: GET
info: 	 pattern

succeed:
  - status code: 201
  - response body: JSON representing found task

failed:
  - status code: 400, 409, 500, ...
  - response body: JSON with error message + time
*/
func (h *HTTPHandlers) HandlerGetTasks(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	task, err := h.todolist.GetTask(title)
	if err != nil {
		errorDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errorDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errorDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(task, "", "   ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks
method: GET
info: 	 pattern

succeed:
  - status code: 201
  - response body: JSON representing found task

failed:
  - status code: 400, 409, 500, ...
  - response body: JSON with error message + time
*/
func (h *HTTPHandlers) HandlerGetALLCompletedTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.todolist.ListTasks()
	b, err := json.MarshalIndent(tasks, "", "   ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks?completed=false
method: GET
info: 	query parameter

succeed:
  - status code: 201
  - response body: JSON representing found task

failed:
  - status code: 400, 409, 500, ...
  - response body: JSON with error message + time
*/
func (h *HTTPHandlers) HandlerGetALLNotCompletedTasks(w http.ResponseWriter, r *http.Request) {
	notcompletestasks := h.todolist.ListNotCompletedTasks()
	b, err := json.MarshalIndent(notcompletestasks, "", "   ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: PATCH
info: 	 pattern + JSON in HTTP body

succeed:
  - status code: 201
  - response body: JSON representing changed task

failed:
  - status code: 400, 409, 500, ...
  - response body: JSON with error message + time
*/
func (h *HTTPHandlers) HandlerCompleteTask(w http.ResponseWriter, r *http.Request) {
	var completeDTO CompleteTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := ErrorDTO{
			Message: "invalid request body",
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	title := mux.Vars(r)["title"]

	var (
		changedTask todo.Task
		err         error
	)

	if completeDTO.Complete {
		changedTask, err = h.todolist.CompleteTask(title)
	} else {
		changedTask, err = h.todolist.UncompleteTask(title)
	}

	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(changedTask, "", "   ")
	if err != nil {
		panic(err)
	}

	_, err = w.Write(b)
	if err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: DELETE
info: 	 pattern

succeed:
  - status code: 204 No content
  - response body: -

failed:
  - status code: 400, 409, 500, ...
  - response body: JSON with error message + time
*/
func (h *HTTPHandlers) HandlerDeleteTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	if err := h.todolist.DeleteTask(title); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
