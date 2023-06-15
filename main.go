package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Task represents the structure of our resource.
type Task struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

var (
	db    *sql.DB
	DEBUG bool = true
)

// initDB initializes the database connection.
func initDB() {
	var err error
	// Establish a connection to the database.
	db, err = sql.Open("mysql", "root:Plantago1$@tcp(127.0.0.1:3306)/todoapp")
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection.
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Database connection initialized")
}

// logOperation logs the operation to the console if DEBUG is enabled.
func logOperation(operation string) {
	if DEBUG {
		log.Println(operation)
	}
}

// getTasksHandler fetches all tasks from the database.
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	logOperation("List all tasks")

	// Query the database.
	rows, err := db.Query("SELECT id, text, completed FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate through the results.
	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Text, &task.Completed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	// Return JSON response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// addTaskHandler adds a new task to the database.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	logOperation("Create a new task")

	// Parse the incoming data.
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert data into the database.
	res, err := db.Exec("INSERT INTO tasks (text, completed) VALUES (?, ?)", task.Text, task.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the ID of the inserted item.
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	task.ID = int(id)

	// Return the newly created item.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// editTaskHandler updates an existing task.
func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ID from the URL.
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	logOperation(fmt.Sprintf("Update task with ID %d", id))

	// Parse the incoming data.
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the item in the database.
	_, err = db.Exec("UPDATE tasks SET text = ?, completed = ? WHERE id = ?", task.Text, task.Completed, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task.ID = id

	// Return the updated item.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// deleteTaskHandler deletes a task from the database.
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ID from the URL.
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	logOperation(fmt.Sprintf("Delete task with ID %d", id))

	// Remove the item from the database.
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a 204 status to indicate success but with no response body.
	w.WriteHeader(http.StatusNoContent)
}

// getTaskByIDHandler retrieves a single task by its ID.
func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ID from the URL.
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	logOperation(fmt.Sprintf("Get task with ID %d", id))

	// Query the database for the task.
	row := db.QueryRow("SELECT id, text, completed FROM tasks WHERE id = ?", id)

	// Populate the task struct with the retrieved data.
	task := &Task{}
	err := row.Scan(&task.ID, &task.Text, &task.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	// Define our route handlers.
	r.HandleFunc("/tasks", getTasksHandler).Methods("GET")
	r.HandleFunc("/tasks", addTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", editTaskHandler).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTaskHandler).Methods("DELETE")
	r.HandleFunc("/tasks/{id}", getTaskByIDHandler).Methods("GET")

	// Print server start message to the terminal.
	fmt.Println("Starting Todo app server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
