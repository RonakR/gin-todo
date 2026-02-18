# Gin Todo API

Simple in-memory todo API using Gin.

## Prereqs

- Go 1.22+

## Run

```bash
go mod tidy
go run .
```

Server starts at `http://localhost:8080`.

## Endpoints

- `GET /health`
- `GET /todos`
- `POST /todos`
- `GET /todos/:id`
- `PATCH /todos/:id`
- `DELETE /todos/:id`

## Example requests

```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"buy milk"}'

curl http://localhost:8080/todos

curl -X PATCH http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'

curl -X DELETE http://localhost:8080/todos/1
```
