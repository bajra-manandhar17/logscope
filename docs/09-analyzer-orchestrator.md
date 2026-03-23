# 09 — Analyzer Orchestrator

**Complexity:** Moderate
**Phase:** 2 — Backend Analyzer
**Blocked by:** 04-format-detector, 05-streaming-parser, 06-summary-aggregator, 07-pattern-detector, 08-time-series
**Blocks:** 10-analyze-handler

## Objective

Wire together detect → parse → summarize → detect patterns → time-series into a single `Analyze(ctx, reader)` call.

## Scope

- `internal/analyzer/analyzer.go`
- Single function: `Analyze(ctx context.Context, r io.Reader, formatHint string) (*AnalysisResult, error)`
- Flow:
  1. Detect format (or use hint if provided)
  2. Stream-parse entries
  3. Cap entries at `MaxEntries` for response
  4. Feed all entries to summary aggregator
  5. Feed all entry messages to pattern detector
  6. Compute time-series from entries
  7. Assemble `AnalysisResult`
- Accept `context.Context` — propagate to all sub-components
- `internal/analyzer/analyzer_test.go` — integration test with sample log data

## Acceptance Criteria

- [x] Returns complete `AnalysisResult` for JSON input
- [x] Returns complete `AnalysisResult` for plaintext input
- [x] Respects entry cap (10K)
- [x] Respects context cancellation
- [x] Error cases return meaningful errors
- [x] Integration test with realistic sample data
