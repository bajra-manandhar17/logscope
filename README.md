# LogScope v2

High-performance, stateless log analyzer + generator. Go backend with streaming memory-safe parsing; React frontend for visualization. Ships as a single Go binary embedding the React SPA.

## Features

- Upload and analyze structured JSON + unstructured plaintext logs (up to 100MB)
- Auto-detect log format; extract levels, time range, top sources, patterns
- Time-series chart of log volume + error rate
- Pattern detection via token-masking (`{NUM}`, `{UUID}`, `{IP}`, `{HEX}`)
- Log generator with configurable format, volume, level distribution, and time range
- SSE streaming preview during generation; send generated logs directly to analyzer

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

## Quick Start

**Development**

```bash
# Terminal 1 — Go backend on :8080
make dev-backend

# Terminal 2 — Vite on :5173 (proxies /api → :8080)
make dev-frontend
```

**Production build**

```bash
make build          # outputs bin/logscope
./bin/logscope      # serves SPA + API on :8080
```

## API

```
POST /api/analyze    multipart/form-data (field: "file"), ?format=auto|json|plaintext
POST /api/generate   application/json → text/event-stream (SSE batches)
GET  /api/health     → { status: "ok" }
GET  /*              embedded React SPA
```

See [`docs/architecture-plan.md`](docs/architecture-plan.md) for full API contracts, request/response shapes, and design decisions.

## Project Structure

```
cmd/logscope/          entry point + go:embed
internal/
  analyzer/            parser, detector, summary, pattern, orchestrator
  generator/           config validation + streaming line generator
  handler/             HTTP handlers (analyze, generate)
frontend/src/
  pages/               AnalyzerPage, GeneratorPage
  components/          analyzer/, generator/, shared/
  stores/              Zustand stores (analyzerStore, generatorStore)
  api/                 fetch wrappers for /api/analyze and /api/generate
docs/
  architecture-plan.md full design doc
```

## Testing

```bash
go test ./internal/...          # backend unit + integration tests
cd frontend && npm test         # frontend Vitest tests
```
