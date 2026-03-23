# LogScope v2 — Architecture Design & Implementation Plan

## Context

LogScope is a greenfield, high-performance, stateless log analyzer + log generator. Local developer tool — single user, no auth, no DB. Go backend streams and parses logs memory-safely; React frontend visualizes results. No disk persistence — all state is request-scoped. Ships as a single Go binary embedding the React SPA.

**Module:** `github.com/bajra-manandhar17/logscope-v2`

**Problem:** Need a tool to upload/analyze logs (structured JSON + unstructured plaintext) with search, filtering, time-series charts, and pattern detection — plus a log generator for testing/demo purposes.

---

## Architecture: REST + SSE

### Why this approach
- SSE handles the only server→client streaming need (log generation preview)
- Analysis returns single JSON response — sufficient for ≤100MB files
- Truly stateless: no session IDs, no server-side hub
- SSE is pure Go stdlib (`http.Flusher` + `fmt.Fprintf`), zero dependencies
- Trivially testable with `httptest`

---

## API Design

```
POST /api/analyze
  Request:  multipart/form-data (field: "file"), optional ?format=auto|json|plaintext
  Limit:    100MB max (enforced via http.MaxBytesReader)
  Response: 200 JSON {
    formatDetected: "json" | "plaintext",
    summary: { totalLines, levels: {error,warn,info,debug}, timeRange, topSources },
    entries: LogEntry[]       (capped at 10,000; totalLines tells full count),
    patterns: Pattern[]       (top N grouped patterns with counts),
    timeSeries: TimeBucket[]  (log volume + error rate bucketed by time interval)
  }

POST /api/generate
  Request:  application/json (GenerateConfig)
  Response: 200 text/event-stream
    event: batch  → data: { lines: string[] }   (~100 lines/batch)
    event: done   → data: { totalLines: int }
  Note: Browser EventSource API only supports GET. Frontend uses fetch() +
        ReadableStream + custom SSE parser to consume POST SSE responses.

GET /api/health → 200 { status: "ok" }
GET /*          → embedded React SPA
```

### Error Response Contract
All error responses use a consistent shape:
```json
{ "error": "human-readable message", "code": "error_code" }
```
Error codes: `file_too_large` (413), `invalid_format` (400), `invalid_config` (400), `internal_error` (500).

### GenerateConfig Definition
```
GenerateConfig {
  format:       "json" | "plaintext"
  totalLines:   int (1..1,000,000)
  levels:       { error: float, warn: float, info: float, debug: float }  // must sum to 1.0
  timeRange:    { start: RFC3339, end: RFC3339 }
}
```

### Time-Series Bucketing
Backend auto-selects bucket interval based on log time range:
- < 1 hour → 1-minute buckets
- < 24 hours → 15-minute buckets
- < 7 days → 1-hour buckets
- ≥ 7 days → 1-day buckets

Response includes `bucketInterval` field so frontend knows the granularity.

### "Send generated logs to analyzer" flow
1. Frontend joins accumulated `lines[]` into a string
2. Creates `new Blob([text], { type: 'text/plain' })`
3. POSTs as `multipart/form-data` to `/api/analyze`
4. Server processes identically to a file upload — zero disk involvement

---

## Backend Structure

```
cmd/logscope/main.go                 — go:embed frontend, wire stdlib ServeMux, start server

internal/
  server/
    server.go                        — http.Server setup, CORS middleware, route registration
    server_test.go
  analyzer/
    types.go                         — LogEntry, Summary, Pattern, TimeBucket structs
    parser.go                        — streaming line-by-line parser (bufio.Scanner over io.Reader)
    parser_test.go
    detector.go                      — auto-detect JSON vs plaintext (peek first lines)
    detector_test.go
    summary.go                       — level counts, time range, top sources (O(1) memory)
    summary_test.go
    pattern.go                       — token-masking: replace numbers/UUIDs/IPs → group by template
    pattern_test.go
    analyzer.go                      — orchestrates: parse → summarize → detect patterns → time-series
    analyzer_test.go
  generator/
    config.go                        — GenerateConfig struct + validation
    generator.go                     — streaming log line generator (yields batches)
    generator_test.go
  handler/
    analyze.go                       — HTTP handler: parse multipart, call analyzer, return JSON
    analyze_test.go
    generate.go                      — HTTP handler: parse config JSON, SSE stream response
    generate_test.go
```

### Key design decisions
- **Router:** Go 1.22 stdlib `http.ServeMux` (method routing, zero dependencies)
- **Upload limit:** 100MB enforced via `http.MaxBytesReader` in analyze handler
- **Memory safety:** `parser.go` wraps `multipart.File` (io.Reader) in `bufio.Scanner` — never `ReadAll`. Summary = O(1) counters. Pattern map capped at 10K unique templates (stop inserting new, only increment existing). Entries capped at 10K. All caps are hardcoded constants.
- **Format detection fallback:** If auto-detection fails (e.g., mixed format), best-effort parse as plaintext. Response includes `formatDetected` field.
- **Pattern detection scope:** Runs on the full file stream (not just capped entries) for better accuracy. 10K template cap bounds memory.
- **Source extraction:** JSON logs: look for first of `source`, `service`, `module`, `logger`, `component` fields. Plaintext: extract bracketed token after level (e.g., `[myservice]`). Skip `topSources` if no source found.
- **Context propagation:** All analyzer/generator functions accept `context.Context` from `r.Context()`. Server sets read/write timeouts. Client disconnect cancels processing.
- **Pattern detection:** Token-masking — replace variable parts ({NUM}, {UUID}, {IP}, {HEX}) with placeholders, group identical templates, rank by count
- **CORS:** Dev mode only — allow `localhost:5173`. Production serves SPA from same origin, no CORS needed. Controlled via env var or flag.
- **No global state:** All processing is per-request. Server struct holds only config (port, limits).

---

## Frontend Structure

```
frontend/
  src/
    main.tsx
    App.tsx                          — React Router: / redirects to /analyze, /generate
    pages/
      AnalyzerPage.tsx               — orchestrates analyzer components
      GeneratorPage.tsx              — orchestrates generator components
    components/
      analyzer/
        FileUpload.tsx               — drag-and-drop + click file upload
        SummaryCards.tsx              — total lines, level breakdown, time range, top sources
        LogTable.tsx                 — virtualized table (tanstack-virtual) with search, level/time filters, sort
        TimeSeriesChart.tsx          — Recharts line/bar chart (volume + error rate over time)
        PatternList.tsx              — grouped patterns with counts
        FilterBar.tsx                — search box, level checkboxes, time range picker
      generator/
        ConfigForm.tsx               — format, volume, level distribution, time range
        LivePreview.tsx              — scrollable log preview during SSE generation
        ActionBar.tsx                — download, copy, send-to-analyzer (auto-navigates to /analyze)
      shared/
        Layout.tsx                   — page shell with nav
        Navigation.tsx               — top/sidebar nav between Analyzer and Generator
    stores/
      analyzerStore.ts               — upload state, result, client-side filters + derived filtered entries
      generatorStore.ts              — config, generation state, accumulated lines, send-to-analyzer action
    api/
      analyze.ts                     — fetch wrapper for POST /api/analyze (multipart)
      generate.ts                    — fetch + ReadableStream SSE parser for POST /api/generate
    types/index.ts                   — shared TS types mirroring Go structs
    lib/utils.ts                     — formatting helpers
  index.html
  vite.config.ts                     — proxy /api → Go backend in dev mode
  tsconfig.json
  tailwind.config.ts
  package.json
```

### Zustand Store Design

```typescript
// analyzerStore
interface AnalyzerState {
  status: 'idle' | 'uploading' | 'done' | 'error'
  error: string | null
  result: AnalysisResult | null
  filters: { search: string; levels: string[]; timeRange: [Date, Date] | null }
  // filteredEntries derived via useMemo selector in components, not in store
  upload: (file: File | Blob) => Promise<void>
  cancel: () => void          // AbortController for in-progress uploads
  setFilters: (f: Partial<Filters>) => void
  reset: () => void
}

// generatorStore
interface GeneratorState {
  config: GenerateConfig
  status: 'idle' | 'generating' | 'done' | 'error'
  lines: string[]
  error: string | null
  generate: () => void        // fetch + ReadableStream, appends batches
  abort: () => void           // AbortController cancels fetch
  sendToAnalyzer: () => void  // Blob → POST /api/analyze, then navigate to /analyze
  setConfig: (c: Partial<GenerateConfig>) => void
  reset: () => void
}
```

### Key frontend decisions
- **Components are dumb** — read from Zustand stores, dispatch actions, no business logic
- **Client-side filtering** on capped 10K entries (search, level filter, time range)
- **Virtualized LogTable** via `@tanstack/react-virtual` — 10K DOM nodes would lag on filter/search
- **Send-to-analyzer** auto-navigates to `/analyze` and displays results
- **shadcn/ui** components copied into repo (no runtime dependency)
- **Recharts** for time-series visualization

---

## Build & Deployment

```makefile
# Makefile
build:
	cd frontend && npm run build          # → frontend/dist/
	go build -o bin/logscope ./cmd/logscope

dev-backend:
	go run ./cmd/logscope                 # serves on :8080

dev-frontend:
	cd frontend && npm run dev            # Vite on :5173, proxy /api → :8080
```

- `cmd/logscope/main.go` uses `//go:embed frontend/dist/*`
- Single binary serves both API + SPA
- Dev mode: Vite proxy to Go backend

---

## Testing Strategy

### Backend (Go standard testing + benchmarks)
- **Unit tests:** parser, detector, summary, pattern, generator — each package independently
- **Handler tests:** `httptest.NewServer` for full HTTP round-trips (upload, SSE stream)
- **Benchmarks:** parser on large inputs (10K, 100K, 1M lines) to verify streaming performance

### Frontend (Vitest + React Testing Library)
- **Store tests:** all Zustand stores — highest value (pure logic, test actions + state transitions)
- **Component tests:** FileUpload, LogTable, ConfigForm only (complex interactions). Skip trivial display components.
- **API tests:** mock fetch/ReadableStream, verify request formation and response handling

---

## Verification Plan

1. **Backend unit tests:** `go test ./internal/...`
2. **Frontend tests:** `cd frontend && npm test`
3. **Integration test:** Start server, upload a sample log file via curl, verify JSON response structure
4. **SSE test:** curl `/api/generate` with config, verify event stream format
5. **E2E manual:** Build single binary, open in browser, upload log → verify dashboard; generate logs → verify preview → send to analyzer
6. **Memory test:** Upload a ~50MB file, monitor Go memory usage stays bounded

---

## Implementation Order

### Phase 1: Project Scaffolding
1. `git init`, `go mod init github.com/bajra-manandhar17/logscope-v2`, frontend scaffolding (Vite + React + TS + Tailwind + shadcn + React Router + @tanstack/react-virtual)
2. Makefile with build targets
3. `cmd/logscope/main.go` — embed + serve static + health endpoint + CORS middleware (dev mode)

### Phase 2: Backend Core — Analyzer
4. `internal/analyzer/types.go` — data types
5. `internal/analyzer/detector.go` — format auto-detection
6. `internal/analyzer/parser.go` — streaming line-by-line parser
7. `internal/analyzer/summary.go` — stats aggregation
8. `internal/analyzer/pattern.go` — token-masking pattern detection
9. `internal/analyzer/analyzer.go` — orchestrator
10. `internal/handler/analyze.go` — HTTP handler + route

### Phase 3: Backend Core — Generator
11. `internal/generator/config.go` — config struct + validation
12. `internal/generator/generator.go` — streaming line generator
13. `internal/handler/generate.go` — SSE handler + route

### Phase 4: Frontend — Analyzer
14. Router, Layout, Navigation
15. FileUpload + analyzerStore + API wrapper
16. SummaryCards
17. LogTable + FilterBar
18. TimeSeriesChart
19. PatternList

### Phase 5: Frontend — Generator
20. ConfigForm + generatorStore + SSE API wrapper
21. LivePreview
22. ActionBar (download, copy, send-to-analyzer)

### Phase 6: Integration & Verification
23. Single binary build verification (`make build` + run)
24. Full test suite pass (`go test ./internal/...` + `npm test`)
25. Manual E2E: upload → dashboard → generate → preview → send-to-analyzer flow

**Note:** Error handling, loading states, and edge cases are built into each phase — not deferred. Each handler/component includes error paths in its tests.
