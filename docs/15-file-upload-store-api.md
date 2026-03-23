# 15 — FileUpload + Analyzer Store + API Wrapper

**Complexity:** Complex
**Phase:** 4 — Frontend Analyzer
**Blocked by:** 10-analyze-handler, 14-router-layout-nav
**Blocks:** 16-summary-cards, 17-log-table-filter, 18-time-series-chart, 19-pattern-list

## Objective

File upload component, Zustand analyzer store, and fetch API wrapper for `POST /api/analyze`.

## Scope

- `src/api/analyze.ts` — fetch wrapper: build `FormData`, POST to `/api/analyze`, return typed `AnalysisResult`
- `src/stores/analyzerStore.ts` — Zustand store:
  - State: status (idle/uploading/done/error), error, result, filters
  - Actions: `upload(file)`, `cancel()`, `setFilters()`, `reset()`
  - `upload()` uses `AbortController` for cancellation
- `src/components/analyzer/FileUpload.tsx` — drag-and-drop + click upload
  - Drag states (hover highlight)
  - File type/size validation (client-side, before upload)
  - Triggers `analyzerStore.upload()`
  - Shows upload progress/loading state

## Acceptance Criteria

- [x] Drag-and-drop uploads file and triggers analysis
- [x] Click-to-browse uploads file and triggers analysis
- [x] Store transitions: idle → uploading → done (with result) or error
- [x] Cancel aborts in-flight request
- [x] API wrapper sends correct multipart request
- [x] Error responses displayed to user
- [x] Store tests cover all state transitions
