# Go CRUD API with Docker & PostgreSQL

A professional, containerized RESTful API for managing a to-do list, built with standard library Go and a PostgreSQL database. This project is designed to showcase best practices in backend development, including a clean project structure, a multi-service Docker setup with Docker Compose, and a production-ready multi-stage Dockerfile.

---

## Features

- **Full CRUD Functionality**: Create, Read, Update, and Delete to-do items.
- **Robust Database**: Uses a networked PostgreSQL database for persistent, scalable storage.
- **RESTful API Design**: Follows standard REST principles for predictable and clean endpoints.
- **Fully Containerized**: A multi-stage `Dockerfile` creates a tiny, secure, static binary.
- **Local Development with Docker Compose**: Includes a `docker-compose.yml` for a one-command local development setup.
- **Production-Ready**: Connects to a cloud database via environment variables for easy deployment.

---

## Prerequisites

Before you begin, ensure you have the following installed:
- [Go](https://go.dev/doc/install) (Version 1.24.6 or later)
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)

---

## 🚀 Running the Application

### Running with Docker Compose (Recommended for Local Development)

This is the easiest way to run the entire application stack (API + Database) on your local machine.

1.  **Start the services:**
    ```bash
    docker-compose up --build
    ```
    This command will build your Go API image, pull the PostgreSQL image, and start both containers.

2.  **Access the API:**
    The API will be available at `http://localhost:8888`.

3.  **To stop the services:**
    Press `Ctrl + C` in the terminal.

### Running Manually for Production (e.g., on Render)

The `Dockerfile` is optimized for production deployment. A cloud provider like Render will use it to build and run your API, which will connect to a managed cloud database via the `DATABASE_URL` environment variable.

---

## 📖 API Endpoints

The base URL for all endpoints is `http://localhost:8888`.

### Get All To-Do Items

- **Method**: `GET`
- **Endpoint**: `/todos`
- **`curl` Example**: `curl http://localhost:8888/todos`

### Create a New To-Do Item

- **Method**: `POST`
- **Endpoint**: `/todos`
- **`curl` Example**: `curl -X POST -H "Content-Type: application/json" -d '{"task":"Deploy to the cloud", "completed":false}' http://localhost:8888/todos`

### Update an Existing To-Do Item

- **Method**: `PUT`
- **Endpoint**: `/todos/{id}`
- **`curl` Example**: `curl -X PUT -H "Content-Type: application/json" -d '{"task":"Deploy to the cloud", "completed":true}' http://localhost:8888/todos/1`

### Delete a To-Do Item

- **Method**: `DELETE`
- **Endpoint**: `/todos/{id}`
- **`curl` Example**: `curl -X DELETE http://localhost:8888/todos/1`