package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	// NOTE: We no longer need the sqlite3 driver
)

type Todo struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

type Server struct {
	db *pgxpool.Pool // Use the pgx connection pool
}

func NewServer() *Server {
	// Get the database connection URL from the environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("DATABASE_URL environment variable is not set")
	}

	// Create a new connection pool
	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}

	// Create the table, using PostgreSQL syntax
	createTableSQL := `CREATE TABLE IF NOT EXISTS todos (
		"id" SERIAL PRIMARY KEY,
		"task" TEXT,
		"completed" BOOLEAN
	);`
	if _, err := db.Exec(context.Background(), createTableSQL); err != nil {
		panic(err)
	}

	return &Server{db: db}
}

// --- Handlers (Updated for PostgreSQL) ---

func (s *Server) getTodosHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(context.Background(), "SELECT id, task, completed FROM todos ORDER BY id")
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

	// Use QueryRow to get the ID back after inserting
	err := s.db.QueryRow(context.Background(), "INSERT INTO todos (task, completed) VALUES ($1, $2) RETURNING id", newTodo.Task, newTodo.Completed).Scan(&newTodo.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	// Use Exec and check the command tag for affected rows
	tag, err := s.db.Exec(context.Background(), "UPDATE todos SET task = $1, completed = $2 WHERE id = $3", updatedTodo.Task, updatedTodo.Completed, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tag.RowsAffected() == 0 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	updatedTodo.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

func (s *Server) deleteTodoHandler(w http.ResponseWriter, _ *http.Request, id int) {
	tag, err := s.db.Exec(context.Background(), "DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tag.RowsAffected() == 0 {
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
