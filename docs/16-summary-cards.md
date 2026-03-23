# 16 — Summary Cards

**Complexity:** Simple
**Phase:** 4 — Frontend Analyzer
**Blocked by:** 15-file-upload-store-api
**Blocks:** 23-integration-build

## Objective

Display analysis summary: total lines, level breakdown, time range, top sources.

## Scope

- `src/components/analyzer/SummaryCards.tsx`
- Read from `analyzerStore.result.summary`
- Cards for:
  - Total lines count
  - Level breakdown (error/warn/info/debug with color coding)
  - Time range (formatted start–end)
  - Top sources list
- Responsive grid layout
- Empty/loading states

## Acceptance Criteria

- [x] Displays all summary fields from analysis result
- [x] Level counts are color-coded
- [x] Handles missing fields gracefully (e.g., no sources)
- [x] Responsive layout on different screen sizes
