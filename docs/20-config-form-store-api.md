# 20 — Config Form + Generator Store + SSE API

**Complexity:** Complex
**Phase:** 5 — Frontend Generator
**Blocked by:** 13-generate-handler, 14-router-layout-nav
**Blocks:** 21-live-preview, 22-action-bar

## Objective

Generator config form, Zustand store, and SSE API wrapper using fetch + ReadableStream.

## Scope

- `src/api/generate.ts`
  - `fetch()` POST to `/api/generate` with JSON body
  - Parse SSE response via `ReadableStream` + custom parser (not `EventSource` — POST not supported by EventSource API)
  - Yield parsed events: `batch` (lines[]) and `done` (totalLines)
  - Handle stream errors and abort
- `src/stores/generatorStore.ts` — Zustand store:
  - State: config, status (idle/generating/done/error), lines[], error
  - Actions: `generate()`, `abort()`, `sendToAnalyzer()`, `setConfig()`, `reset()`
  - `generate()` appends batches to `lines[]` as they arrive
  - `abort()` cancels fetch via `AbortController`
  - `sendToAnalyzer()` creates Blob from lines, calls analyzerStore.upload(), navigates to `/analyze`
- `src/components/generator/ConfigForm.tsx`
  - Format selector (json/plaintext)
  - Total lines input (1–1,000,000)
  - Level distribution sliders/inputs (must sum to 1.0)
  - Time range pickers (start/end)
  - Client-side validation before submit

## Acceptance Criteria

- [x] SSE parser correctly handles batch and done events
- [x] Store accumulates lines from batches
- [x] Abort cancels in-flight generation
- [x] Config form validates all fields
- [x] Level sliders enforce sum-to-1.0 constraint
- [x] Store tests cover state transitions
- [x] API tests mock ReadableStream
