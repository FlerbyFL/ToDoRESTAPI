# ToDo REST API

This repository contains a simple REST API for managing ToDo tasks. The API is designed to be lightweight and easy to use, allowing clients to create, read, update, and delete tasks using standard HTTP requests.

## Description
The project implements a backend for managing a task list. A client can send requests to add new tasks, view all tasks, get only incomplete tasks, mark a task as completed, and delete a task.

The API returns data in JSON format and follows basic REST principles.

## API Endpoints
- `GET /tasks` — get all tasks
- `GET /tasks?completed=false` — get uncompleted tasks
- `POST /tasks` — add a task
- `PATCH /tasks/:id` — update task status
- `DELETE /tasks/:id` — delete a task

## Используемые технологии
- Go 
- Gorilla Mux
- JSON

## Функции
- Add a task
- View all tasks
- View uncompleted tasks
- Complete a task
- Delete a task