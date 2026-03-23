# 03 — Analyzer Types

**Complexity:** Simple
**Phase:** 2 — Backend Analyzer
**Blocked by:** 01-project-scaffolding
**Blocks:** 04-format-detector, 05-streaming-parser, 06-summary-aggregator, 07-pattern-detector, 08-time-series, 09-analyzer-orchestrator

## Objective

Define all shared data types for the analyzer package.

## Scope

- `internal/analyzer/types.go`
- Structs: `LogEntry`, `Summary`, `Pattern`, `TimeBucket`, `AnalysisResult`
- `LogEntry`: timestamp, level, message, source, raw line, line number
- `Summary`: totalLines, level counts (error/warn/info/debug), timeRange, topSources
- `Pattern`: template string, count, sample line
- `TimeBucket`: timestamp, count, errorCount
- `AnalysisResult`: formatDetected, summary, entries (capped 10K), patterns, timeSeries, bucketInterval
- Constants: `MaxEntries = 10_000`, `MaxPatterns = 10_000`

## Acceptance Criteria

- [x] All structs compile with correct JSON tags
- [x] Constants defined for entry/pattern caps
- [x] Types align with API response contract in architecture doc
