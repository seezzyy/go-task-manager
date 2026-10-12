# Go Task Manager

REST API for managing tasks, built with Go and PostgreSQL.

![Go CI](https://github.com/seezzyy/go-task-manager/actions/workflows/ci.yml/badge.svg)

## Tech Stack

- Go
- PostgreSQL
- Docker & Docker Compose
- golang-migrate
- GitHub Actions

## Features

- Create tasks
- Get tasks by ID
- List all tasks
- Update tasks
- Delete tasks
- Task validation
- PostgreSQL persistence
- Database migrations
- Unit and integration tests

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /tasks | Get all tasks |
| GET | /tasks/{id} | Get task by ID |
| POST | /tasks | Create a new task |
| PUT | /tasks/{id} | Update a task |
| DELETE | /tasks/{id} | Delete a task |

### Create a task

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go",
    "description": "Build a REST API"
  }'
```

### Get all tasks

```bash
curl http://localhost:8080/tasks
```

### Get task by ID

```bash
curl http://localhost:8080/tasks/1
```

### Delete a task

```bash
curl -X DELETE http://localhost:8080/tasks/1
```

## Getting Started

### Requirements

- Docker
- Docker Compose

### Installation

1. Clone the repository:

```bash
git clone https://github.com/seezzyy/go-task-manager.git
cd go-task-manager
```

2. Create a `.env` file:

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
POSTGRES_DB=taskmanager
```

3. Start the application:

```bash
docker compose up -d --build
```

4. Check the API:

```bash
curl http://localhost:8080/tasks
```

The API runs at `http://localhost:8080`.

Database migrations are applied automatically by Docker Compose.

## Testing

Run unit tests:

```bash
go test ./...
```

Run integration tests:

```bash
go test -tags=integration ./internal/storage -v
```

Integration tests require a running PostgreSQL database named `taskmanager_test` with migrations applied.

Run static analysis:

```bash
go vet ./...
```

## CI/CD

GitHub Actions automatically runs:

- Go build
- Unit tests
- Static analysis with go vet

Workflow configuration: `.github/workflows/ci.yml`