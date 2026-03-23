# 04 — Format Auto-Detector

**Complexity:** Moderate
**Phase:** 2 — Backend Analyzer
**Blocked by:** 03-analyzer-types
**Blocks:** 09-analyzer-orchestrator

## Objective

Auto-detect whether uploaded log file is JSON or plaintext by peeking at first lines.

## Scope

- `internal/analyzer/detector.go`
- Peek first N non-empty lines (e.g., 10)
- Try `json.Valid()` on each line — if majority valid JSON → "json", else "plaintext"
- Must work on `io.Reader` without consuming the stream (use `io.TeeReader` or buffer+replay)
- Return detected format + a new reader that replays peeked bytes
- Fallback: if ambiguous/mixed, default to "plaintext"
- `internal/analyzer/detector_test.go`

## Edge Cases

- Empty file → plaintext
- Single-line file
- JSON lines with leading whitespace
- Mixed format (some JSON, some not)

## Acceptance Criteria

- [x] Correctly detects pure JSON logs
- [x] Correctly detects plaintext logs
- [x] Falls back to plaintext on mixed input
- [x] Does not consume bytes from the reader (downstream parser gets full stream)
- [x] Tests cover all edge cases above
