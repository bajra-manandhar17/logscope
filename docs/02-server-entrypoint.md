# 02 — Server Entrypoint & Health Endpoint

**Complexity:** Simple
**Phase:** 1 — Scaffolding
**Blocked by:** 01-project-scaffolding
**Blocks:** 10-analyze-handler, 13-generate-handler, 14-router-layout-nav

## Objective

Create `cmd/logscope/main.go` with embedded frontend, stdlib ServeMux, health endpoint, and CORS middleware (dev mode).

## Scope

- `cmd/logscope/main.go` — `//go:embed frontend/dist/*`, serve static files on `GET /*`
- `internal/server/server.go` — `http.Server` setup, CORS middleware, route registration
- `GET /api/health` → `200 { status: "ok" }`
- CORS: allow `localhost:5173` in dev mode (env var or flag controlled)
- Read/write timeouts on server
- `internal/server/server_test.go` — health endpoint test, CORS header test

## Acceptance Criteria

- [x] `go run ./cmd/logscope` starts server on `:8080`
- [x] `curl localhost:8080/api/health` returns `{"status":"ok"}`
- [x] CORS headers present when dev mode enabled
- [x] Embedded SPA served on `GET /`
- [x] Server tests pass
