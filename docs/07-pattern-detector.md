# 07 — Pattern Detector

**Complexity:** Complex
**Phase:** 2 — Backend Analyzer
**Blocked by:** 03-analyzer-types
**Blocks:** 09-analyzer-orchestrator

## Objective

Detect recurring log patterns via token-masking: replace variable parts with placeholders, group identical templates, rank by count.

## Scope

- `internal/analyzer/pattern.go`
- Token masking replacements:
  - Numbers → `{NUM}`
  - UUIDs → `{UUID}`
  - IP addresses → `{IP}`
  - Hex strings → `{HEX}`
- Group messages by masked template string
- Track count + one sample original line per template
- **Cap:** 10K unique templates max (stop inserting new, only increment existing)
- Runs on full file stream (not just capped entries) for accuracy
- Return top N patterns sorted by count
- `internal/analyzer/pattern_test.go`

## Edge Cases

- Lines with no variable parts (template = original)
- Lines that are all variable parts
- Template map at capacity — new unique templates ignored

## Acceptance Criteria

- [x] Correctly masks numbers, UUIDs, IPs, hex strings
- [x] Groups identical templates and counts occurrences
- [x] Retains one sample line per pattern
- [x] Respects 10K template cap
- [x] Returns patterns sorted by count descending
- [x] Tests cover masking, grouping, cap behavior
