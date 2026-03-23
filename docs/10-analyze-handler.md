# 10 — Analyze HTTP Handler

**Complexity:** Moderate
**Phase:** 2 — Backend Analyzer
**Blocked by:** 02-server-entrypoint, 09-analyzer-orchestrator
**Blocks:** 15-file-upload-store-api, 22-action-bar

## Objective

HTTP handler for `POST /api/analyze` — parse multipart upload, call analyzer, return JSON.

## Scope

- `internal/handler/analyze.go`
- Parse `multipart/form-data`, field name `"file"`
- Enforce 100MB limit via `http.MaxBytesReader`
- Optional query param `?format=auto|json|plaintext`
- Pass `multipart.File` (io.Reader) + format hint to `analyzer.Analyze()`
- Return `200 JSON` with `AnalysisResult`
- Error responses use contract: `{ "error": "msg", "code": "error_code" }`
  - 413 `file_too_large`
  - 400 `invalid_format`
  - 500 `internal_error`
- Register route on ServeMux: `POST /api/analyze`
- `internal/handler/analyze_test.go` — httptest round-trips

## Acceptance Criteria

- [x] Accepts multipart file upload and returns analysis JSON
- [x] Enforces 100MB limit with 413 response
- [x] Returns proper error codes for bad requests
- [x] Content-Type is `application/json`
- [x] Handler tests pass with httptest
