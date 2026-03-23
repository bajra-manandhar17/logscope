# 23 — Integration & Build Verification

**Complexity:** Moderate
**Phase:** 6 — Integration
**Blocked by:** 16-summary-cards, 17-log-table-filter, 18-time-series-chart, 19-pattern-list, 21-live-preview, 22-action-bar
**Blocks:** 24-test-suite, 25-e2e-verification

## Objective

Verify single binary build works end-to-end: Go embeds frontend, serves SPA + API.

## Scope

- `make build` produces single binary at `bin/logscope`
- Binary serves:
  - Embedded React SPA on `GET /*`
  - All API endpoints (`/api/analyze`, `/api/generate`, `/api/health`)
- Verify with curl:
  - `curl localhost:8080/api/health` → 200
  - `curl -F "file=@sample.log" localhost:8080/api/analyze` → valid JSON
  - `curl -X POST -d '{"format":"json","totalLines":100,...}' localhost:8080/api/generate` → SSE stream
- SPA loads in browser, routing works

## Acceptance Criteria

- [x] `make build` succeeds
- [x] Single binary runs and serves SPA + API
- [x] Health endpoint responds
- [x] Analyze endpoint processes file upload
- [x] Generate endpoint streams SSE
- [x] SPA loads and navigates between pages
