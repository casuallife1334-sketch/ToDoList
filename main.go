package main

import (
	"ToDoListNilchan/internal/core"
	"ToDoListNilchan/internal/todo/repository"
	"ToDoListNilchan/internal/todo/service"
	"ToDoListNilchan/internal/todo/transport"
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	logger, logFileClose, err := core.NewLogger("INFO")
	if err != nil {
		panic(err)
	}
	defer logFileClose()

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM)

	conn, err := repository.ConnectRepository(ctx)
	if err != nil {
		log.Fatal(err)
	}

	postgres := repository.NewRepository(conn)
	service := service.NewService(postgres)
	handlers := transport.NewHTTPHandlers(service, logger)
	server := transport.NewHTTPServer(handlers)

	if err := server.StartServer(); err != nil {
		log.Fatal(err)
	}
}
