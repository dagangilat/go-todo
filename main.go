package main

import (
	"database/sql"
	"encoding/json"
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

var db *sql.DB

// Initialize database connection.
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
}

// Fetch all tasks from database.
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
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

// Add a new task to the database.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
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

// Edit an existing task.
func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ID from the URL.
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

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

// Delete a task from the database.
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the ID from the URL.
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Remove the item from the database.
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a 204 status to indicate success but with no response body.
	w.WriteHeader(http.StatusNoContent)
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

	log.Fatal(http.ListenAndServe(":8000", r))
}
