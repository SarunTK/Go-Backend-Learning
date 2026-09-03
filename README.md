# Go Backend Learning

A practical backend starter built with Go, PostgreSQL, JWT, gRPC, and Kubernetes-friendly project structure.

## Stack

- Go
- PostgreSQL
- HTTP REST API
- JWT authentication
- gRPC contract example
- Docker Compose for local development
- Kubernetes-ready folder structure

## Project structure

- `cmd/server` — entry point of the app
- `internal/config` — environment configuration
- `internal/db` — PostgreSQL connection and migration runner
- `internal/model` — data models
- `internal/repository` — database access layer
- `internal/service` — business logic and validation
- `internal/httpserver` — REST handlers
- `internal/auth` — JWT helpers
- `internal/grpcserver` — gRPC server scaffold
- `db/migrations` — SQL migration files
- `proto` — gRPC .proto definition
- `k8s` — Kubernetes example manifest

## Prerequisites

- Go 1.22+
- Docker Desktop or Docker Engine
- PostgreSQL (or use Docker Compose)

## Local run with Docker Compose

```bash
docker compose up --build
```

This starts:
- PostgreSQL on `localhost:5432`
- App on `localhost:8080`
- gRPC server on `localhost:9090`

## Local run without Docker

1. Create `.env` from `.env.example`
2. Start PostgreSQL manually
3. Run:

```bash
go mod tidy
go run ./cmd/server
```

## Health check

```bash
curl http://localhost:8080/healthz
```

## Login example

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@example.com"}'
```

## Example authenticated request

```bash
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

## Database migration

The project automatically runs SQL files from `db/migrations` when the app starts.

## CI workflow

GitHub Actions runs on push and pull request:

- `go vet ./...`
- `go test ./...`

The workflow file is in `.github/workflows/ci.yml`.

## Notes

- `.env` is ignored by git
- `.env.example` is kept as a template
- Docker image build is included, but Kubernetes deployment is only scaffolded and not a full production cluster config
