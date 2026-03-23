# 08 — Time-Series Bucketing

**Complexity:** Moderate
**Phase:** 2 — Backend Analyzer
**Blocked by:** 03-analyzer-types
**Blocks:** 09-analyzer-orchestrator

## Objective

Bucket log entries by time interval for volume + error rate visualization.

## Scope

- `internal/analyzer/timeseries.go` (or extend `summary.go`)
- Auto-select bucket interval based on log time range:
  - < 1 hour → 1-minute buckets
  - < 24 hours → 15-minute buckets
  - < 7 days → 1-hour buckets
  - ≥ 7 days → 1-day buckets
- Each `TimeBucket`: timestamp (bucket start), count, errorCount
- Return `[]TimeBucket` + `bucketInterval` string
- Two-pass or deferred: needs time range (from summary) before bucketing — may need to bucket after streaming, using stored timestamps or a second pass
- `internal/analyzer/timeseries_test.go`

## Design Decision

Since we stream entries and cap at 10K, time-series can be computed from the capped entries list. Alternatively, track bucket map incrementally during streaming (requires estimating range from first entry or using dynamic rebucketing). Simpler: compute from capped entries post-stream.

## Acceptance Criteria

- [x] Correct bucket interval selection for each time range
- [x] Accurate count and errorCount per bucket
- [x] Buckets sorted chronologically
- [x] `bucketInterval` field matches selected interval
- [x] Handles entries with no timestamps (skip them)
- [x] Tests for each interval threshold
