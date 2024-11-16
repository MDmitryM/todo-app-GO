package main

import (
	"log"

	"github.com/MDmitryM/todo-app-GO"
	"github.com/MDmitryM/todo-app-GO/pkg/handler"
)

func main() {
	handler := new(handler.Handler)
	srv := new(todo.Server)

	if err := srv.Run("8000", handler.InitRoutes()); err != nil {
		log.Fatalf("server run error! %s", err.Error())
	}
}
