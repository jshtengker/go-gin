# AGENTS.md

## Build & Run

- **Run directly:** `go run ./cmd/api` (runs on `:9001`)
- **Live-reload:** `air` (uses `.air.toml`, builds to `tmp/main.exe`)
- **Module name:** `backend` — imports use `backend/...` prefix
- **Go version:** `1.25.5` (per `go.mod`)

## Env

- Env file lives at `internal/configs/.env` (not root `.env`)
- Loaded via `godotenv.Load("internal/configs/.env")`
- Default DB: `127.0.0.1:5433`, user `admin`, pass `admin123`, db `localDB`
- `SERVER_PORT` env var is defined (defaults to `9001` if omitted)

## Architecture

```
cmd/api/main.go             — entrypoint
internal/configs/           — env loading + GORM Postgres connection
internal/interfaces/        — domain contracts (Repository & Service interfaces)
internal/api/router/        — Gin engine setup
internal/api/routes/        — dependency wiring + route registration
internal/api/handlers/      — HTTP handlers (Gin)
internal/api/middlewares/    — Gin middlewares
internal/api/requests/      — request DTOs
internal/api/responses/     — response DTOs
internal/models/            — domain & GORM model structs
internal/repositories/      — data access layer (GORM queries)
internal/services/          — business logic layer
pkg/apperrors/              — domain error definitions
pkg/helpers/                — pagination, response helpers
pkg/logger/                 — centralized Uber Zap JSON logger
```

Clean architecture: `handler (API) → service (logic) → repository (data access)` decoupled via `internal/interfaces`.

## API

- Versioned under `/api/v1` — e.g. `/api/v1/users`, `/api/v1/products`, `/api/v1/orders`
- **Only GET endpoints exist** — `GetAll` (paginated) and `GetByID` per resource
- Pagination via `?page=1&page_size=10` (default page=1, page_size=10, max 100)
- Response shape: `{"error":bool,"status":int,"message":string,"data":any,"errors":any}`
- All IDs are UUIDv4 (DB default `gen_random_uuid()`)

## Auth

- `internal/api/middlewares/middleware.go` has a simple `Authenticate()` middleware (checks `Token: auth` header)
- **Not wired into any route** — routes register handlers directly without middleware

## State

- No tests exist (no `_test.go` files)
- Linter configured via `.golangci.yml` (`golangci-lint`)
- `tmp/` (Air build output) is gitignored

