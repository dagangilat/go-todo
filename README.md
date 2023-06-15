# Todo App

This is a simple Todo app server built with Go. It allows you to manage tasks by providing various API endpoints.
Data is persisted in a MySQL database, which need to be installed and configured before running the application.

## Requirements

- Go 1.16+
- MySQL 8+

## API Endpoints

- `GET /tasks`: Retrieves all tasks.
- `POST /tasks`: Adds a new task.
- `PUT /tasks/{id}`: Updates an existing task.
- `DELETE /tasks/{id}`: Deletes a task.
- `GET /tasks/{id}`: Retrieves a single task by its ID.

## Usage

1. Make sure you have a MySQL database set up.
2. Update the database connection credentials in the code.
3. Build and run the Go application.

## Configuration

You can enable/disable debug logging by setting the `DEBUG` variable in the code.

## Dependencies

This application relies on the following dependencies:

- `github.com/gorilla/mux`: A powerful URL router and dispatcher for Go.
- `github.com/go-sql-driver/mysql`: A MySQL driver for Go's database/sql package.

Make sure to install these dependencies using Go's package manager before running the application.

## Examples

Here are some examples of how to use the API endpoints with cURL:

### Get all tasks

```bash
curl -X GET http://localhost:8000/tasks
```

### Add a new task

```bash
curl -X POST -H "Content-Type: application/json" -d '{"text":"Task 1", "completed": false}' http://localhost:8000/tasks
```

### Update an existing task

```bash
curl -X PUT -H "Content-Type: application/json" -d '{"text":"Updated Task", "completed": true}' http://localhost:8000/tasks/{id}
```

### Delete a task

```bash
curl -X DELETE http://localhost:8000/tasks/{id}
```

### Get a single task by ID

```bash
curl -X GET http://localhost:8000/tasks/{id}
```

Replace `{id}` with the actual ID of the task you want to modify or retrieve.

Feel free to customize and modify these cURL examples to suit your needs.