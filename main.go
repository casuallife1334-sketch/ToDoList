package main

import (
	"ToDoListNilchan/http"
	"ToDoListNilchan/todo"
	"fmt"
)

func main() {
	todoList := todo.NewList()
	httpHandlers := http.NewHTTPHandlers(todoList)
	httpServer := http.NewHTTPServer(httpHandlers)

	if err := httpServer.StartServer(); err != nil {
		fmt.Println("failed to start HTTP server:", err)
	}
}
