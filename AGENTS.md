# AGENTS.md

## Build & Run

- **Run directly:** `go run ./cmd/api` (runs on `:9001` or `SERVER_PORT`)
- **Run migrations:** `go run ./cmd/migrate [up|down|version|force]`
- **Live-reload:** `air` (uses `.air.toml`, builds to `tmp/main.exe`)
- **Module name:** `backend` — imports use `backend/...` prefix
- **Go version:** `1.26.0` (per `go.mod`)

## Testing

- **Run all unit tests:** `go test ./...` (fast, offline, uses mocks, zero DB dependency)
- **Run unit tests with coverage:** `go test -v -cover ./...`
- **Run integration tests:** `go test -tags=integration -v ./tests/integration/...` (requires PostgreSQL)

## Env

- Env file lives at `internal/configs/.env` (not root `.env`)
- Loaded via `godotenv.Load("internal/configs/.env")`
- Default DB: `127.0.0.1:5433`, user `admin`, pass `admin123`, db `localDB`
- Auth configs: `JWT_SECRET`, `JWT_EXPIRY_MINUTES`, `REFRESH_TOKEN_SECRET`, `REFRESH_TOKEN_EXPIRY_HOURS`, `REFRESH_TOKEN_GRACE_PERIOD_SECONDS`
- `SERVER_PORT` env var is defined (defaults to `9001` if omitted)

## Architecture

```
cmd/api/main.go             — web server entrypoint
cmd/migrate/main.go         — database migration CLI runner (golang-migrate)
migrations/                 — versioned SQL migrations (000001 to 000005)
internal/configs/           — env loading + GORM Postgres connection
internal/interfaces/        — domain contracts (Repository & Service interfaces)
internal/mocks/             — shared mock repositories and services for unit tests
internal/api/router/        — Gin engine setup
internal/api/routes/        — dependency wiring + route registration
internal/api/handlers/      — HTTP handlers (Gin) + co-located unit tests (*_test.go)
internal/api/middlewares/   — Gin middlewares + co-located unit tests (*_test.go)
internal/api/requests/      — request DTOs
internal/api/responses/     — response DTOs
internal/models/            — domain & GORM model structs
internal/repositories/      — data access layer (GORM queries)
internal/services/          — business logic layer + co-located unit tests (*_test.go)
pkg/apperrors/              — domain error definitions
pkg/helpers/                — pagination, response, context helpers
pkg/logger/                 — centralized Uber Zap JSON logger
pkg/security/               — bcrypt password hashing & JWT token handling + unit tests
tests/integration/          — full end-to-end integration tests hitting PostgreSQL
```

Clean architecture: `handler (API) → service (logic) → repository (data access)` decoupled via `internal/interfaces`.

## API

- Versioned under `/api/v1` — e.g. `/api/v1/auth`, `/api/v1/users`, `/api/v1/products`, `/api/v1/orders`
- **Auth Endpoints:**
  - `POST /api/v1/auth/register` (public)
  - `POST /api/v1/auth/login` (public)
  - `POST /api/v1/auth/refresh` (public — token rotation with reuse detection)
  - `POST /api/v1/auth/logout` (protected — revokes refresh token in DB)
  - `GET /api/v1/auth/me` (protected — profile of current authenticated user)
- **Protected Resource Endpoints (Require Bearer Token):**
  - `/api/v1/users` & `/api/v1/users/:id`
  - `/api/v1/products` & `/api/v1/products/:id`
  - `/api/v1/orders` & `/api/v1/orders/:id`
- Pagination via `?page=1&page_size=10` (default page=1, page_size=10, max 100)
- Response shape: `{"error":bool,"status":int,"message":string,"data":any,"errors":any}`
- All IDs are UUIDv4 (DB default `gen_random_uuid()`)

## Auth

- Dual-Token JWT with stateful refresh tokens stored in PostgreSQL (`refresh_tokens` table)
- Access token validated statelessly in `internal/api/middlewares/middleware.go` via `Authenticate(jwtSecret)`
- Token rotation with grace period (`REFRESH_TOKEN_GRACE_PERIOD_SECONDS`, default 30s) to absorb concurrent frontend requests
- Replay / reuse attack detection: presentation of rotated token past grace period revokes all user sessions
- Explicit logout (`POST /api/v1/auth/logout`) revokes immediately with zero grace period
- Injects `user_id` (string), `user_uuid` (uuid.UUID), `user_email`, `username` into `*gin.Context`

## State

- Unit tests in `internal/api/handlers`, `internal/api/middlewares`, `internal/services`, `pkg/security`
- Integration test in `tests/integration` (`//go:build integration`)
- Linter configured via `.golangci.yml` (`golangci-lint`)
- `tmp/` (Air build output) is gitignored
