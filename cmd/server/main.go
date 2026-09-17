package main

import (
	"fmt"

	"restapi/internal/http"
	"restapi/internal/todo"
)

func main() {
	todolist := todo.NewList()
	httpHandlers := http.NewHTTPHandlers(todolist)
	httpSrever := http.NewHTTPServer(httpHandlers)

	if err := httpSrever.StartServer(); err != nil {
		fmt.Println("failed to start server", err)
	}

}
