# Go Backend Learning

This project is a practical backend starter based on the stack:

- Go
- PostgreSQL
- REST API
- JWT authentication
- gRPC contract example
- Kubernetes-ready structure

This repository intentionally skips Docker image building and Kubernetes deployment execution, but includes the code structure and conventions used in a production backend.

## Project layout

- `cmd/server` - application entry point
- `internal/config` - environment configuration
- `internal/db` - PostgreSQL and migration runner
- `internal/model` - domain models
- `internal/repository` - database access layer
- `internal/service` - business logic and validation
- `internal/httpserver` - REST handlers
- `internal/auth` - JWT helpers
- `db/migrations` - database schema files
- `proto` - gRPC contract definition
- `k8s` - Kubernetes manifest example

## Run locally

1. Create a local `.env` file based on `.env.example`
2. Start PostgreSQL
3. Run:

```bash
go mod tidy
go run ./cmd/server
```

The HTTP API will be available on `:8080` and the gRPC service scaffold is on `:9090`.

## Example API

- `GET /healthz`
- `POST /api/v1/login`
- `POST /api/v1/users` with `Authorization: Bearer <token>`
- `GET /api/v1/users`
- `GET /api/v1/users/{id}`

## Example login

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@example.com"}'
```

## Notes

The gRPC contract is defined under `proto/user.proto` and can be generated with `protoc` later.
