# 17 — Log Table + Filter Bar

**Complexity:** Complex
**Phase:** 4 — Frontend Analyzer
**Blocked by:** 15-file-upload-store-api
**Blocks:** 23-integration-build

## Objective

Virtualized log table with client-side search, level filtering, time range filtering, and sorting.

## Scope

- `src/components/analyzer/LogTable.tsx`
  - Virtualized via `@tanstack/react-virtual` (10K entries would lag without it)
  - Columns: line number, timestamp, level, source, message
  - Level cells color-coded
  - Sortable columns (at minimum: timestamp, level)
- `src/components/analyzer/FilterBar.tsx`
  - Search box (filters message content)
  - Level checkboxes (error/warn/info/debug)
  - Time range picker (start/end)
- Client-side filtering: `useMemo` selector over `analyzerStore.result.entries` + `analyzerStore.filters`
- Filters dispatch to `analyzerStore.setFilters()`

## Performance

- Filtering on 10K entries must remain responsive
- `useMemo` with proper dependency tracking
- Virtualization renders only visible rows

## Acceptance Criteria

- [x] Table renders 10K entries without lag
- [x] Search filters entries by message content
- [x] Level checkboxes filter by selected levels
- [x] Time range picker filters by timestamp
- [x] Filters combine (AND logic)
- [x] Sortable by timestamp and level
- [x] Empty state when no entries match filters
