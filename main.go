package main

import (
	"fmt"

	"restapi/http"
	"restapi/todo"
)

func main() {
	todolist := todo.NewList()
	httpHandlers := http.NewHTTPHandlers(todolist)
	httpSrever := http.NewHTTPServer(httpHandlers)

	if err := httpSrever.StartServer(); err != nil {
		fmt.Println("failde to start server", err)
	}

}
