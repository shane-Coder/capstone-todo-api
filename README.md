# Go CRUD API with Docker & SQLite

A simple, robust, and containerized RESTful API for managing a to-do list, built with standard library Go and SQLite. This project serves as a practical demonstration of backend development fundamentals, including API design, database integration, and containerization with Docker.

---

## Features

- **Full CRUD Functionality**: Create, Read, Update, and Delete to-do items.
- **Persistent Storage**: Data is saved in a local SQLite database, ensuring it persists across server restarts.
- **RESTful API Design**: Follows standard REST principles for predictable and clean endpoints.
- **Containerized**: Fully containerized with a multi-stage `Dockerfile` for easy and reproducible deployment.
- **Professional Configuration**: Manages configuration (like the port) via environment variables.

---

## Prerequisites

Before you begin, ensure you have the following installed:
- [Go](https://go.dev/doc/install) (Version 1.24.6 or later)
- [Docker](https://docs.docker.com/get-docker/)

---

## 🚀 Running the Application

There are two ways to run this project: locally with Go, or inside a Docker container.

### 1. Running with Docker (Recommended)

This is the easiest and most reliable way to run the application.

1.  **Build the Docker image:**
    ```bash
    docker build -t todo-api .
    ```

2.  **Run the Docker container:**
    ```bash
    docker run -p 8888:8888 todo-api
    ```
The API will be available at `http://localhost:8888`.

### 2. Running Locally with Go

1.  **Create an environment file:**
    Create a file named `.env` in the root of the project and add the port number:
    ```
    PORT=8888
    ```

2.  **Run the server:**
    ```bash
    go run main.go
    ```
The API will be available at `http://localhost:8888`.

---

## 📖 API Endpoints

The base URL for all endpoints is `http://localhost:8888`.

### Get All To-Do Items

- **Method**: `GET`
- **Endpoint**: `/todos`
- **Description**: Retrieves a list of all to-do items.
- **`curl` Example**:
  ```bash
  curl http://localhost:8888/todos

### Create a New To-Do Item

- **Method**: `POST`
- **Endpoint**: `/todos`
- **Description**: Adds a new to-do item to the database.
- **`curl` Example**:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"task":"Learn Docker", "completed":false}' http://localhost:8888/todos
  
### Update an Existing To-Do Item

- **Method**: `PUT`
- **Endpoint**: `/todos/{id}`
- **Description**: Updates the task and completed status of an existing to-do item.
- **`curl` Example**:
  ```bash
  curl -X PUT -H "Content-Type: application/json" -d '{"task":"Master Docker", "completed":true}' http://localhost:8888/todos/1
 
### Delete a To-Do Item

- **Method**: `DELETE`
- **Endpoint**: `/todos/{id}`
- **Description**: Removes a to-do item from the database.
- **`curl` Example**:
  ```bash  
    curl -X DELETE http://localhost:8888/todos/1