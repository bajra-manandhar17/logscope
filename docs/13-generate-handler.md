# 13 — Generate SSE Handler

**Complexity:** Moderate
**Phase:** 3 — Backend Generator
**Blocked by:** 02-server-entrypoint, 12-streaming-generator
**Blocks:** 20-config-form-store-api

## Objective

HTTP handler for `POST /api/generate` — parse config JSON, stream SSE response.

## Scope

- `internal/handler/generate.go`
- Parse `application/json` body → `GenerateConfig`
- Validate config, return 400 `invalid_config` on failure
- Set response headers: `Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive`
- Stream batches as SSE events:
  - `event: batch` → `data: { "lines": [...] }`
  - `event: done` → `data: { "totalLines": N }`
- Use `http.Flusher` to flush after each event
- Pass `r.Context()` to generator — client disconnect stops generation
- Register route: `POST /api/generate`
- `internal/handler/generate_test.go` — httptest SSE verification

## Acceptance Criteria

- [x] Returns SSE stream with correct headers
- [x] Each batch event contains lines array
- [x] Done event contains correct totalLines
- [x] Flushes after each event (no buffering)
- [x] Returns 400 for invalid config
- [x] Stops generation on client disconnect
- [x] Handler tests pass
