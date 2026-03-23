# 06 — Summary Aggregator

**Complexity:** Moderate
**Phase:** 2 — Backend Analyzer
**Blocked by:** 03-analyzer-types
**Blocks:** 09-analyzer-orchestrator

## Objective

Aggregate O(1) memory summary stats from a stream of log entries: level counts, time range, top sources.

## Scope

- `internal/analyzer/summary.go`
- Incremental aggregation — call per-entry, extract summary at end
- Track: totalLines, per-level counts, min/max timestamp, source frequency map
- Top sources: maintain map, extract top N at end (N=10)
- Skip `topSources` if no source found in any entry
- `internal/analyzer/summary_test.go`

## Acceptance Criteria

- [x] Correct level counts across mixed-level input
- [x] Correct time range (earliest to latest)
- [x] Top sources ranked by frequency
- [x] O(1) memory (bounded source map — cap if needed)
- [x] Handles entries with missing timestamps/levels/sources gracefully
- [x] Tests pass
