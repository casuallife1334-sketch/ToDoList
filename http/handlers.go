package http

import (
	"ToDoListNilchan/todo"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	todoList *todo.List
}

func NewHTTPHandlers(todoList *todo.List) *HTTPHandlers {
	return &HTTPHandlers{
		todoList: todoList,
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
	// Здесь, приняв запрос, мы бы могли обращаться к структуре Task из пакета todo,
	// НО там есть куча полей, которые нам не нужны,
	// тк по сети мы принимаем только заголовок и описание

	// Создадим новую структуру DTO (data transfer object)
	var taskDTO TaskDTO

	// Читаем из тела запроса, если ошибка - обрабатываем
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		// Передаем в теле ответа ошибку таким образом, сразу со статус-кодом
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	// Проверяем - все ли провалидировалось?
	if err := taskDTO.ValidateForCreate(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	todoTask := todo.NewTask(taskDTO.Title, taskDTO.Description)
	if err := h.todoList.AddTask(todoTask); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		if errors.Is(err, todo.ErrTaskAlreadyExists) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(todoTask, "", "	")
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
	// Таким образом получаем мапу из паттерна
	// где ключ - {title}, а значение - то, что написал клиент

	// На наличие ключа можно не проверять, тк если вызвался этот хендлер
	// то клиент ПОЛЮБОМУ должен был передать что-то в эту мапу
	title := mux.Vars(r)["title"]

	task, err := h.todoList.GetTask(title)
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
	tasks := h.todoList.ListTask()
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
	uncompletedTasks := h.todoList.ListUncompletedTasks()
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
func (h *HTTPHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	var completeDTO CompleteTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	title := mux.Vars(r)["title"]

	var (
		// Создаем переменные за рамками условных ветвлений, чтобы не дублировать код
		changedTask todo.Task
		err         error
	)
	if completeDTO.Complete {
		// Будем прям из метода CompleteTask возвращать задачу, чтобы ее вывести в
		// HTTP ответе
		changedTask, err = h.todoList.CompleteTask(title)
	} else {
		changedTask, err = h.todoList.UncompleteTask(title)
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

	b, err := json.MarshalIndent(changedTask, "", "    ")
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
	title := mux.Vars(r)["title"]

	if err := h.todoList.DeleteTask(title); err != nil {
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
