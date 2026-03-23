# 19 — Pattern List

**Complexity:** Simple
**Phase:** 4 — Frontend Analyzer
**Blocked by:** 15-file-upload-store-api
**Blocks:** 23-integration-build

## Objective

Display detected log patterns grouped by template with occurrence counts.

## Scope

- `src/components/analyzer/PatternList.tsx`
- Read from `analyzerStore.result.patterns`
- Display each pattern:
  - Masked template (with placeholders highlighted)
  - Occurrence count
  - One sample original line (expandable/collapsible)
- Sorted by count descending (already sorted from backend)
- Scrollable list if many patterns
- Empty state when no patterns detected

## Acceptance Criteria

- [x] Displays pattern templates with counts
- [x] Placeholders visually distinct (color/badge)
- [x] Sample line expandable
- [x] Sorted by count descending
- [x] Empty state for no patterns
