# 05 — Streaming Line-by-Line Parser

**Complexity:** Complex
**Phase:** 2 — Backend Analyzer
**Blocked by:** 03-analyzer-types, 04-format-detector
**Blocks:** 09-analyzer-orchestrator

## Objective

Parse log files line-by-line via `bufio.Scanner` over `io.Reader` — never `ReadAll`. Supports both JSON and plaintext formats.

## Scope

- `internal/analyzer/parser.go`
- Accept `io.Reader` + detected format → stream `LogEntry` values
- **JSON parsing:** extract `timestamp`, `level`, `message`, `source` (check fields: `source`, `service`, `module`, `logger`, `component`)
- **Plaintext parsing:** regex-based extraction of timestamp, level (INFO/WARN/ERROR/DEBUG), bracketed source (e.g., `[myservice]`), message remainder
- Return channel or callback pattern for streaming entries
- Accept `context.Context` — stop on cancellation
- `internal/analyzer/parser_test.go`

## Memory Safety

- `bufio.Scanner` wraps `io.Reader` — no full-file buffering
- Each line parsed independently, no accumulation except capped output

## Edge Cases

- Malformed JSON lines (skip, don't fail)
- Lines without timestamps or levels
- Very long lines (Scanner max token size)
- Empty lines

## Acceptance Criteria

- [x] Parses JSON log lines extracting all fields
- [x] Parses plaintext log lines extracting all fields
- [x] Skips malformed lines without failing
- [x] Respects context cancellation
- [x] Never loads full file into memory
- [x] Tests for both formats + edge cases
