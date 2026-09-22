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
- `SERVER_PORT` env var is defined but **not used** — the server hardcodes `:9001`

## Architecture

```
cmd/api/main.go             — entrypoint
internal/configs/           — env loading + GORM Postgres connection
internal/api/router/        — Gin engine setup
internal/api/routes/        — dependency wiring + route registration
internal/api/handlers/      — HTTP handlers (Gin)
internal/api/middlewares/    — Gin middlewares
internal/api/requests/      — request DTOs
internal/api/responses/     — response DTOs
internal/models/            — GORM model structs
internal/repositories/      — data access layer (GORM queries)
internal/services/          — business logic layer
pkg/helpers/                — pagination, response helpers
```

Clean architecture: `handler → service → repository`.

## API

- Versioned under `/v1` — e.g. `/v1/users`, `/v1/products`, `/v1/orders`
- **Only GET endpoints exist** — `GetAll` (paginated) and `GetById` per resource
- Pagination via `?page=1&page_size=10` (default page=1, page_size=10, max 100)
- Response shape: `{"error":bool,"status":int,"message":string,"data":any,"errors":any}`
- All IDs are UUIDv4 (DB default `gen_random_uuid()`)

## Auth

- `internal/api/middlewares/middleware.go` has a simple `Authenticate()` middleware (checks `Token: auth` header)
- **Not wired into any route** — routes register handlers directly without middleware

## State

- No tests exist (no `_test.go` files)
- No linter, formatter, or CI config present
- `tmp/` (Air build output) is gitignored
