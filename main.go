package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

// The Server struct will hold dependencies like the database connection.
type Server struct {
	db *sql.DB
}

func NewServer() *Server {
	db, err := sql.Open("sqlite3", "./todos.db")
	if err != nil {
		panic(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS todos (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"task" TEXT,
		"completed" BOOLEAN
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		panic(err)
	}

	return &Server{db: db}
}

func (s *Server) getTodosHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query("SELECT id, task, completed FROM todos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Task, &todo.Completed); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (s *Server) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	result, err := s.db.Exec("INSERT INTO todos (task, completed) VALUES (?, ?)", newTodo.Task, newTodo.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	newTodo.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func (s *Server) updateTodoHandler(w http.ResponseWriter, r *http.Request, id int) {
	var updatedTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	result, err := s.db.Exec("UPDATE todos SET task = ?, completed = ? WHERE id = ?", updatedTodo.Task, updatedTodo.Completed, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	updatedTodo.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

func (s *Server) deleteTodoHandler(w http.ResponseWriter, _ *http.Request, id int) {
	_, err := s.db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) todosRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/todos")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		switch r.Method {
		case "GET":
			s.getTodosHandler(w, r)
		case "POST":
			s.createTodoHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "PUT":
		s.updateTodoHandler(w, r, id)
	case "DELETE":
		s.deleteTodoHandler(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	server := NewServer()

	http.HandleFunc("/todos/", server.todosRouter)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	fmt.Printf("✅ API Server starting on port: %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
