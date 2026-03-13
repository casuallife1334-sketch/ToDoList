package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(httpHandlers *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandlers,
	}
}

/*
Для запуска сервера нам нужно зарегистрировать какие-либо хендлеры
НО: у нас в программе достаточно сложные правила роутинга

Роутинг - по входящим параметрам (pattern, метод) понять какой именно хендлер
нужно вызвать, поэтому пакет net/http НЕ очень то и подходит под такую задачу
(нужно было бы писать ЧЕРЕСЧУР много условных ветвлений для проверки вызова нужного хендлера)

Выход - удобная внешняя библиотека github.com/gorilla/mux
*/
func (s *HTTPServer) StartServer() error {
	// Создаем роутер из пакета mux
	router := mux.NewRouter()

	// Задаем правила для вызова конкретного хендлера
	router.Path("/tasks").Methods("POST").HandlerFunc(s.httpHandlers.HandleCreateTask)
	router.Path("/tasks/{title}").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetTask)

	// Стоит обратить внимания на эти два роутера:
	// Роутер с query параметрами должен стоять СВЕРХУ (быть зарегистрирован раньше),
	// потому что поиск нужного роутера идет СВЕРХУ ВНИЗ
	// Иначе нам просто будут приходить все задачи, а не только выполненные
	router.Path("/tasks").Methods("GET").Queries("completed", "false").HandlerFunc(s.httpHandlers.HandleGetAllUncompletedTasks)
	router.Path("/tasks").Methods("GET").HandlerFunc(s.httpHandlers.HandleGetAllTasks)

	router.Path("/tasks/{title}").Methods("PATCH").HandlerFunc(s.httpHandlers.HandleCompleteTask)
	router.Path("/tasks/{title}").Methods("DELETE").HandlerFunc(s.httpHandlers.HandleDeleteTask)

	// Наполнили правилами роутинга роутеры и передали его в ListenAndServe
	// Под капотом это все правильно маршрутизируется

	// Дело в том, что функция ListenAndServe ВСЕГДА возвращает ошибку, поэтому
	// сделаем такой хак
	if err := http.ListenAndServe(":9091", router); err != nil {
		// Ошибка такого вида означает, что мы закрываем сервер (она будет при успешно запуске)
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
