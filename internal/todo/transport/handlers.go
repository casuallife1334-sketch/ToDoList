package transport

import (
	"ToDoListNilchan/internal/core"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type TaskService interface {
	CreateTask(ctx context.Context, title, description string) (core.TaskDomain, error)
	GetTask(ctx context.Context, id int) (core.TaskDomain, error)
	GetAllTasks(ctx context.Context) ([]core.TaskDomain, error)
	GetAllUncompletedTasks(ctx context.Context) ([]core.TaskDomain, error)
	CompleteTask(ctx context.Context, id int) (core.TaskDomain, error)
	UncompleteTask(ctx context.Context, id int) (core.TaskDomain, error)
	DeleteTask(ctx context.Context, id int) error
}

type HTTPHandlers struct {
	taskService TaskService
}

func NewHTTPHandlers(taskService TaskService) *HTTPHandlers {
	return &HTTPHandlers{
		taskService: taskService,
	}
}

/*
Контракт метода HandleCreateTask, удовлетворяющий парадигме REST API
Входящая информация HandleCreateTask:
pattern: /tasks
method: POST
info: JSON in HTTP request body

Ответная информация:
succeed:
  - status code: 201 created
  - response body: JSON represent created task | JSON представляет созданную задачу

failed:
  - status code: 400, 409, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO

	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := taskDTO.ValidateForCreate(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	task, err := h.taskService.CreateTask(r.Context(), taskDTO.Title, taskDTO.Description)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(task, "", "	")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP-response:", err)
		return
	}
}

/*
Контракт метода HandleGetTask, удовлетворяющий парадигме REST API
Входящая информация HandleGetTask:
pattern: /tasks/(title)
method: GET
info: pattern

Ответная информация:
succeed:
  - status code: 200 OK
  - response body: JSON represent found task

failed:
  - status code: 400, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	task, err := h.taskService.GetTask(r.Context(), id)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, core.ErrNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(task, "", "	")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP-response:", err)
		return
	}
}

/*
Контракт метода HandleGetAllTasks, удовлетворяющий парадигме REST API
Входящая информация HandleGetAllTasks:
pattern: /tasks
method: GET
info: -

Ответная информация:
succeed:
  - status code: 200 OK
  - response body: JSON represent found tasks

failed:
  - status code: 400, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskService.GetAllTasks(r.Context())
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(tasks, "", "	")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP-response:", err)
		return
	}
}

/*
Контракт метода HandleGetAllUncompletedTasks, удовлетворяющий парадигме REST API
Входящая информация HandleGetAllUncompletedTasks:
pattern: /tasks?completed=true
method: GET
info: query parameters

Ответная информация:
succeed:
  - status code: 200 OK
  - response body: JSON represent found tasks

failed:
  - status code: 400, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleGetAllUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	uncompletedTasks, err := h.taskService.GetAllUncompletedTasks(r.Context())
	b, err := json.MarshalIndent(uncompletedTasks, "", "	")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP-response:", err)
		return
	}
}

/*
Контракт метода HandleCompleteTask, удовлетворяющий парадигме REST API
Входящая информация HandleCompleteTask:
pattern: /tasks/(title)
method: PATCH
info: pattern + JSON in HTTP request body

Ответная информация:
succeed:
  - status code: 200 OK
  - response body: JSON represent changed task

failed:
  - status code: 400, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var completeDTO CompleteTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	var (
		updatedTask core.TaskDomain
		err         error
	)

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if completeDTO.Complete {
		updatedTask, err = h.taskService.CompleteTask(r.Context(), id)
	} else {
		updatedTask, err = h.taskService.UncompleteTask(r.Context(), id)
	}

	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, core.ErrNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(updatedTask, "", "    ")
	if err != nil {
		// Так делать очень не желательно, но мы пока не изучили логирование
		// поэтому пока так
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("fail to write HTTP-response:", err)
		return
	}
}

/*
Контракт метода HandleDeleteTask, удовлетворяющий парадигме REST API
Входящая информация HandleDeleteTask:
pattern: /tasks/(title)
method: DELETE
info: pattern

Ответная информация:
succeed:
  - status code: 204 No Content
  - response body: -

failed:
  - status code: 400, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := h.taskService.DeleteTask(r.Context(), id); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, core.ErrNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
