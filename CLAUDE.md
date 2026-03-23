# LogScope v2 — Claude Instructions

Stateless log analyzer + generator. Go backend (streaming, memory-safe) + React frontend. No DB, no disk persistence — all state is request-scoped. Full design in [`docs/architecture-plan.md`](docs/architecture-plan.md).

## Key Constraints

- **Never `ReadAll` on uploads** — always `bufio.Scanner` over `io.Reader` for streaming
- **Memory caps (hardcoded constants):** entries ≤ 10K, pattern templates ≤ 10K
- **Upload limit:** 100MB via `http.MaxBytesReader` in analyze handler
- **No global state** — all processing is per-request; Server struct holds only config
- **Context propagation:** all analyzer/generator funcs accept `context.Context` from `r.Context()`
- **CORS:** dev mode only (`localhost:5173`), controlled via env var/flag — not in production

## Error Response Shape

```json
{ "error": "human-readable message", "code": "error_code" }
```
Codes: `file_too_large` (413), `invalid_format` (400), `invalid_config` (400), `internal_error` (500).

## Frontend Rules

- Components are dumb — no business logic; read Zustand stores, dispatch actions only
- `filteredEntries` derived via `useMemo` selector in components, not in store
- Client-side filtering runs on capped 10K entries (search, level, time range)
- `LogTable` must use `@tanstack/react-virtual` — 10K DOM nodes will lag

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.22+ |
| Frontend | Vite + React + TypeScript |
| Styling | Tailwind CSS + shadcn/ui |
| Charts | Recharts |
| State | Zustand |
| Testing (Backend) | Go standard testing + benchmarks |
| Testing (Frontend) | Vitest + React Testing Library |
